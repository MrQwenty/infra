package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
)

type AddPhoneRequest struct {
	PhoneNumber        string `json:"phoneNumber" binding:"required"`
	VerificationMethod string `json:"verificationMethod,omitempty"`
}

type ChangePhoneRequest struct {
	NewPhoneNumber     string `json:"newPhoneNumber" binding:"required"`
	VerificationMethod string `json:"verificationMethod,omitempty"`
}

type WhatsAppVerificationRequest struct {
	Token       string `json:"token" binding:"required"`
	Code        string `json:"code" binding:"required"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type InitiateVerificationRequest struct {
	PhoneNumber        string `json:"phoneNumber" binding:"required"`
	VerificationMethod string `json:"verificationMethod" binding:"required"`
}

type ResendVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}

type CancelVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}

type ApiResponse struct {
	Success             bool        `json:"success"`
	Data                interface{} `json:"data,omitempty"`
	Message             string      `json:"message,omitempty"`
	VerificationToken   string      `json:"verificationToken,omitempty"`
	ExpiresAt           *time.Time  `json:"expiresAt,omitempty"`
	AttemptsRemaining   *int        `json:"attemptsRemaining,omitempty"`
}

func (h *UserManagementHandlers) AddPhoneNumberHandler(c *gin.Context) {
	var req AddPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, ApiResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	response, err := h.userManagementClient.AddPhoneNumber(c, &api.AddPhoneNumberRequest{
		Token:       token,
		PhoneNumber: req.PhoneNumber,
		VerificationMethod: req.VerificationMethod,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to add phone number",
		})
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Success:           response.Success,
		VerificationToken: response.VerificationToken,
		Message:           "Phone number added successfully",
	})
}

func (h *UserManagementHandlers) ChangePhoneNumberHandler(c *gin.Context) {
	var req ChangePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, ApiResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	response, err := h.userManagementClient.EditPhoneNumber(c, &api.EditPhoneNumberRequest{
		Token:          token,
		NewPhoneNumber: req.NewPhoneNumber,
		VerificationMethod: req.VerificationMethod,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to change phone number",
		})
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Success:           response.Success,
		VerificationToken: response.VerificationToken,
		Message:           "Phone number changed successfully",
	})
}

func (h *UserManagementHandlers) InitiateWhatsAppVerificationHandler(c *gin.Context) {
	var req InitiateVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, ApiResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	response, err := h.whatsappClient.InitiateVerification(c, &whatsapp.InitiateVerificationRequest{
		PhoneNumber:        req.PhoneNumber,
		VerificationMethod: req.VerificationMethod,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to initiate WhatsApp verification",
		})
		return
	}

	expiresAt := time.Unix(response.ExpiresAt, 0)
	attemptsRemaining := int(response.AttemptsRemaining)

	c.JSON(http.StatusOK, ApiResponse{
		Success:           response.Success,
		VerificationToken: response.VerificationToken,
		Message:           response.Message,
		ExpiresAt:         &expiresAt,
		AttemptsRemaining: &attemptsRemaining,
	})
}

func (h *UserManagementHandlers) VerifyWhatsAppHandler(c *gin.Context) {
	var req WhatsAppVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	response, err := h.whatsappClient.VerifyCode(c, &whatsapp.VerifyCodeRequest{
		Token: req.Token,
		Code:  req.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to verify WhatsApp code",
		})
		return
	}

	attemptsRemaining := int(response.AttemptsRemaining)

	c.JSON(http.StatusOK, ApiResponse{
		Success:           response.Success,
		Message:           response.Message,
		AttemptsRemaining: &attemptsRemaining,
	})
}

func (h *UserManagementHandlers) ResendWhatsAppCodeHandler(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	response, err := h.whatsappClient.ResendCode(c, &whatsapp.ResendCodeRequest{
		Token: req.Token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to resend WhatsApp code",
		})
		return
	}

	expiresAt := time.Unix(response.ExpiresAt, 0)
	attemptsRemaining := int(response.AttemptsRemaining)

	c.JSON(http.StatusOK, ApiResponse{
		Success:           response.Success,
		VerificationToken: response.VerificationToken,
		Message:           response.Message,
		ExpiresAt:         &expiresAt,
		AttemptsRemaining: &attemptsRemaining,
	})
}

func (h *UserManagementHandlers) CancelWhatsAppVerificationHandler(c *gin.Context) {
	var req CancelVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	response, err := h.whatsappClient.CancelVerification(c, &whatsapp.CancelVerificationRequest{
		Token: req.Token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to cancel WhatsApp verification",
		})
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Success: response.Success,
		Message: response.Message,
	})
}

func (h *HttpEndpoints) verifyPhoneNumber(c *gin.Context) {
	h.grpcCallHandler(
		c,
		func(c *gin.Context) (protoreflect.ProtoMessage, error) {
			var req umAPI.VerifyPhoneNumberRequest
			if err := h.JsonToProto(c, &req); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			return h.clients.UserManagement.VerifyPhoneNumber(context.Background(), &req)
		},
	)
}

func (h *HttpEndpoints) resendPhoneVerificationCode(c *gin.Context) {
	h.grpcCallHandler(
		c,
		func(c *gin.Context) (protoreflect.ProtoMessage, error) {
			var req umAPI.ResendVerificationCodeRequest
			if err := h.JsonToProto(c, &req); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			return h.clients.UserManagement.ResendVerificationCode(context.Background(), &req)
		},
	)
}

func (h *HttpEndpoints) cancelPhoneVerification(c *gin.Context) {
	h.grpcCallHandler(
		c,
		func(c *gin.Context) (protoreflect.ProtoMessage, error) {
			var req umAPI.CancelVerificationRequest
			if err := h.JsonToProto(c, &req); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			return h.clients.UserManagement.CancelVerification(context.Background(), &req)
		},
	)
}