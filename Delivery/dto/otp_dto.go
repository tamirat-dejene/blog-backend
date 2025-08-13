package dto

type VerifyOTPRequest struct {
	Code string `json:"code" binding:"required"`
}
