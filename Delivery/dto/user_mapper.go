package dto

import (
	domain "g6/blog-api/Domain"
	"time"
)

type UserRequest struct {
	ID        string    `json:"id" validate:"omitempty"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6,max=100"`
	FirstName string    `json:"first_name" validate:"required,alpha,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"required,alpha,min=2,max=50"`
	Role      string    `json:"role" validate:"required,oneof=admin user superadmin"`
	Bio       string    `json:"bio" validate:"max=500"`
	AvatarURL string    `json:"avatar_url" validate:"omitempty,url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// user registration request mapper
func ToDomainUser(req UserRequest) domain.User {
	return domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // Ensure to hash the password before saving to the domain
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.UserRole(req.Role),
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
	}
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// to response list
func ToUserResponseList(users []*domain.User) []UserResponse {
	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, ToUserResponse(*user))
	}
	return responses
}

// / 	USER UPDATE REQUEST
// user update profile request
type UserUpdateProfileRequest struct {
	Bio       string `form:"bio" validate:"omitempty,max=500"`
	FirstName string `form:"first_name" validate:"omitempty,alpha,min=2,max=50"`
	LastName  string `form:"last_name" validate:"omitempty,alpha,min=2,max=50"`
}
