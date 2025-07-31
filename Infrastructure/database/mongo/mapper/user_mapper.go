package mapper

import (
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Role      string             `bson:"role"`
	Bio       string             `bson:"bio"`
	AvatarURL string             `bson:"avatar_url"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func UserToDomain(user *UserModel) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Role:      domain.UserRole(user.Role),
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserFromDomain(user *domain.User) *UserModel {

	return &UserModel{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Role:      string(user.Role),
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToDomainList(userModels []*UserModel) []*domain.User {
	var users []*domain.User
	for _, userModel := range userModels {
		users = append(users, UserToDomain(userModel))
	}
	return users
}
