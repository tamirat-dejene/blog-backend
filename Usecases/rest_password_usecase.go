package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"time"

	"github.com/google/uuid"
)

type PasswordResetUsecase struct {
	PasswordResetRepo domain.IPasswordResetRepository
	PasswordExpiry    time.Duration
}

func NewPasswordResetUsecase(repo domain.IPasswordResetRepository, expiry time.Duration) domain.IPasswordResetUsecase {
	return &PasswordResetUsecase{
		PasswordResetRepo: repo,
		PasswordExpiry:    expiry,
	}
}

func (u *PasswordResetUsecase) SaveResetToken(email string) error {
	plainToken := uuid.NewString()
	hashedToken, _ := security.HashToken(plainToken)

	resetToken := &domain.PasswordResetToken{
		Email:     email,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(u.PasswordExpiry),
		Used:      false,
	}
	return u.PasswordResetRepo.SaveResetToken(context.Background(), resetToken)
}

func (u *PasswordResetUsecase) FindByEmail(email string) (*domain.PasswordResetToken, error) {
	return u.PasswordResetRepo.FindByEmail(context.Background(), email)
}

func (u *PasswordResetUsecase) MarkAsUsed(token *domain.PasswordResetToken) error {
	return u.PasswordResetRepo.MarkAsUsed(context.Background(), token)
}
