package domain

import "errors"

type DomainError struct {
	Err  error
	Code int
}

var (
	ErrNotFound          = errors.New("not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("User not authenticated or authorized")
	ErrInvalidFile       = errors.New("invalid file format")
	ErrOTPNotFound       = errors.New("OTP not found")
	ErrOTPExpired        = errors.New("OTP expired")
	ErrOTPMaxAttempts    = errors.New("maximum OTP attempts exceeded")
	ErrOTPStillValid     = errors.New("OTP is still valid")
	ErrOTPInvalidCode    = errors.New("invalid OTP code")
	ErrOTPInvalid        = errors.New("invalid OTP")
	ErrOTPFailedToDelete = errors.New("failed to delete OTP")
)
