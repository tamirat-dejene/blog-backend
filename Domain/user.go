package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID
	Username  string
	Email     string
	FirstName string
	LastName  string
	Password  string
	Role      UserRole
	Bio       string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleUser       UserRole = "user"
	RoleSuperAdmin UserRole = "superadmin"
)

type IUserUsecase interface {
	FindByUsernameOrEmail(context.Context, string) (*User, error)
	FindUserByID(string) (*User, error)
	// GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	// UpdateUser(id primitive.ObjectID, user *User) (*User, error)
	// DeleteUser(id primitive.ObjectID) error
	// GetAllUsers() ([]*User, error)

	// anti
	Register(request *User) error
	ChangeRole(initiatorRole string, targetUserID string, request User) error
	Logout(userID string) error
}

type IUserRepository interface {
	CreateUser(context.Context, *User) error
	FindUserByID(context.Context, string) (*User, error)
	GetUserByUsername(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	// UpdateUser(id primitive.ObjectID, user *User) error
	// DeleteUser(id primitive.ObjectID) error
	GetAllUsers(context.Context) ([]*User, error)
	// FindUserByUsername(username string) (*User, error)
	// FindUserByEmail(email string) (*User, error)
	// FindUserByID(id primitive.ObjectID) (*User, error)
	// FindUserByRole(role string) ([]*User, error)

	// anti
	FindByUsernameOrEmail(context.Context, string) (User, error)
	InvalidateTokens(context.Context, string) error
	ChangeRole(context.Context, string, string, string) error
}
