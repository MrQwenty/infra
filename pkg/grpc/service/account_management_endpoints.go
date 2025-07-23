package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/influenzanet/user-management-service/pkg/api"
	"github.com/influenzanet/user-management-service/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	err = s.userDBservice.AddPhoneNumber(instanceID, userID, req.PhoneNumber, verificationToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.sendWhatsAppVerification(req.PhoneNumber, verificationCode)
	if err != nil {
		log.Printf("Error sending WhatsApp verification: %v", err)
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

	err = s.userDBservice.UpdatePhoneNumber(instanceID, userID, req.NewPhoneNumber, verificationToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.sendWhatsAppVerification(req.NewPhoneNumber, verificationCode)
	if err != nil {
		log.Printf("Error sending WhatsApp verification: %v", err)
	}

	return &api.EditPhoneNumberResponse{
		Success:           true,
		VerificationToken: verificationToken,
	}, nil
}

func (s *userManagementServer) generateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("whatsapp_%s_%d", hex.EncodeToString(bytes), time.Now().Unix()), nil
}

func (s *userManagementServer) generateVerificationCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%02d", int(b)%100)
	}
	return code[:6]
}

func (s *userManagementServer) sendWhatsAppVerification(phoneNumber, code string) error {
	message := fmt.Sprintf("Your InfluenzaNet verification code is: %s", code)
	
	log.Printf("Sending WhatsApp verification to %s: %s", phoneNumber, message)
	
	return nil
}
