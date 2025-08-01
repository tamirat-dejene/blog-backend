package domain

import "context"

type IEmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type IPasswordResetUsecase interface {
	ResetPassword(email, token, newPassword string) error
	RequestReset(email string) error
}

type IPasswordResetRepository interface {
	SaveResetToken(ctx context.Context, token *PasswordResetToken) error
	FindByEmail(ctx context.Context, email string) (*PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, token *PasswordResetToken) error
	DeleteResetToken(ctx context.Context, email string) error
	UpdateResetToken(ctx context.Context, token *PasswordResetToken) error
}
