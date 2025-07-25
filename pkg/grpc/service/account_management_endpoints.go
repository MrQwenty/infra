package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"regexp"

	"github.com/influenzanet/user-management-service/pkg/api"
	"github.com/influenzanet/user-management-service/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_VERIFICATION_ATTEMPTS = 3
	VERIFICATION_CODE_EXPIRY_MINUTES = 10
	MAX_RETRY_ATTEMPTS = 3
)

type VerificationAttempt struct {
	PhoneNumber    string
	Code           string
	Token          string
	Attempts       int
	MaxAttempts    int
	CreatedAt      time.Time
	ExpiresAt      time.Time
	Status         string // "pending", "verified", "expired", "failed"
	RetryCount     int
	MaxRetries     int
}

// In-memory storage for verification attempts (in production, use Redis or database)
var verificationAttempts = make(map[string]*VerificationAttempt)

func (s *userManagementServer) AddPhoneNumber(ctx context.Context, req *api.AddPhoneNumberRequest) (*api.AddPhoneNumberResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	userID, instanceID, err := s.ValidateToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if req.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number cannot be empty")
	}

	// Validate phone number format
	if !s.isValidPhoneNumber(req.PhoneNumber) {
		return nil, status.Error(codes.InvalidArgument, "phone not valid")
	}

	user, err := s.userDBservice.GetUser(instanceID, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Account.AccountConfirmedAt <= 0 {
		return nil, status.Error(codes.FailedPrecondition, "account not confirmed")
	}

	verificationToken, err := s.generateVerificationToken()
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate verification token")
	}

	verificationCode := s.generateVerificationCode()

	// Create verification attempt
	attempt := &VerificationAttempt{
		PhoneNumber: req.PhoneNumber,
		Code:        verificationCode,
		Token:       verificationToken,
		Attempts:    0,
		MaxAttempts: MAX_VERIFICATION_ATTEMPTS,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(VERIFICATION_CODE_EXPIRY_MINUTES * time.Minute),
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  MAX_RETRY_ATTEMPTS,
	}

	verificationAttempts[verificationToken] = attempt

	// Send verification code based on method
	verificationMethod := req.VerificationMethod
	if verificationMethod == "" {
		verificationMethod = "whatsapp" // default
	}

	err = s.sendVerificationCode(req.PhoneNumber, verificationCode, verificationMethod, 0)
	if err != nil {
		log.Printf("Error sending verification: %v", err)
		delete(verificationAttempts, verificationToken)
		return nil, status.Error(codes.Internal, "failed to send verification code")
	}

	return &api.AddPhoneNumberResponse{
		Success:           true,
		VerificationToken: verificationToken,
	}, nil
}

func (s *userManagementServer) EditPhoneNumber(ctx context.Context, req *api.EditPhoneNumberRequest) (*api.EditPhoneNumberResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	userID, instanceID, err := s.ValidateToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if req.NewPhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "new phone number cannot be empty")
	}

	// Validate phone number format
	if !s.isValidPhoneNumber(req.NewPhoneNumber) {
		return nil, status.Error(codes.InvalidArgument, "phone not valid")
	}

	user, err := s.userDBservice.GetUser(instanceID, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Account.AccountConfirmedAt <= 0 {
		return nil, status.Error(codes.FailedPrecondition, "account not confirmed")
	}

	// Check if user has a phone number to edit
	hasPhone := false
	for _, contact := range user.ContactInfos {
		if contact.Type == "phone" {
			hasPhone = true
			break
		}
	}
	if !hasPhone {
		return nil, status.Error(codes.InvalidArgument, "user has no phone number to edit")
	}

	verificationToken, err := s.generateVerificationToken()
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate verification token")
	}

	verificationCode := s.generateVerificationCode()

	// Create verification attempt
	attempt := &VerificationAttempt{
		PhoneNumber: req.NewPhoneNumber,
		Code:        verificationCode,
		Token:       verificationToken,
		Attempts:    0,
		MaxAttempts: MAX_VERIFICATION_ATTEMPTS,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(VERIFICATION_CODE_EXPIRY_MINUTES * time.Minute),
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  MAX_RETRY_ATTEMPTS,
	}

	verificationAttempts[verificationToken] = attempt

	// Send verification code based on method
	verificationMethod := req.VerificationMethod
	if verificationMethod == "" {
		verificationMethod = "whatsapp" // default
	}

	err = s.sendVerificationCode(req.NewPhoneNumber, verificationCode, verificationMethod, 0)
	if err != nil {
		log.Printf("Error sending verification: %v", err)
		delete(verificationAttempts, verificationToken)
		return nil, status.Error(codes.Internal, "failed to send verification code")
	}

	return &api.EditPhoneNumberResponse{
		Success:           true,
		VerificationToken: verificationToken,
	}, nil
}

func (s *userManagementServer) VerifyPhoneNumber(ctx context.Context, req *api.VerifyPhoneNumberRequest) (*api.VerifyPhoneNumberResponse, error) {
	if req == nil || req.Token == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	attempt, exists := verificationAttempts[req.Token]
	if !exists {
		return &api.VerifyPhoneNumberResponse{
			Success:           false,
			Message:           "Invalid or expired verification token",
			Verified:          false,
			AttemptsRemaining: 0,
		}, nil
	}

	// Check if verification has expired
	if time.Now().After(attempt.ExpiresAt) {
		attempt.Status = "expired"
		delete(verificationAttempts, req.Token)
		return &api.VerifyPhoneNumberResponse{
			Success:           false,
			Message:           "Verification code has expired",
			Verified:          false,
			AttemptsRemaining: 0,
		}, nil
	}

	// Check if max attempts reached
	if attempt.Attempts >= attempt.MaxAttempts {
		attempt.Status = "failed"
		delete(verificationAttempts, req.Token)
		return &api.VerifyPhoneNumberResponse{
			Success:           false,
			Message:           "Maximum verification attempts exceeded",
			Verified:          false,
			AttemptsRemaining: 0,
		}, nil
	}

	// Increment attempt counter
	attempt.Attempts++

	// Verify the code
	if attempt.Code == req.Code {
		attempt.Status = "verified"
		
		// Update user's phone number in database
		userID, instanceID, err := s.ValidateToken(req.Token)
		if err == nil {
			// For add phone operation
			if err := s.userDBservice.AddPhoneNumber(instanceID, userID, attempt.PhoneNumber, req.Token); err != nil {
				log.Printf("Error updating phone number in database: %v", err)
			}
		}
		
		// Clean up successful verification
		go func() {
			time.Sleep(5 * time.Second)
			delete(verificationAttempts, req.Token)
		}()

		return &api.VerifyPhoneNumberResponse{
			Success:           true,
			Message:           "Phone number verified successfully",
			Verified:          true,
			AttemptsRemaining: 0,
		}, nil
	}

	attemptsRemaining := attempt.MaxAttempts - attempt.Attempts
	
	if attemptsRemaining == 0 {
		attempt.Status = "failed"
		delete(verificationAttempts, req.Token)
		return &api.VerifyPhoneNumberResponse{
			Success:           false,
			Message:           "Invalid verification code. Maximum attempts exceeded.",
			Verified:          false,
			AttemptsRemaining: 0,
		}, nil
	}

	return &api.VerifyPhoneNumberResponse{
		Success:           false,
		Message:           fmt.Sprintf("Invalid verification code. %d attempts remaining.", attemptsRemaining),
		Verified:          false,
		AttemptsRemaining: int32(attemptsRemaining),
	}, nil
}

func (s *userManagementServer) ResendVerificationCode(ctx context.Context, req *api.ResendVerificationCodeRequest) (*api.ResendVerificationCodeResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	attempt, exists := verificationAttempts[req.Token]
	if !exists {
		return nil, status.Error(codes.NotFound, "Invalid or expired verification token")
	}

	if time.Now().After(attempt.ExpiresAt) {
		attempt.Status = "expired"
		delete(verificationAttempts, req.Token)
		return nil, status.Error(codes.DeadlineExceeded, "Verification session has expired")
	}

	// Generate new code and reset attempts
	newCode := s.generateVerificationCode()
	attempt.Code = newCode
	attempt.Attempts = 0
	attempt.RetryCount = 0
	attempt.ExpiresAt = time.Now().Add(VERIFICATION_CODE_EXPIRY_MINUTES * time.Minute)

	err := s.sendVerificationCode(attempt.PhoneNumber, newCode, "whatsapp", 0)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to resend verification code")
	}

	return &api.ResendVerificationCodeResponse{
		Success:           true,
		VerificationToken: req.Token,
		Message:           "New verification code sent",
		ExpiresAt:         attempt.ExpiresAt.Unix(),
		AttemptsRemaining: int32(MAX_VERIFICATION_ATTEMPTS),
	}, nil
}

func (s *userManagementServer) CancelVerification(ctx context.Context, req *api.CancelVerificationRequest) (*api.CancelVerificationResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	_, exists := verificationAttempts[req.Token]
	if exists {
		delete(verificationAttempts, req.Token)
	}

	return &api.CancelVerificationResponse{
		Success: true,
		Message: "Verification cancelled successfully",
	}, nil
}

func (s *userManagementServer) isValidPhoneNumber(phoneNumber string) bool {
	// E.164 format validation: +[country code][number]
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phoneNumber)
}

func (s *userManagementServer) generateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("whatsapp_%s_%d", hex.EncodeToString(bytes), time.Now().Unix()), nil
}

func (s *userManagementServer) generateVerificationCode() string {
	// Generate a 6-digit code
	code := make([]byte, 6)
	for i := range code {
		code[i] = byte('0' + (rand.Int() % 10))
	}
	return string(code)
}

func (s *userManagementServer) sendVerificationCode(phoneNumber, code, method string, retryCount int) error {
	switch method {
	case "whatsapp":
		return s.sendWhatsAppVerification(phoneNumber, code, retryCount)
	case "sms":
		return s.sendSMSVerification(phoneNumber, code, retryCount)
	default:
		return fmt.Errorf("unsupported verification method: %s", method)
	}
}

func (s *userManagementServer) sendWhatsAppVerification(phoneNumber, code string, retryCount int) error {
	// Simulate WhatsApp API call with retry logic
	message := fmt.Sprintf("Your InfluenzaNet verification code is: %s. This code will expire in %d minutes.", code, VERIFICATION_CODE_EXPIRY_MINUTES)
	
	log.Printf("Sending WhatsApp verification to %s (attempt %d): %s", phoneNumber, retryCount+1, message)
	
	// Simulate API call failure for demonstration
	if retryCount < MAX_RETRY_ATTEMPTS {
		// In real implementation, make actual WhatsApp API call here
		// For now, we'll simulate success
		return nil
	}
	
	return fmt.Errorf("failed to send WhatsApp verification after %d attempts", MAX_RETRY_ATTEMPTS)
}

func (s *userManagementServer) sendSMSVerification(phoneNumber, code string, retryCount int) error {
	// Simulate SMS API call with retry logic
	message := fmt.Sprintf("Your InfluenzaNet verification code is: %s. This code will expire in %d minutes.", code, VERIFICATION_CODE_EXPIRY_MINUTES)
	
	log.Printf("Sending SMS verification to %s (attempt %d): %s", phoneNumber, retryCount+1, message)
	
	// In real implementation, make actual SMS API call here
	return nil
}