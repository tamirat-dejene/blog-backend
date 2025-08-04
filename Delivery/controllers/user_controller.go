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

func (ctrl *UserController) Register(ctx *gin.Context) {
	var req dto.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := dto.ToDomainUser(req)
	if err := ctrl.uc.Register(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (ctrl *UserController) Logout(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if err := ctrl.uc.Logout(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (ctrl *UserController) ChangeRole(ctx *gin.Context) {
	var req dto.ChangeRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	initiator := ctx.GetString("role")
	target := ctx.Param("userID")
	user := domain.User{
		Role: domain.UserRole(req.Role),
	}
	if err := ctrl.uc.ChangeRole(initiator, target, user); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User role changed successfully"})
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
