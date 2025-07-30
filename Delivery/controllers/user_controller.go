package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/DTO"
	"g6/blog-api/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type UserController struct {
	UserUsecase Domain.IUserUsecase
}

func NewUserController(userUsecase Domain.IUserUsecase) *UserController {
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
	c.JSON(http.StatusOK, DTO.FromUserEntityToDTOList(users))
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user DTO.UserRequest
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
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": DTO.FromUserEntityToDTO(newUser)})
}
