package domain

import (
	"context"
	"time"
)

type OTP struct {
	ID        string
	Email     string
	CodeHash  string
	ExpiresAt time.Time
	Attempts  int
	CreatedAt time.Time
}

type IOTPUsecase interface {
	RequestOTP(email string) error
	VerifyOTP(email, code string) (*OTP, error)
	DeleteByID(id string) error
}

type IOTPRepository interface {
	SaveOTP(ctx context.Context, otp *OTP) error
	FindOTPByEmail(ctx context.Context, email string) (*OTP, error)
	DeleteOTPByID(ctx context.Context, id string) error
	UpdateOTPByID(ctx context.Context, otp *OTP) error
}
