package DTO

import (
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
