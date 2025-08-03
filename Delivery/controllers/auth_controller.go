package controllers

import (
	"fmt"
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	utils "g6/blog-api/Utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUsecase          domain.IUserUsecase
	AuthService          domain.IAuthService
	RefreshTokenUsecase  domain.IRefreshTokenUsecase
	PasswordResetUsecase domain.IPasswordResetUsecase
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
	// check if user is really exist
	user, err := ac.UserUsecase.FindByUsernameOrEmail(c.Request.Context(), loginRequest.Identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// validate the password
	err = security.ValidatePassword(user.Password, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// if every thing is fine generate access and refresh token
	response, err := ac.AuthService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save the refresh token in the database
	refreshToken := &domain.RefreshToken{
		Token:     response.RefreshToken,
		UserID:    user.ID,
		Revoked:   false,
		ExpiresAt: response.ExpiresAt,
		CreatedAt: time.Now(),
	}

	// Check if the user have token on db
	existingToken, _ := ac.RefreshTokenUsecase.FindByUserID(user.ID)
	if existingToken != nil {
		// refresh token shouldn't used twice so we have to make revoke true
		err := ac.RefreshTokenUsecase.RevokedToken(existingToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke existing refresh token"})
			return
		}
		// if that so just replace it instead of creating new one
		err = ac.RefreshTokenUsecase.ReplaceToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing refresh token"})
			return
		}
	} else {
		// Save the refresh token using the usecase
		if err := ac.RefreshTokenUsecase.Save(refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
			return
		}
	}

	// Set the refresh token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   int(time.Until(response.ExpiresAt).Seconds()),
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
		MaxAge:   int(time.Until(response.ExpiresAt).Seconds()),
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

// Refresh token
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
	// revoked the token on db
	if err := ac.RefreshTokenUsecase.RevokedToken(tokenDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
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

	// generate new token
	response, err := ac.AuthService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	// Save the refresh token in the database
	refreshToken := &domain.RefreshToken{
		Token:     response.RefreshToken,
		UserID:    user.ID,
		Revoked:   false,
		ExpiresAt: response.ExpiresAt,
		CreatedAt: time.Now(),
	}

	// Set the refresh token in the cookies
	utils.SetCookie(c, utils.CookieOptions{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   int(time.Until(response.ExpiresAt).Seconds()),
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
		MaxAge:   int(time.Until(response.ExpiresAt).Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Optional update the existing refresh token
	if err := ac.RefreshTokenUsecase.ReplaceToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
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
	fmt.Println("Part 1")

	tokenDoc, err := ac.RefreshTokenUsecase.FindByToken(refreshToken)
	if err != nil || tokenDoc == nil || tokenDoc.Revoked || time.Now().After(tokenDoc.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in or your session has expired"})
		return
	}
	fmt.Println("Part 2")
	// revoke the token
	if err := ac.RefreshTokenUsecase.RevokedToken(tokenDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}
	fmt.Println("Part 3")

	// delete the token from the database
	if err := ac.RefreshTokenUsecase.DeleteByUserID(tokenDoc.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
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
