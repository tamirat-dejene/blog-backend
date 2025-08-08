package usecases

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"math/rand"
	"time"
)

type OTPUsecase struct {
	OTPRepo            domain.IOTPRepository
	EmailService       domain.IEmailService
	ctxtimeout         time.Duration
	otpExpiration      time.Duration
	otpMaximumAttempts int
	secretSalt         string
}

func NewOTPUsecase(repo domain.IOTPRepository, emailService domain.IEmailService, timeout time.Duration, expiration time.Duration, maxAttempts int, secretSalt string) domain.IOTPUsecase {
	return &OTPUsecase{
		OTPRepo:            repo,
		ctxtimeout:         timeout,
		EmailService:       emailService,
		otpExpiration:      expiration,
		otpMaximumAttempts: maxAttempts,
		secretSalt:         secretSalt,
	}
}

// RequestOTP
func (otpuc *OTPUsecase) RequestOTP(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), otpuc.ctxtimeout)
	defer cancel()

	otp := &domain.OTP{
		Email: email,
	}
	code, otpExist, err := otpuc.generateRegistrationOTP(ctx, otp)
	if err != nil {
		return err
	}

	// send otp for email
	body := fmt.Sprintf(`
		<html>
		<body>
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; background-color: #f9f9f9;">
				<h2 style="color: #333;">Your OTP Code</h2>
				<p>Use the following OTP code to complete your registration:</p>
				<h1 style="color: #007BFF; text-align: center;">%s</h1>
				<p>This code will expire in %d minutes.</p>
				<p style="color: #666; text-align: center;">If you did not request this code, please ignore this email.</p>
			</div>
		</body>
		</html>
	`, code, int(otpuc.otpExpiration.Minutes()))

	err = otpuc.EmailService.SendEmail(ctx, email, "Your OTP Code", body)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}
	if otpExist {
		err = otpuc.OTPRepo.UpdateOTPByID(ctx, otp)
		if err != nil {
			return err
		}
		return nil
	}

	err = otpuc.OTPRepo.SaveOTP(ctx, otp)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}
	return nil
}

// GenerateOTP generates a new OTP for the given email and purpose
func (otpuc *OTPUsecase) generateRegistrationOTP(ctx context.Context, otp *domain.OTP) (string, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, otpuc.ctxtimeout)
	defer cancel()

	// Check if the user already requested an OTP
	existingOTP, err := otpuc.OTPRepo.FindOTPByEmail(ctx, otp.Email)
	if err != nil && err != domain.ErrOTPNotFound {
		return "", existingOTP != nil, err
	}

	if existingOTP != nil {
		otp.ID = existingOTP.ID
		// Check if the OTP is still valid
		if time.Now().Before(existingOTP.ExpiresAt) {
			return "", existingOTP != nil, domain.ErrOTPStillValid
		}

		// Check if the user has reached the maximum attempts
		if existingOTP.Attempts >= otpuc.otpMaximumAttempts {
			if time.Since(existingOTP.CreatedAt) < 24*time.Hour {
				return "", existingOTP != nil, domain.ErrOTPMaxAttempts
			}
			// Reset attempts after 24 hours
			existingOTP.Attempts = 0
		}
	}

	// Set the expiration time for the new OTP
	otpRandCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	code := otpRandCode + otpuc.secretSalt
	otp.ExpiresAt = time.Now().Add(otpuc.otpExpiration)
	otp.CreatedAt = time.Now()
	otp.CodeHash = security.HashOTPCode(code)

	// Increment the attempts
	if existingOTP != nil {
		otp.Attempts = existingOTP.Attempts + 1
	} else {
		otp.Attempts = 1
	}
	return otpRandCode, existingOTP != nil, nil
}

// delete by id
func (otpuc *OTPUsecase) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), otpuc.ctxtimeout)
	defer cancel()

	return otpuc.OTPRepo.DeleteOTPByID(ctx, id)
}

// find by email
func (otpuc *OTPUsecase) findByEmail(ctx context.Context, email string) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, otpuc.ctxtimeout)
	defer cancel()

	otp, err := otpuc.OTPRepo.FindOTPByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if otp == nil {
		return nil, domain.ErrOTPNotFound
	}

	return otp, nil
}

// VerifyOTP verifies the OTP for the given email
func (otpuc *OTPUsecase) VerifyOTP(email, code string) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), otpuc.ctxtimeout)
	defer cancel()

	otp, err := otpuc.findByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrOTPNotFound
	}

	if time.Now().After(otp.ExpiresAt) {
		if err := otpuc.DeleteByID(otp.ID); err != nil {
			return nil, fmt.Errorf("failed to delete expired OTP: %w", err)
		}
		return nil, domain.ErrOTPExpired
	}

	if security.VerifyOTPCode(otp.CodeHash, code+otpuc.secretSalt) {
		return otp, nil
	}

	return nil, domain.ErrOTPInvalidCode
}
