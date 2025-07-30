package Domain

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
	Role      string
	Bio       string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IUserUsecase interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	// UpdateUser(id primitive.ObjectID, user *User) (*User, error)
	// DeleteUser(id primitive.ObjectID) error
	GetAllUsers() ([]*User, error)
}

type IUserRepository interface {
	CreateUser(context.Context, *User) error
	GetUserByUsername(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	// UpdateUser(id primitive.ObjectID, user *User) error
	// DeleteUser(id primitive.ObjectID) error
	GetAllUsers(context.Context) ([]*User, error)
	// FindUserByUsername(username string) (*User, error)
	// FindUserByEmail(email string) (*User, error)
	// FindUserByID(id primitive.ObjectID) (*User, error)
	// FindUserByRole(role string) ([]*User, error)
}
