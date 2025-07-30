package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	UserUsecase domain.IUserUsecase
}

func NewUserController(userUsecase domain.IUserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.UserUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, dto.FromUserEntityToDTOList(users))
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user dto.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}
	newUser := user.ToUserEntity()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	if err := uc.UserUsecase.CreateUser(newUser); err != nil {
		fmt.Println("Error creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": dto.FromUserEntityToDTO(newUser)})
}
