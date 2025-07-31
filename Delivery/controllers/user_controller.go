package controllers

import (
	"net/http"

	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	

	"github.com/gin-gonic/gin"
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

	user := dto.ToDomainUser(req)
	if err := ctrl.uc.Register(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
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
