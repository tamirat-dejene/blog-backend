package controllers

import (
	"context"
	"g6/blog-api/Delivery/bootstrap"
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/oauth"
	"g6/blog-api/Infrastructure/security"
	utils "g6/blog-api/Utils"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/oauth2/v2"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUsecase          domain.IUserUsecase
	AuthService          domain.IAuthService
	OTP                  domain.IOTPUsecase
	RefreshTokenUsecase  domain.IRefreshTokenUsecase
	PasswordResetUsecase domain.IPasswordResetUsecase
	Env                  *bootstrap.Env
}

func (ac *AuthController) RegisterRequest(c *gin.Context) {
	var newUser dto.UserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := dto.ToDomainUser(newUser)
	err := ac.UserUsecase.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (ac *AuthController) LoginRequest(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if loginRequest.Identifier == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username/email and password are required"})
		return
	}

	// Check if user exists
	user, err := ac.UserUsecase.FindByUsernameOrEmail(c.Request.Context(), loginRequest.Identifier)
	if err != nil {
		// Use generic error message for both user not found and password mismatch
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Validate the password
	if err := security.ValidatePassword(user.Password, loginRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate access and refresh tokens
	response, err := ac.AuthService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Prepare refresh token for DB
	refreshToken := &domain.RefreshToken{
		Token:     response.RefreshToken,
		UserID:    user.ID,
		Revoked:   false,
		ExpiresAt: response.RefreshTokenExpiresAt,
		CreatedAt: time.Now(),
	}

	// Check if the user already has a refresh token in DB
	existingToken, findErr := ac.RefreshTokenUsecase.FindByUserID(user.ID)
	if findErr != nil && findErr.Error() != "refresh token not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing refresh token"})
		return
	}
	if findErr != nil && findErr.Error() == "refresh token not found" {
		existingToken = nil // Explicitly set to nil if no token is found
	}

	if existingToken != nil {
		// Revoke the old token
		if err := ac.RefreshTokenUsecase.RevokedToken(existingToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke existing refresh token"})
			return
		}
		// Replace with the new token
		if err := ac.RefreshTokenUsecase.ReplaceToken(refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing refresh token"})
			return
		}
	} else {
		// Save the new refresh token
		if err := ac.RefreshTokenUsecase.Save(refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
			return
		}
	}

	// Set the refresh token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   int(time.Until(response.RefreshTokenExpiresAt).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	// Set the access token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "access_token",
		Value:    response.AccessToken,
		MaxAge:   int(time.Until(response.AccessTokenExpiresAt).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// find token from db
	tokenDoc, err := ac.RefreshTokenUsecase.FindByToken(req.RefreshToken)
	if err != nil || tokenDoc == nil || tokenDoc.Revoked || time.Now().After(tokenDoc.ExpiresAt) {
		if tokenDoc != nil {
			_ = ac.RefreshTokenUsecase.DeleteByUserID(tokenDoc.UserID)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// check both the token is valid and not expired
	_, err = utils.GetCookie(c, "refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token found in cookies, please login again"})
		return
	}

	// Validate the refresh token
	_, err = ac.AuthService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token provided"})
		return
	}

	// find the user of the token
	user, err := ac.UserUsecase.FindUserByID(tokenDoc.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// generate new access token
	response, err := ac.AuthService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Decide whether to rotate the refresh token
	rotateThreshold := 2 * time.Hour // 2 hours before expiry
	shouldRotate := time.Until(tokenDoc.ExpiresAt) < rotateThreshold

	var refreshTokenValue string
	var refreshTokenExpiry time.Time

	if shouldRotate {
		// Rotate: generate and store a new refresh token
		refreshToken := &domain.RefreshToken{
			Token:     response.RefreshToken,
			UserID:    user.ID,
			Revoked:   false,
			ExpiresAt: response.RefreshTokenExpiresAt,
			CreatedAt: time.Now(),
		}
		if err := ac.RefreshTokenUsecase.RevokedToken(tokenDoc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old refresh token"})
			return

		}
		if err := ac.RefreshTokenUsecase.ReplaceToken(refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
			return
		}
		refreshTokenValue = response.RefreshToken
		refreshTokenExpiry = response.RefreshTokenExpiresAt
	} else {
		// Do not rotate: keep the old refresh token
		refreshTokenValue = tokenDoc.Token
		refreshTokenExpiry = tokenDoc.ExpiresAt
	}

	// Set the refresh token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "refresh_token",
		Value:    refreshTokenValue,
		MaxAge:   int(time.Until(refreshTokenExpiry).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	// Set the access token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "access_token",
		Value:    response.AccessToken,
		MaxAge:   int(time.Until(response.AccessTokenExpiresAt).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: refreshTokenValue,
	})
}

// log out here
func (ac *AuthController) LogoutRequest(c *gin.Context) {

	// get the refresh token from cookies
	refreshToken, err := utils.GetCookie(c, "refresh_token")
	utils.DeleteCookie(c, "refresh_token")
	utils.DeleteCookie(c, "access_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in or your session has expired"})
		return
	}

	tokenDoc, err := ac.RefreshTokenUsecase.FindByToken(refreshToken)
	if err != nil || tokenDoc == nil || tokenDoc.Revoked || time.Now().After(tokenDoc.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in or your session has expired"})
		return
	}
	// revoke the token
	if err := ac.RefreshTokenUsecase.RevokedToken(tokenDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	// delete the token from the database
	if err := ac.RefreshTokenUsecase.DeleteByUserID(tokenDoc.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) ChangeRoleRequest(c *gin.Context) {
	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initiator := c.GetString("role")
	err := ac.UserUsecase.ChangeRole(initiator, req.UserID, domain.User{
		Role: domain.UserRole(req.Role),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to change user role", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role changed successfully"})
}

// forget password here
func (ac *AuthController) ForgotPasswordRequest(c *gin.Context) {

	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.PasswordResetUsecase.RequestReset(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process request", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to email"})
}

// Reset password here
func (ac *AuthController) ResetPasswordRequest(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.PasswordResetUsecase.ResetPassword(req.Email, req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reset password", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// verify email request
func (ac *AuthController) VerifyEmailRequest(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.UserUsecase.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	}

	err = ac.OTP.RequestOTP(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Please check your email for the OTP to verify your account"})
}

// verify otp
func (ac *AuthController) VerifyOTPRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.UserUsecase.FindUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	}

	otp, err := ac.OTP.VerifyOTP(user.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// update user verification status
	user.IsVerified = true
	if _, err := ac.UserUsecase.UpdateUser(user.ID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user verification status"})
		return
	}

	// delete the OTP after successful verification
	if err := ac.OTP.DeleteByID(otp.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrOTPFailedToDelete})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// otp resend request
func (ac *AuthController) ResendOTPRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := ac.UserUsecase.FindUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	}

	err = ac.OTP.RequestOTP(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to resend OTP", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}

//Oauth Google handlers

func (ac *AuthController) GoogleLogin(c *gin.Context) {
	conf := oauth.GetGoogleOAuthConfig(ac.Env.GoogleClientID, ac.Env.GoogleClientSecret, ac.Env.GoogleRedirectURL)
	url := conf.AuthCodeURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ac *AuthController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	conf := oauth.GetGoogleOAuthConfig(ac.Env.GoogleClientID, ac.Env.GoogleClientSecret, ac.Env.GoogleRedirectURL)
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	client := conf.Client(context.Background(), token)
	service, err := oauth2.New(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create oauth service"})
		return
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// Try finding user by email
	user, err := ac.UserUsecase.FindByUsernameOrEmail(c.Request.Context(), userInfo.Email)
	if err == nil && user != nil {
		// User exists - login flow
		response, err := ac.AuthService.GenerateTokens(*user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		refreshToken := &domain.RefreshToken{
			Token:     response.RefreshToken,
			UserID:    user.ID,
			ExpiresAt: response.RefreshTokenExpiresAt,
			Revoked:   false,
			CreatedAt: time.Now(),
		}

		if existingToken, err := ac.RefreshTokenUsecase.FindByUserID(user.ID); err == nil {
			ac.RefreshTokenUsecase.RevokedToken(existingToken)
			ac.RefreshTokenUsecase.ReplaceToken(refreshToken)
		} else {
			ac.RefreshTokenUsecase.Save(refreshToken)
		}

		utils.SetCookie(c, utils.CookieOptions{
			Name:     "refresh_token",
			Value:    response.RefreshToken,
			MaxAge:   int(time.Until(response.RefreshTokenExpiresAt).Seconds()),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})

		utils.SetCookie(c, utils.CookieOptions{
			Name:     "access_token",
			Value:    response.AccessToken,
			MaxAge:   int(time.Until(response.AccessTokenExpiresAt).Seconds()),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Login successful via Google", "user": dto.ToUserResponse(*user)})
		return
	}

	// If user not found – register new one
	newUser := &domain.User{
		Username:  strings.Split(userInfo.Email, "@")[0],
		Email:     userInfo.Email,
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Role:      "user",
		Bio:       "",
		AvatarURL: userInfo.Picture,
		Provider:  "google",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ac.UserUsecase.Register(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	// Generate token for new user
	response, err := ac.AuthService.GenerateTokens(*newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken := &domain.RefreshToken{
		Token:     response.RefreshToken,
		UserID:    newUser.ID,
		ExpiresAt: response.RefreshTokenExpiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if existingToken, err := ac.RefreshTokenUsecase.FindByUserID(newUser.ID); err == nil {
		ac.RefreshTokenUsecase.RevokedToken(existingToken)
		ac.RefreshTokenUsecase.ReplaceToken(refreshToken)
	} else {
		ac.RefreshTokenUsecase.Save(refreshToken)
	}

	utils.SetCookie(c, utils.CookieOptions{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   int(time.Until(response.RefreshTokenExpiresAt).Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	utils.SetCookie(c, utils.CookieOptions{
		Name:     "access_token",
		Value:    response.AccessToken,
		MaxAge:   int(time.Until(response.AccessTokenExpiresAt).Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully via Google", "user": dto.ToUserResponse(*newUser)})
}
