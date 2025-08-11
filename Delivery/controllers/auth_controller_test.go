package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"g6/blog-api/Infrastructure/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AuthControllerSuite defines the test suite for AuthController
type AuthControllerSuite struct {
	suite.Suite
	mockUserUsecase          *domain_mocks.MockIUserUsecase
	mockAuthService          *domain_mocks.MockIAuthService
	mockOTPUsecase           *domain_mocks.MockIOTPUsecase
	mockRefreshTokenUsecase  *domain_mocks.MockIRefreshTokenUsecase
	mockPasswordResetUsecase *domain_mocks.MockIPasswordResetUsecase
	handler                  *AuthController
	validate                 *validator.Validate
}

// SetupTest initializes the mocks and handler before each test
func (s *AuthControllerSuite) SetupTest() {
	s.mockUserUsecase = domain_mocks.NewMockIUserUsecase(s.T())
	s.mockAuthService = domain_mocks.NewMockIAuthService(s.T())
	s.mockOTPUsecase = domain_mocks.NewMockIOTPUsecase(s.T())
	s.mockRefreshTokenUsecase = domain_mocks.NewMockIRefreshTokenUsecase(s.T())
	s.mockPasswordResetUsecase = domain_mocks.NewMockIPasswordResetUsecase(s.T())

	s.handler = &AuthController{
		UserUsecase:          s.mockUserUsecase,
		AuthService:          s.mockAuthService,
		OTP:                  s.mockOTPUsecase,
		RefreshTokenUsecase:  s.mockRefreshTokenUsecase,
		PasswordResetUsecase: s.mockPasswordResetUsecase,
	}
	s.validate = validator.New()
}

// TestAuthControllerSuite runs the test suite
func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}

// TestRegisterRequest tests the RegisterRequest method
func (s *AuthControllerSuite) TestRegisterRequest() {
	s.Run("Success", func() {
		userRequest := dto.UserRequest{
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
			Provider:  "manual",
		}

		user := dto.ToDomainUser(userRequest)
		s.mockUserUsecase.On("Register", mock.MatchedBy(func(u *domain.User) bool {
			return u.Username == user.Username &&
				u.Email == user.Email &&
				u.FirstName == user.FirstName &&
				u.LastName == user.LastName &&
				u.Provider == user.Provider &&
				len(u.Password) > 0
		})).Return(nil)

		c, w := s.createTestRequest(http.MethodPost, "/register", userRequest, nil)
		s.handler.RegisterRequest(c)

		s.Equal(http.StatusCreated, w.Code)
		var response dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(userRequest.Username, response.Username)
		s.Equal(userRequest.Email, response.Email)
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/register", "{invalid json}", nil)
		s.handler.RegisterRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		userRequest := dto.UserRequest{
			Username:  "ab",      // Too short
			Email:     "invalid", // Invalid email
			Password:  "pass",    // Too short
			FirstName: "",
			LastName:  "",
		}
		c, w := s.createTestRequest(http.MethodPost, "/register", userRequest, nil)
		s.handler.RegisterRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "Field validation")
		s.resetMocks()
	})

	s.Run("RegisterError", func() {
		userRequest := dto.UserRequest{
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "password123",
			Provider:  "manual",
			FirstName: "Test",
			LastName:  "User",
		}
		user := dto.ToDomainUser(userRequest)
		s.mockUserUsecase.On("Register", &user).Return(errors.New("registration failed"))
		c, w := s.createTestRequest(http.MethodPost, "/register", userRequest, nil)
		s.handler.RegisterRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("registration failed", response["error"])
		s.resetMocks()
	})
}

// TestLoginRequest tests the LoginRequest method
func (s *AuthControllerSuite) TestLoginRequest() {
	s.Run("Success", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "test@example.com",
			Password:   "password123",
		}
		hashedPassword, _ := security.HashPassword("password123")
		user := &domain.User{
			ID:       "1",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}
		tokenResponse := domain.RefreshTokenResponse{
			AccessToken:           "access-token",
			RefreshToken:          "refresh-token",
			AccessTokenExpiresAt:  time.Now().Add(time.Hour),
			RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
		}
		s.mockUserUsecase.On("FindByUsernameOrEmail", mock.Anything, loginRequest.Identifier).Return(user, nil)
		s.mockAuthService.On("GenerateTokens", *user).Return(tokenResponse, nil)
		s.mockRefreshTokenUsecase.On("FindByUserID", user.ID).Return(nil, errors.New("refresh token not found"))
		s.mockRefreshTokenUsecase.On("Save", mock.Anything).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/login", loginRequest, nil)

		s.handler.LoginRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Println(response)
		s.Equal("Login successful", response["message"])
		userResponse := response["user"].(map[string]any)
		s.Equal(user.Email, userResponse["email"])
		tokens := response["tokens"].(map[string]any)
		s.Equal(tokenResponse.AccessToken, tokens["access_token"])
		s.Equal(tokenResponse.RefreshToken, tokens["refresh_token"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {

		c, w := s.createTestRequest(http.MethodPost, "/login", "{invalid json}", nil)
		s.handler.LoginRequest(c)
		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["message"], "invalid character")
		s.resetMocks()
	})

	s.Run("EmptyIdentifierOrPassword", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "",
			Password:   "",
		}
		c, w := s.createTestRequest(http.MethodPost, "/login", loginRequest, nil)
		s.handler.LoginRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Username/email and password are required", response["message"])
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "test@example.com",
			Password:   "password123",
		}
		s.mockUserUsecase.On("FindByUsernameOrEmail", mock.Anything, loginRequest.Identifier).Return(nil, errors.New("user not found"))
		c, w := s.createTestRequest(http.MethodPost, "/login", loginRequest, nil)

		s.handler.LoginRequest(c)
		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid email or password", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidPassword", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "test@example.com",
			Password:   "wrongpassword",
		}
		hashedPassword, _ := security.HashPassword("password123")
		user := &domain.User{
			ID:       "1",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}
		s.mockUserUsecase.On("FindByUsernameOrEmail", mock.Anything, loginRequest.Identifier).Return(user, nil)

		c, w := s.createTestRequest(http.MethodPost, "/login", loginRequest, nil)
		s.handler.LoginRequest(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid email or password", response["error"])
		s.resetMocks()
	})

	s.Run("TokenGenerationError", func() {
		loginRequest := dto.LoginRequest{
			Identifier: "test@example.com",
			Password:   "password123",
		}
		hashedPassword, _ := security.HashPassword("password123")
		user := &domain.User{
			ID:       "1",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}
		s.mockUserUsecase.On("FindByUsernameOrEmail", mock.Anything, loginRequest.Identifier).Return(user, nil)
		s.mockAuthService.On("GenerateTokens", *user).Return(domain.RefreshTokenResponse{}, errors.New("token generation failed"))

		c, w := s.createTestRequest(http.MethodPost, "/login", loginRequest, nil)

		s.handler.LoginRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to generate token", response["error"])
		s.resetMocks()
	})
}

// TestRefreshToken tests the RefreshToken method
func (s *AuthControllerSuite) TestRefreshToken() {
	s.Run("SuccessWithoutRotation", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "refresh-token",
		}
		tokenDoc := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(3 * time.Hour),
			Revoked:   false,
		}
		user := &domain.User{
			ID:    "1",
			Email: "test@example.com",
		}
		tokenResponse := &domain.RefreshTokenResponse{
			AccessToken:           "new-access-token",
			RefreshToken:          "refresh-token",
			AccessTokenExpiresAt:  time.Now().Add(time.Hour),
			RefreshTokenExpiresAt: tokenDoc.ExpiresAt,
		}
		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(tokenDoc, nil)
		s.mockAuthService.On("ValidateRefreshToken", refreshRequest.RefreshToken).Return(jwt.MapClaims{"user_id": "1"}, nil)
		s.mockUserUsecase.On("FindUserByID", tokenDoc.UserID).Return(user, nil)
		s.mockAuthService.On("GenerateTokens", *user).Return(*tokenResponse, nil)

		cookie := []*http.Cookie{{Name: "refresh_token", Value: "refresh-token"}}
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, cookie)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(tokenResponse.AccessToken, response.AccessToken)
		s.Equal(tokenResponse.RefreshToken, response.RefreshToken)
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/refresh", "{invalid json}", nil)
		s.handler.RefreshToken(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid payload", response["error"])
		s.resetMocks()
	})

	s.Run("ExpiredToken", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "expired-refresh-token",
		}

		tokenDoc := &domain.RefreshToken{
			Token:     "expired-refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(-time.Hour), // Token is ExpiredToken
			Revoked:   false,
		}
		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(tokenDoc, nil)
		s.mockRefreshTokenUsecase.On("DeleteByUserID", mock.Anything).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, nil)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains("Invalid or expired token", response["error"])
		s.resetMocks()
	})

	s.Run("NoTokenOnDB", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "refresh-token",
		}
		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(nil, errors.New("token not found"))
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, nil)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid or expired token", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidToken", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "invalid-token",
		}
		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(nil, errors.New("token not found"))
		s.mockRefreshTokenUsecase.On("DeleteByUserID", mock.Anything).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, nil)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid or expired token", response["error"])
		s.resetMocks()
	})

	s.Run("NoCookie", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "refresh-token",
		}
		tokenDoc := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(3 * time.Hour),
			Revoked:   false,
		}
		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(tokenDoc, nil)
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, nil)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("No refresh token found in cookies, please login again", response["error"])
		s.resetMocks()
	})

	s.Run("TokenValidationError", func() {
		refreshRequest := dto.RefreshTokenRequest{
			RefreshToken: "invalid-token",
		}
		tokenDoc := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(3 * time.Hour),
			Revoked:   false,
		}

		s.mockRefreshTokenUsecase.On("FindByToken", refreshRequest.RefreshToken).Return(tokenDoc, nil)
		s.mockAuthService.On("ValidateRefreshToken", refreshRequest.RefreshToken).Return(jwt.MapClaims{}, errors.New("invalid token"))
		cookie := []*http.Cookie{{Name: "refresh_token", Value: "refresh-token"}}
		c, w := s.createTestRequest(http.MethodPost, "/refresh", refreshRequest, cookie)

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains("Invalid or expired token provided", response["error"])
		s.resetMocks()
	})

}

// TestLogoutRequest tests the LogoutRequest method
func (s *AuthControllerSuite) TestLogoutRequest() {
	s.Run("Success", func() {
		tokenDoc := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRefreshTokenUsecase.On("FindByToken", "refresh-token").Return(tokenDoc, nil)
		s.mockRefreshTokenUsecase.On("RevokedToken", tokenDoc).Return(nil)
		s.mockRefreshTokenUsecase.On("DeleteByUserID", tokenDoc.UserID).Return(nil)

		cookie := []*http.Cookie{{Name: "refresh_token", Value: "refresh-token"}, {Name: "access_token", Value: "access-token"}}
		c, w := s.createTestRequest(http.MethodPost, "/logout", nil, cookie)

		s.handler.LogoutRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Logged out successfully", response["message"])
		s.resetMocks()
	})

	s.Run("NoCookie", func() {
		s.resetMocks()
		c, w := s.createTestRequest(http.MethodPost, "/logout", nil, nil)
		s.handler.LogoutRequest(c)
		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You are not logged in or your session has expired", response["error"])
		s.resetMocks()
	})

	s.Run("InvalidToken", func() {
		s.mockRefreshTokenUsecase.On("FindByToken", "invalid-token").Return(nil, errors.New("token not found"))

		cookie := []*http.Cookie{{Name: "refresh_token", Value: "invalid-token"}}
		c, w := s.createTestRequest(http.MethodPost, "/logout", nil, cookie)

		s.handler.LogoutRequest(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("You are not logged in or your session has expired", response["error"])
		s.resetMocks()
	})

	s.Run("RevokeError", func() {
		tokenDoc := &domain.RefreshToken{
			Token:     "refresh-token",
			UserID:    "1",
			ExpiresAt: time.Now().Add(time.Hour),
			Revoked:   false,
		}
		s.mockRefreshTokenUsecase.On("FindByToken", "refresh-token").Return(tokenDoc, nil)
		s.mockRefreshTokenUsecase.On("RevokedToken", tokenDoc).Return(errors.New("revoke failed"))

		cookie := []*http.Cookie{{Name: "refresh_token", Value: "refresh-token"}}
		c, w := s.createTestRequest(http.MethodPost, "/logout", nil, cookie)

		s.handler.LogoutRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to revoke token", response["error"])
		s.resetMocks()
	})
}

// TestChangeRoleRequest tests the ChangeRoleRequest method
func (s *AuthControllerSuite) TestChangeRoleRequest() {
	s.Run("Success", func() {
		changeRoleRequest := dto.ChangeRoleRequest{
			UserID: "1",
			Role:   "admin",
		}
		s.mockUserUsecase.On("ChangeRole", "admin", changeRoleRequest.UserID, mock.MatchedBy(func(u domain.User) bool {
			return u.Role == domain.UserRole(changeRoleRequest.Role)
		})).Return(nil)

		c, w := s.createTestRequest(http.MethodPost, "/change-role", changeRoleRequest, nil)
		c.Set("role", "admin")

		s.handler.ChangeRoleRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("User role changed successfully", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/change-role", "{invalid json}", nil)

		s.handler.ChangeRoleRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid request", response["error"])
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		changeRoleRequest := dto.ChangeRoleRequest{
			UserID: "",
			Role:   "invalid",
		}
		c, w := s.createTestRequest(http.MethodPost, "/change-role", changeRoleRequest, nil)

		s.handler.ChangeRoleRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Field validation")
		s.resetMocks()
	})

	s.Run("ChangeRoleError", func() {
		changeRoleRequest := dto.ChangeRoleRequest{
			UserID: "1",
			Role:   "admin",
		}
		s.mockUserUsecase.On("ChangeRole", "admin", changeRoleRequest.UserID, mock.MatchedBy(func(u domain.User) bool {
			return u.Role == domain.UserRole(changeRoleRequest.Role)
		})).Return(errors.New("change role failed"))

		c, w := s.createTestRequest(http.MethodPost, "/change-role", changeRoleRequest, nil)
		c.Set("role", "admin")

		s.handler.ChangeRoleRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to change user role", response["message"])
		s.resetMocks()
	})
}

// TestForgotPasswordRequest tests the ForgotPasswordRequest method
func (s *AuthControllerSuite) TestForgotPasswordRequest() {
	s.Run("Success", func() {
		forgotPasswordRequest := dto.ForgotPasswordRequest{
			Email: "test@example.com",
		}
		s.mockPasswordResetUsecase.On("RequestReset", forgotPasswordRequest.Email).Return(nil)

		c, w := s.createTestRequest(http.MethodPost, "/forgot-password", forgotPasswordRequest, nil)

		s.handler.ForgotPasswordRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Password reset link sent to email", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/forgot-password", "{invalid json}", nil)

		s.handler.ForgotPasswordRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid email", response["error"])
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		forgotPasswordRequest := dto.ForgotPasswordRequest{
			Email: "invalid",
		}
		c, w := s.createTestRequest(http.MethodPost, "/forgot-password", forgotPasswordRequest, nil)

		s.handler.ForgotPasswordRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Field validation")
		s.resetMocks()
	})

	s.Run("RequestResetError", func() {
		forgotPasswordRequest := dto.ForgotPasswordRequest{
			Email: "test@example.com",
		}
		s.mockPasswordResetUsecase.On("RequestReset", forgotPasswordRequest.Email).Return(errors.New("reset request failed"))
		c, w := s.createTestRequest(http.MethodPost, "/forgot-password", forgotPasswordRequest, nil)

		s.handler.ForgotPasswordRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to process request", response["message"])
		s.resetMocks()
	})
}

// TestResetPasswordRequest tests the ResetPasswordRequest method
func (s *AuthControllerSuite) TestResetPasswordRequest() {
	s.Run("Success", func() {
		resetPasswordRequest := dto.ResetPasswordRequest{
			Email:       "test@example.com",
			Token:       "reset-token",
			NewPassword: "newpassword123",
		}
		s.mockPasswordResetUsecase.On("ResetPassword", resetPasswordRequest.Email, resetPasswordRequest.Token, resetPasswordRequest.NewPassword).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/reset-password", resetPasswordRequest, nil)

		s.handler.ResetPasswordRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Password reset successfully", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/reset-password", "{invalid json}", nil)

		s.handler.ResetPasswordRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid request", response["error"])
		s.resetMocks()
	})

	s.Run("ResetPasswordError", func() {
		resetPasswordRequest := dto.ResetPasswordRequest{
			Email:       "test@example.com",
			Token:       "reset-token",
			NewPassword: "newpassword123",
		}
		s.mockPasswordResetUsecase.On("ResetPassword", resetPasswordRequest.Email, resetPasswordRequest.Token, resetPasswordRequest.NewPassword).Return(errors.New("reset failed"))

		c, w := s.createTestRequest(http.MethodPost, "/reset-password", resetPasswordRequest, nil)

		s.handler.ResetPasswordRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to reset password", response["message"])
		s.resetMocks()
	})
}

// TestVerifyEmailRequest tests the VerifyEmailRequest method
func (s *AuthControllerSuite) TestVerifyEmailRequest() {
	s.Run("Success", func() {
		verifyEmailRequest := dto.VerifyEmailRequest{
			Email: "test@example.com",
		}
		user := &domain.User{
			ID:    "1",
			Email: "test@example.com",
		}
		s.mockUserUsecase.On("GetUserByEmail", verifyEmailRequest.Email).Return(user, nil)
		s.mockOTPUsecase.On("RequestOTP", user.Email).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/verify-email", verifyEmailRequest, nil)

		s.handler.VerifyEmailRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Please check your email for the OTP to verify your account", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/verify-email", "{invalid json}", nil)

		s.handler.VerifyEmailRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid request", response["error"])
		s.resetMocks()
	})

	s.Run("ValidationError", func() {
		verifyEmailRequest := dto.VerifyEmailRequest{
			Email: "invalid",
		}
		c, w := s.createTestRequest(http.MethodPost, "/verify-email", verifyEmailRequest, nil)

		s.handler.VerifyEmailRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Field validation")
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		verifyEmailRequest := dto.VerifyEmailRequest{
			Email: "test@example.com",
		}
		s.mockUserUsecase.On("GetUserByEmail", verifyEmailRequest.Email).Return(nil, domain.ErrNotFound)
		c, w := s.createTestRequest(http.MethodPost, "/verify-email", verifyEmailRequest, nil)

		s.handler.VerifyEmailRequest(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(domain.ErrNotFound.Error(), response["error"])
		s.resetMocks()
	})
}

// TestVerifyOTPRequest tests the VerifyOTPRequest method
func (s *AuthControllerSuite) TestVerifyOTPRequest() {
	s.Run("Success", func() {
		verifyOTPRequest := dto.VerifyOTPRequest{
			Code: "123456",
		}
		user := &domain.User{
			ID:         "1",
			Email:      "test@example.com",
			IsVerified: false,
		}
		otp := &domain.OTP{
			ID:       "otp1",
			Email:    "test@example.com",
			CodeHash: "123456",
		}
		s.mockUserUsecase.On("FindUserByID", "1").Return(user, nil)
		s.mockOTPUsecase.On("VerifyOTP", user.Email, verifyOTPRequest.Code).Return(otp, nil)
		s.mockUserUsecase.On("UpdateUser", user.ID, mock.MatchedBy(func(u *domain.User) bool {
			return u.IsVerified
		})).Return(user, nil)
		s.mockOTPUsecase.On("DeleteByID", otp.ID).Return(nil)
		c, w := s.createTestRequest(http.MethodPost, "/verify-otp", verifyOTPRequest, nil)
		c.Set("user_id", "1")

		s.handler.VerifyOTPRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Email verified successfully", response["message"])
		s.resetMocks()
	})

	s.Run("InvalidJSON", func() {
		c, w := s.createTestRequest(http.MethodPost, "/verify-otp", "{invalid_json}", nil)
		c.Set("user_id", "1")

		s.handler.VerifyOTPRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Invalid request", response["error"])
		s.resetMocks()
	})

	s.Run("BindingError", func() {
		verifyOTPRequest := dto.VerifyOTPRequest{
			Code: "",
		}
		c, w := s.createTestRequest(http.MethodPost, "/verify-otp", verifyOTPRequest, nil)
		c.Set("user_id", "1")

		s.handler.VerifyOTPRequest(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "Invalid request")
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		verifyOTPRequest := dto.VerifyOTPRequest{
			Code: "123456",
		}
		s.mockUserUsecase.On("FindUserByID", "1").Return(nil, domain.ErrNotFound)
		c, w := s.createTestRequest(http.MethodPost, "/verify-otp", verifyOTPRequest, nil)
		c.Set("user_id", "1")

		s.handler.VerifyOTPRequest(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Println(response)
		s.Equal(domain.ErrNotFound.Error(), response["error"])
		s.resetMocks()
	})

	s.Run("InvalidOTP", func() {
		verifyOTPRequest := dto.VerifyOTPRequest{
			Code: "123456",
		}
		user := &domain.User{
			ID:    "1",
			Email: "test@example.com",
		}
		s.mockUserUsecase.On("FindUserByID", "1").Return(user, nil)
		s.mockOTPUsecase.On("VerifyOTP", user.Email, verifyOTPRequest.Code).Return(nil, errors.New("invalid OTP"))

		c, w := s.createTestRequest(http.MethodPost, "/verify-otp", verifyOTPRequest, nil)
		c.Set("user_id", "1")

		s.handler.VerifyOTPRequest(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("invalid OTP", response["error"])
		s.resetMocks()
	})
}

// TestResendOTPRequest tests the ResendOTPRequest method
func (s *AuthControllerSuite) TestResendOTPRequest() {
	s.Run("Success", func() {
		user := &domain.User{
			ID:    "1",
			Email: "test@example.com",
		}
		s.mockUserUsecase.On("FindUserByID", "1").Return(user, nil)
		s.mockOTPUsecase.On("RequestOTP", user.Email).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/resend-otp", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "1")

		s.handler.ResendOTPRequest(c)

		s.Equal(http.StatusOK, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("OTP resent successfully", response["message"])
		s.resetMocks()
	})

	s.Run("UserNotFound", func() {
		s.mockUserUsecase.On("FindUserByID", "1").Return(nil, domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodPost, "/resend-otp", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "1")

		s.handler.ResendOTPRequest(c)

		s.Equal(http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(domain.ErrNotFound.Error(), response["error"])
		s.resetMocks()
	})

	s.Run("RequestOTPError", func() {
		user := &domain.User{
			ID:    "1",
			Email: "test@example.com",
		}
		s.mockUserUsecase.On("FindUserByID", "1").Return(user, nil)
		s.mockOTPUsecase.On("RequestOTP", user.Email).Return(errors.New("OTP request failed"))

		req := httptest.NewRequest(http.MethodPost, "/resend-otp", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", "1")

		s.handler.ResendOTPRequest(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("Failed to resend OTP", response["message"])
		s.resetMocks()
	})
}

func (s *AuthControllerSuite) resetMocks() {
	s.mockUserUsecase.ExpectedCalls = nil
	s.mockUserUsecase.Calls = nil
	s.mockAuthService.ExpectedCalls = nil
	s.mockAuthService.Calls = nil
	s.mockOTPUsecase.ExpectedCalls = nil
	s.mockOTPUsecase.Calls = nil
	s.mockRefreshTokenUsecase.ExpectedCalls = nil
	s.mockRefreshTokenUsecase.Calls = nil
	s.mockPasswordResetUsecase.ExpectedCalls = nil
	s.mockPasswordResetUsecase.Calls = nil
}

func (s *AuthControllerSuite) createTestRequest(method, url string, body interface{}, cookies []*http.Cookie) (*gin.Context, *httptest.ResponseRecorder) {
	var requestBody *bytes.Reader
	switch b := body.(type) {
	case string:
		requestBody = bytes.NewReader([]byte(b))
	case nil:
		requestBody = bytes.NewReader(nil)
	default:
		jsonBody, _ := json.Marshal(b)
		requestBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}
