package controllers

import (
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserUsecase         domain.IUserUsecase
	AuthService         domain.IAuthService
	RefreshTokenUsecase domain.IRefreshTokenUsecase
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

	user, err := ac.UserUsecase.FindByUsernameOrEmail(c.Request.Context(), loginRequest.Identifier)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = security.ValidatePassword(user.Password, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	response, err := ac.AuthService.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save the refresh token in the database
	refreshToken := &domain.RefreshToken{
		Token:     response.RefreshToken,
		UserID:    user.ID.Hex(),
		ExpiresAt: response.ExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := ac.RefreshTokenUsecase.Save(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (ac *AuthController) LogoutRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uid, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	if err := ac.RefreshTokenUsecase.DeleteByUserID(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}