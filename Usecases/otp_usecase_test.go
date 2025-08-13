package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"g6/blog-api/Infrastructure/security"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// OTPUsecaseSuite defines the test suite for OTPUsecase
type OTPUsecaseSuite struct {
	suite.Suite
	mockOTPRepo   *domain_mocks.MockIOTPRepository
	mockEmail     *domain_mocks.MockIEmailService
	usecase       *OTPUsecase
	timeout       time.Duration
	otpExpiration time.Duration
	maxAttempts   int
	secretSalt    string
}

func (s *OTPUsecaseSuite) SetupTest() {
	s.mockOTPRepo = domain_mocks.NewMockIOTPRepository(s.T())
	s.mockEmail = domain_mocks.NewMockIEmailService(s.T())
	s.timeout = 3 * time.Second
	s.otpExpiration = 10 * time.Minute
	s.maxAttempts = 5
	s.secretSalt = "test-salt"
	s.usecase = &OTPUsecase{
		OTPRepo:            s.mockOTPRepo,
		EmailService:       s.mockEmail,
		ctxtimeout:         s.timeout,
		otpExpiration:      s.otpExpiration,
		otpMaximumAttempts: s.maxAttempts,
		secretSalt:         s.secretSalt,
	}
}

func TestOTPUsecaseSuite(t *testing.T) {
	suite.Run(t, new(OTPUsecaseSuite))
}

func (s *OTPUsecaseSuite) TestRequestOTP() {
	s.Run("SuccessNewOTP", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Your OTP Code", mock.Anything).Return(nil)
		s.mockOTPRepo.On("SaveOTP", mock.Anything, mock.MatchedBy(func(otp *domain.OTP) bool {
			return otp.Email == email &&
				otp.Attempts == 1 &&
				!otp.CreatedAt.IsZero() &&
				time.Until(otp.ExpiresAt) <= s.otpExpiration
		})).Return(nil)

		err := s.usecase.RequestOTP(email)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("SuccessExistingExpiredOTP", func() {
		email := "test@example.com"

		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "old-hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Your OTP Code", mock.Anything).Return(nil)
		s.mockOTPRepo.On("UpdateOTPByID", mock.Anything, mock.MatchedBy(func(otp *domain.OTP) bool {
			return otp.ID == existingOTP.ID &&
				otp.Email == email &&
				otp.Attempts == 2 &&
				!otp.CreatedAt.IsZero() &&
				time.Until(otp.ExpiresAt) <= s.otpExpiration
		})).Return(nil)
		err := s.usecase.RequestOTP(email)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("ExistingValidOTP", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "hashed-code",
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)

		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Equal(domain.ErrOTPStillValid, err)
		s.resetMocks()
	})

	s.Run("MaxAttemptsExceeded", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  s.maxAttempts,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)

		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Equal(domain.ErrOTPMaxAttempts, err)
		s.resetMocks()
	})

	s.Run("FindOTPError", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, errors.New("find failed"))

		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Equal("find failed", err.Error())
		s.resetMocks()
	})

	s.Run("SaveOTPError", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)
		s.mockOTPRepo.On("SaveOTP", mock.Anything, mock.Anything).Return(errors.New("save failed"))
		s.mockEmail.On("SendEmail", mock.Anything, email, "Your OTP Code", mock.Anything).Return(nil)
		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Contains(err.Error(), "failed to save OTP: save failed")
		s.resetMocks()
	})

	s.Run("UpdateOTPError", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "old-hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)
		s.mockOTPRepo.On("UpdateOTPByID", mock.Anything, mock.Anything).Return(errors.New("update failed"))
		s.mockEmail.On("SendEmail", mock.Anything, email, "Your OTP Code", mock.Anything).Return(nil)
		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})

	s.Run("SendEmailError", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Your OTP Code", mock.Anything).Return(errors.New("email failed"))
		s.mockOTPRepo.On("SaveOTP", mock.Anything, mock.Anything).Return(nil)
		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Contains(err.Error(), "failed to send OTP email: email failed")
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(nil, errors.New("context deadline exceeded"))

		err := s.usecase.RequestOTP(email)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *OTPUsecaseSuite) TestGenerateRegistrationOTP() {
	s.Run("SuccessNewOTP", func() {
		email := "test@example.com"
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)
		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.NoError(err)
		s.Equal(6, len(code))
		s.False(otpExist)
		s.Equal(email, otp.Email)
		s.Equal(security.HashOTPCode(code+s.secretSalt), otp.CodeHash)
		s.Equal(1, otp.Attempts)
		s.False(otp.CreatedAt.IsZero())
		s.True(time.Until(otp.ExpiresAt) <= s.otpExpiration)
		s.resetMocks()
	})

	s.Run("SuccessExistingExpiredOTP", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "old-hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		}
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)
		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.NoError(err)
		s.Equal(6, len(code))
		s.True(otpExist)
		s.Equal(existingOTP.ID, otp.ID)
		s.Equal(email, otp.Email)
		s.Equal(security.HashOTPCode(code+s.secretSalt), otp.CodeHash)
		s.Equal(2, otp.Attempts)
		s.False(otp.CreatedAt.IsZero())
		s.True(time.Until(otp.ExpiresAt) <= s.otpExpiration)
		s.resetMocks()
	})

	s.Run("ValidOTP", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "hashed-code",
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)

		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.Error(err)
		s.Equal(domain.ErrOTPStillValid, err)
		s.Equal("", code)
		s.True(otpExist)
		s.resetMocks()
	})

	s.Run("MaxAttemptsExceeded", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  s.maxAttempts,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)

		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.Error(err)
		s.Equal(domain.ErrOTPMaxAttempts, err)
		s.Equal("", code)
		s.True(otpExist)
		s.resetMocks()
	})

	s.Run("ResetAttemptsAfter24Hours", func() {
		email := "test@example.com"
		existingOTP := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "old-hashed-code",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  s.maxAttempts,
			CreatedAt: time.Now().Add(-25 * time.Hour),
		}
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(existingOTP, nil)
		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.NoError(err)
		s.Equal(6, len(code))
		s.True(otpExist)
		s.Equal(existingOTP.ID, otp.ID)
		s.Equal(email, otp.Email)
		s.Equal(security.HashOTPCode(code+s.secretSalt), otp.CodeHash)
		s.Equal(1, otp.Attempts)
		s.False(otp.CreatedAt.IsZero())
		s.True(time.Until(otp.ExpiresAt) <= s.otpExpiration)
		s.resetMocks()
	})

	s.Run("FindOTPError", func() {
		email := "test@example.com"
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, errors.New("find failed"))

		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.Error(err)
		s.Equal("find failed", err.Error())
		s.Equal("", code)
		s.False(otpExist)
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		email := "test@example.com"
		otp := &domain.OTP{Email: email}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(nil, errors.New("context deadline exceeded"))

		code, otpExist, err := s.usecase.generateRegistrationOTP(context.Background(), otp)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.Equal("", code)
		s.False(otpExist)
		s.resetMocks()
	})
}

func (s *OTPUsecaseSuite) TestDeleteByID() {
	s.Run("Success", func() {
		id := "otp1"
		s.mockOTPRepo.On("DeleteOTPByID", mock.Anything, id).Return(nil)

		err := s.usecase.DeleteByID(id)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("DeleteError", func() {
		id := "otp1"
		s.mockOTPRepo.On("DeleteOTPByID", mock.Anything, id).Return(errors.New("delete failed"))

		err := s.usecase.DeleteByID(id)

		s.Error(err)
		s.Equal("delete failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		id := "otp1"
		s.mockOTPRepo.On("DeleteOTPByID", mock.Anything, id).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(errors.New("context deadline exceeded"))

		err := s.usecase.DeleteByID(id)

		s.Error(err)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *OTPUsecaseSuite) TestFindByEmail() {
	s.Run("Success", func() {
		email := "test@example.com"
		otp := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  "hashed-code",
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(otp, nil)

		result, err := s.usecase.findByEmail(context.Background(), email)

		s.NoError(err)
		s.Equal(otp, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)

		result, err := s.usecase.findByEmail(context.Background(), email)

		s.Error(err)
		s.Nil(result)
		s.Equal(domain.ErrOTPNotFound, err)
		s.resetMocks()
	})

	s.Run("FindError", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, errors.New("find failed"))

		result, err := s.usecase.findByEmail(context.Background(), email)

		s.Error(err)
		s.Nil(result)
		s.Equal("find failed", err.Error())
		s.resetMocks()
	})

	s.Run("ContextTimeout", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Run(func(args mock.Arguments) {
			time.Sleep(4 * time.Second) // Exceed the 3-second timeout
		}).Return(nil, errors.New("context deadline exceeded"))

		result, err := s.usecase.findByEmail(context.Background(), email)

		s.Error(err)
		s.Nil(result)
		s.Contains(err.Error(), "context deadline exceeded")
		s.resetMocks()
	})
}

func (s *OTPUsecaseSuite) TestVerifyOTP() {
	s.Run("Success", func() {
		email := "test@example.com"
		code := "123456"
		hashedCode := security.HashOTPCode(code + s.secretSalt)
		otp := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  hashedCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(otp, nil)

		result, err := s.usecase.VerifyOTP(email, code)

		s.NoError(err)
		s.Equal(otp, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		email := "test@example.com"
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(nil, domain.ErrOTPNotFound)

		result, err := s.usecase.VerifyOTP(email, "123456")

		s.Error(err)
		s.Nil(result)
		s.Equal(domain.ErrOTPNotFound, err)
		s.resetMocks()
	})

	s.Run("ExpiredOTP", func() {
		email := "test@example.com"
		code := "123456"
		hashedCode := security.HashOTPCode(code + s.secretSalt)
		otp := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  hashedCode,
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(otp, nil)
		s.mockOTPRepo.On("DeleteOTPByID", mock.Anything, otp.ID).Return(nil)

		result, err := s.usecase.VerifyOTP(email, code)

		s.Error(err)
		s.Nil(result)
		s.Equal(domain.ErrOTPExpired, err)
		s.resetMocks()
	})

	s.Run("InvalidCode", func() {
		email := "test@example.com"
		code := "123456"
		hashedCode := security.HashOTPCode("654321" + s.secretSalt)
		otp := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  hashedCode,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(otp, nil)

		result, err := s.usecase.VerifyOTP(email, code)

		s.Error(err)
		s.Nil(result)
		s.Equal(domain.ErrOTPInvalidCode, err)
		s.resetMocks()
	})

	s.Run("DeleteExpiredOTPError", func() {
		email := "test@example.com"
		code := "123456"
		hashedCode := security.HashOTPCode(code + s.secretSalt)
		otp := &domain.OTP{
			ID:        "otp1",
			Email:     email,
			CodeHash:  hashedCode,
			ExpiresAt: time.Now().Add(-1 * time.Minute),
			Attempts:  1,
			CreatedAt: time.Now(),
		}
		s.mockOTPRepo.On("FindOTPByEmail", mock.Anything, email).Return(otp, nil)
		s.mockOTPRepo.On("DeleteOTPByID", mock.Anything, otp.ID).Return(errors.New("delete failed"))

		result, err := s.usecase.VerifyOTP(email, code)

		s.Error(err)
		s.Nil(result)
		s.Contains(err.Error(), "failed to delete expired OTP: delete failed")
		s.resetMocks()
	})
}

func (s *OTPUsecaseSuite) resetMocks() {
	s.mockOTPRepo.ExpectedCalls = nil
	s.mockOTPRepo.Calls = nil
	s.mockEmail.ExpectedCalls = nil
	s.mockEmail.Calls = nil
}
