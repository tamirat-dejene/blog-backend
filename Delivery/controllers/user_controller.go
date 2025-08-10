package controllers

import (
	"io"
	"net/http"

	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	uc domain.IUserUsecase
}

func NewUserController(uc domain.IUserUsecase) *UserController {
	return &UserController{uc: uc}
}

func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized})
		return
	}

	var req dto.UserUpdateProfileRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "details": err.Error()})
		return
	}

	// validate
	if err := validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var avatarData []byte
	var fileName string

	file, err := c.FormFile("avatar_file")
	if err == nil {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer f.Close()

		avatarData, err = io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read uploaded file"})
			return
		}
		fileName = file.Filename
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error uploading file"})
		return
	} else {
		// No file uploaded, use default avatar or leave it empty
		avatarData = nil
		fileName = ""
	}

	update := domain.UserProfileUpdate{
		Bio:        req.Bio,
		AvatarData: avatarData,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
	}

	updatedUser, err := ctrl.uc.UpdateProfile(userID.(string), update, fileName)
	if err == domain.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err == domain.ErrInvalidFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(*updatedUser))
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := ctrl.uc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
