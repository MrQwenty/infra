package types

type AddPhoneNumberRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type ChangePhoneNumberRequest struct {
	NewPhoneNumber string `json:"newPhoneNumber" binding:"required"`
}

type VerifyWhatsAppRequest struct {
	Token string `json:"token" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type ResendWhatsAppCodeRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Token       string `json:"token" binding:"required"`
}
