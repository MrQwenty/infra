package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
)

type AddPhoneRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type ChangePhoneRequest struct {
	NewPhoneNumber string `json:"newPhoneNumber" binding:"required"`
}

type WhatsAppVerificationRequest struct {
	Token       string `json:"token" binding:"required"`
	Code        string `json:"code" binding:"required"`
	PhoneNumber string `json:"phoneNumber"`
}

type ApiResponse struct {
	Success           bool        `json:"success"`
	Data              interface{} `json:"data,omitempty"`
	Message           string      `json:"message,omitempty"`
	VerificationToken string      `json:"verificationToken,omitempty"`
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

	c.JSON(http.StatusOK, ApiResponse{
		Success: response.Success,
		Message: "WhatsApp verification successful",
	})
}

func (h *UserManagementHandlers) SendWhatsAppCodeHandler(c *gin.Context) {
	var req WhatsAppVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	response, err := h.whatsappClient.SendCode(c, &whatsapp.SendCodeRequest{
		PhoneNumber: req.PhoneNumber,
		Token:       req.Token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Message: "Failed to send WhatsApp code",
		})
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Success: response.Success,
		Message: "WhatsApp code sent successfully",
	})
}
