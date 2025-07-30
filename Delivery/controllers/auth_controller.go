package controllers

import (
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/security"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	existedUser, _ := ac.UserUsecase.GetUserByUsername(newUser.Username)
	if existedUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	existedUser, _ = ac.UserUsecase.GetUserByEmail(newUser.Email)
	if existedUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := newUser.ToUserEntity()
	user.Password = string(hashPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = ac.UserUsecase.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, dto.FromUserEntityToDTO(user))
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
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, loginRequest.Identifier)

	var user *domain.User
	var err error

	if isEmail {
		user, err = ac.UserUsecase.GetUserByEmail(loginRequest.Identifier)
	} else {
		user, err = ac.UserUsecase.GetUserByUsername(loginRequest.Identifier)
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	match := security.ValidatePassword(user.Password, loginRequest.Password)
	if !match {
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
