package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetUsecase struct {
	UserRepo          domain.IUserRepository
	EmailService      domain.IEmailService
	PasswordResetRepo domain.IPasswordResetRepository
	PasswordExpiry    time.Duration
}

func NewPasswordResetUsecase(repo domain.IPasswordResetRepository, userRepo domain.IUserRepository, emailService domain.IEmailService, expiry time.Duration) domain.IPasswordResetUsecase {
	return &PasswordResetUsecase{
		PasswordResetRepo: repo,
		UserRepo:          userRepo,
		EmailService:      emailService,
		PasswordExpiry:    expiry,
	}
}

func (u *PasswordResetUsecase) FindByEmail(email string) (*domain.PasswordResetToken, error) {
	return u.PasswordResetRepo.FindByEmail(context.Background(), email)
}

func (u *PasswordResetUsecase) MarkAsUsed(token *domain.PasswordResetToken) error {
	return u.PasswordResetRepo.MarkAsUsed(context.Background(), token)
}

func (u *PasswordResetUsecase) RequestReset(email string) error {
	user, err := u.UserRepo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return err
	}

	var token string
	// Check if a reset token already exists for the user
	existingToken, err := u.PasswordResetRepo.FindByEmail(context.Background(), user.Email)
	if err == nil && existingToken != nil {
		if !existingToken.Used && existingToken.ExpiresAt.After(time.Now()) {
			return errors.New("a reset token has already been requested. Please wait until it expires or use the existing token.")
		}
		if existingToken.RateLimit >= 5 {
			if time.Since(existingToken.CreatedAt) < 24*time.Hour {
				return errors.New("rate limit exceeded. You can only request up to 5 password resets in a day.")
			}
			existingToken.RateLimit = 0
		}

		// Update the existing token
		existingToken.TokenHash, _ = security.HashToken(uuid.NewString())
		existingToken.ExpiresAt = time.Now().Add(u.PasswordExpiry)
		existingToken.Used = false
		existingToken.RateLimit++

		token = existingToken.TokenHash
		if err := u.PasswordResetRepo.UpdateResetToken(context.Background(), existingToken); err != nil {
			return err
		}
	} else {
		// Generate a new reset token
		plainToken := uuid.NewString()
		hashedToken, _ := security.HashToken(plainToken)

		token = plainToken

		// Store token
		resetToken := &domain.PasswordResetToken{
			Email:     user.Email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(u.PasswordExpiry),
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now(),
		}
		if err := u.PasswordResetRepo.SaveResetToken(context.Background(), resetToken); err != nil {
			return err
		}
	}

	// Send email
	body := `<h1 style="color: #333; font-family: Arial, sans-serif;">Password Reset Request</h1>
<p style="font-family: Arial, sans-serif; color: #555;">Dear ` + user.FirstName + " " + user.LastName + `,</p>
<p style="font-family: Arial, sans-serif; color: #555;">Use the following token to reset your password for the Blog Platform:</p>
<p style="font-family: Arial, sans-serif; color: #1a73e8; font-weight: bold;">` + token + `</p>
<p style="font-family: Arial, sans-serif; color: #555;">This token expires in ` + u.PasswordExpiry.String() + `. If you did not request this, please ignore this email.</p>
<p style="font-family: Arial, sans-serif; color: #555;">Best regards,</p>
<p style="font-family: Arial, sans-serif; color: #555;">The Blog Platform Team</p>`
	return u.EmailService.SendEmail(context.Background(), user.Email, "Password Reset Request", body)
}

func (u *PasswordResetUsecase) ResetPassword(email, token, newPassword string) error {
	// Validate token
	resetToken, err := u.PasswordResetRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return err
	}

	// Check if token is expired or already used
	if resetToken.Used || resetToken.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired reset token")
	}
	// validate db token with provided token
	if resetToken.TokenHash != token {
		return errors.New("invalid reset token")
	}

	// Get user
	user, err := u.UserRepo.GetUserByEmail(context.Background(), resetToken.Email)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Update user
	if err := u.UserRepo.UpdateUser(context.Background(), user.ID, user); err != nil {
		return err
	}

	// Mark token as used
	resetToken.Used = true
	if err := u.PasswordResetRepo.MarkAsUsed(context.Background(), resetToken); err != nil {
		return err
	}

	// Delete token
	return u.PasswordResetRepo.DeleteResetToken(context.Background(), token)
}
