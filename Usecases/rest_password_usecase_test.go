package usecases

import (
	"errors"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"g6/blog-api/Infrastructure/security"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// PasswordResetUsecaseSuite defines the test suite for PasswordResetUsecase
type PasswordResetUsecaseSuite struct {
	suite.Suite
	mockUserRepo  *domain_mocks.MockIUserRepository
	mockEmail     *domain_mocks.MockIEmailService
	mockResetRepo *domain_mocks.MockIPasswordResetRepository
	usecase       *PasswordResetUsecase
	expiry        time.Duration
}

func (s *PasswordResetUsecaseSuite) SetupTest() {
	s.mockUserRepo = domain_mocks.NewMockIUserRepository(s.T())
	s.mockEmail = domain_mocks.NewMockIEmailService(s.T())
	s.mockResetRepo = domain_mocks.NewMockIPasswordResetRepository(s.T())
	s.expiry = 1 * time.Hour
	s.usecase = &PasswordResetUsecase{
		UserRepo:          s.mockUserRepo,
		EmailService:      s.mockEmail,
		PasswordResetRepo: s.mockResetRepo,
		PasswordExpiry:    s.expiry,
	}
}

func TestPasswordResetUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PasswordResetUsecaseSuite))
}

func (s *PasswordResetUsecaseSuite) TestFindByEmail() {
	s.Run("Success", func() {
		email := "test@example.com"
		expectedToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: "hashed-token",
			ExpiresAt: time.Now().Add(s.expiry),
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now(),
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(expectedToken, nil)

		result, err := s.usecase.FindByEmail(email)

		s.NoError(err)
		s.Equal(expectedToken, result)
		s.resetMocks()
	})

	s.Run("NotFound", func() {
		email := "test@example.com"
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("token not found"))

		result, err := s.usecase.FindByEmail(email)

		s.Error(err)
		s.Nil(result)
		s.Equal("token not found", err.Error())
		s.resetMocks()
	})
}

func (s *PasswordResetUsecaseSuite) TestMarkAsUsed() {
	s.Run("Success", func() {
		token := &domain.PasswordResetToken{
			Email:     "test@example.com",
			TokenHash: "hashed-token",
			ExpiresAt: time.Now().Add(s.expiry),
			Used:      false,
		}
		s.mockResetRepo.On("MarkAsUsed", mock.Anything, token).Return(nil)

		err := s.usecase.MarkAsUsed(token)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("Error", func() {
		token := &domain.PasswordResetToken{
			Email:     "test@example.com",
			TokenHash: "hashed-token",
			ExpiresAt: time.Now().Add(s.expiry),
			Used:      false,
		}
		s.mockResetRepo.On("MarkAsUsed", mock.Anything, token).Return(errors.New("mark failed"))

		err := s.usecase.MarkAsUsed(token)

		s.Error(err)
		s.Equal("mark failed", err.Error())
		s.resetMocks()
	})
}

func (s *PasswordResetUsecaseSuite) TestRequestReset() {
	s.Run("SuccessNewToken", func() {
		email := "test@example.com"
		user := &domain.User{
			ID:        "1",
			Email:     email,
			FirstName: "Test",
			LastName:  "User",
		}

		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("not found"))
		s.mockResetRepo.On("SaveResetToken", mock.Anything, mock.MatchedBy(func(t *domain.PasswordResetToken) bool {
			return t.Email == email &&
				time.Until(t.ExpiresAt) > 0 &&
				!t.Used &&
				t.RateLimit == 1 &&
				!t.CreatedAt.IsZero()
		})).Return(nil)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Password Reset Request", mock.Anything).Return(nil)

		err := s.usecase.RequestReset(email)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("SuccessExistingToken", func() {
		email := "test@example.com"
		user := &domain.User{
			ID:        "1",
			Email:     email,
			FirstName: "Test",
			LastName:  "User",
		}
		existingToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: "old-hashed-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(existingToken, nil)
		s.mockResetRepo.On("UpdateResetToken", mock.Anything, mock.MatchedBy(func(t *domain.PasswordResetToken) bool {
			return t.Email == email &&
				t.Used == false &&
				t.RateLimit == 2 &&
				!t.CreatedAt.IsZero() &&
				time.Until(t.ExpiresAt) <= s.expiry
		})).Return(nil)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Password Reset Request", mock.Anything).Return(nil)
		err := s.usecase.RequestReset(email)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		email := "test@example.com"
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("user not found"))

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("user not found", err.Error())
		s.resetMocks()
	})

	s.Run("ExistingValidToken", func() {
		email := "test@example.com"
		user := &domain.User{ID: "1", Email: email}
		existingToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: "hashed-token",
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now(),
		}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(existingToken, nil)

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("a reset token has already been requested. Please wait until it expires or use the existing token", err.Error())
		s.resetMocks()
	})

	s.Run("RateLimitExceeded", func() {
		email := "test@example.com"
		user := &domain.User{ID: "1", Email: email}
		existingToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: "hashed-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Used:      false,
			RateLimit: 5,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(existingToken, nil)

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("rate limit exceeded. You can only request up to 5 password resets in a day", err.Error())
		s.resetMocks()
	})

	s.Run("SaveTokenError", func() {
		email := "test@example.com"
		user := &domain.User{ID: "1", Email: email}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("not found"))
		s.mockResetRepo.On("SaveResetToken", mock.Anything, mock.Anything).Return(errors.New("save failed"))

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("save failed", err.Error())
		s.resetMocks()
	})

	s.Run("UpdateTokenError", func() {
		email := "test@example.com"
		user := &domain.User{ID: "1", Email: email}
		existingToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: "old-hashed-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(existingToken, nil)
		s.mockResetRepo.On("UpdateResetToken", mock.Anything, mock.Anything).Return(errors.New("update failed"))

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})

	s.Run("SendEmailError", func() {
		email := "test@example.com"
		user := &domain.User{ID: "1", Email: email, FirstName: "Test", LastName: "User"}
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("not found"))
		s.mockResetRepo.On("SaveResetToken", mock.Anything, mock.Anything).Return(nil)
		s.mockEmail.On("SendEmail", mock.Anything, email, "Password Reset Request", mock.Anything).Return(errors.New("email failed"))

		err := s.usecase.RequestReset(email)

		s.Error(err)
		s.Equal("email failed", err.Error())
		s.resetMocks()
	})
}

func (s *PasswordResetUsecaseSuite) TestResetPassword() {
	s.Run("Success", func() {
		email := "test@example.com"
		plainToken := uuid.NewString()
		hashedToken, _ := security.HashToken(plainToken)
		newPassword := "newpassword123"
		user := &domain.User{ID: "1", Email: email}
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
			RateLimit: 1,
			CreatedAt: time.Now(),
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, resetToken.Email).Return(user, nil)
		user.Password = newPassword
		s.mockUserRepo.On("UpdateUser", mock.Anything, user.ID, user).Return(nil)
		s.mockResetRepo.On("MarkAsUsed", mock.Anything, mock.MatchedBy(func(t *domain.PasswordResetToken) bool {
			return t.Email == email && t.Used == true
		})).Return(nil)
		s.mockResetRepo.On("DeleteResetToken", mock.Anything, hashedToken).Return(nil)

		err := s.usecase.ResetPassword(email, hashedToken, newPassword)

		s.NoError(err)
		s.resetMocks()
	})

	s.Run("TokenNotFound", func() {
		email := "test@example.com"
		token := "invalid-token"
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("token not found"))

		err := s.usecase.ResetPassword(email, token, "newpassword123")

		s.Error(err)
		s.Equal("token not found", err.Error())
		s.resetMocks()
	})

	s.Run("ExpiredToken", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)

		err := s.usecase.ResetPassword(email, token, "newpassword123")

		s.Error(err)
		s.Equal("invalid or expired reset token", err.Error())
		s.resetMocks()
	})

	s.Run("UsedToken", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      true,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		err := s.usecase.ResetPassword(email, token, "newpassword123")

		s.Error(err)
		s.Equal("invalid or expired reset token", err.Error())
		s.resetMocks()
	})

	s.Run("InvalidToken", func() {
		email := "test@example.com"
		token := "wrong-token"
		hashedToken, _ := security.HashToken("correct-token")
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)

		err := s.usecase.ResetPassword(email, token, "newpassword123")

		s.Error(err)
		s.Equal("invalid reset token", err.Error())
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, resetToken.Email).Return(nil, errors.New("user not found"))

		err := s.usecase.ResetPassword(email, hashedToken, "newpassword123")

		s.Error(err)
		s.Equal("user not found", err.Error())
		s.resetMocks()
	})

	s.Run("HashPasswordError", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		user := &domain.User{ID: "1", Email: email}
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		newPassword := strings.Repeat("newpassword123", 6)
		err := s.usecase.ResetPassword(email, hashedToken, newPassword)

		s.Error(err)
		s.Equal("bcrypt: password length exceeds 72 bytes", err.Error())
		s.resetMocks()
	})

	s.Run("UpdateUserError", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		user := &domain.User{ID: "1", Email: email}
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, user.ID, mock.Anything).Return(errors.New("update failed"))

		err := s.usecase.ResetPassword(email, hashedToken, "newpassword123")

		s.Error(err)
		s.Equal("update failed", err.Error())
		s.resetMocks()
	})

	s.Run("MarkAsUsedError", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		user := &domain.User{ID: "1", Email: email}
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, user.ID, mock.Anything).Return(nil)
		s.mockResetRepo.On("MarkAsUsed", mock.Anything, mock.Anything).Return(errors.New("mark failed"))

		err := s.usecase.ResetPassword(email, hashedToken, "newpassword123")

		s.Error(err)
		s.Equal("mark failed", err.Error())
		s.resetMocks()
	})

	s.Run("DeleteTokenError", func() {
		email := "test@example.com"
		token := "reset-token"
		hashedToken, _ := security.HashToken(token)
		user := &domain.User{ID: "1", Email: email}
		resetToken := &domain.PasswordResetToken{
			Email:     email,
			TokenHash: hashedToken,
			ExpiresAt: time.Now().Add(30 * time.Minute),
			Used:      false,
		}
		s.mockResetRepo.On("FindByEmail", mock.Anything, email).Return(resetToken, nil)
		s.mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
		s.mockUserRepo.On("UpdateUser", mock.Anything, user.ID, mock.Anything).Return(nil)
		s.mockResetRepo.On("MarkAsUsed", mock.Anything, mock.Anything).Return(nil)
		s.mockResetRepo.On("DeleteResetToken", mock.Anything, hashedToken).Return(errors.New("delete failed"))

		err := s.usecase.ResetPassword(email, hashedToken, "newpassword123")

		s.Error(err)
		s.Equal("delete failed", err.Error())
		s.resetMocks()
	})
}

func (s *PasswordResetUsecaseSuite) resetMocks() {
	s.mockUserRepo.ExpectedCalls = nil
	s.mockUserRepo.Calls = nil
	s.mockEmail.ExpectedCalls = nil
	s.mockEmail.Calls = nil
	s.mockResetRepo.ExpectedCalls = nil
	s.mockResetRepo.Calls = nil
}
