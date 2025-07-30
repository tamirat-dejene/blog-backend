package dto

import (
	domain "g6/blog-api/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ur *UserRequest) ToUserEntity() *domain.User {
	return &domain.User{
		ID:        primitive.NewObjectID(),
		Username:  ur.Username,
		Email:     ur.Email,
		Password:  ur.Password, // Note: Password should be hashed before storing
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		Role:      ur.Role,
		Bio:       ur.Bio,
		AvatarURL: ur.AvatarURL,
	}
}

func FromUserEntityToDTO(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FromUserEntityToDTOList(users []*domain.User) []*UserResponse {
	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, FromUserEntityToDTO(user))
	}
	return userResponses
}
