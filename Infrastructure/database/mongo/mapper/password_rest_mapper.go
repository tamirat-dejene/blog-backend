package mapper

import (
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordResetTokenDB struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	RateLimit int64              `bson:"rate_limit"`
	TokenHash string             `bson:"token_hash"`
	ExpiresAt time.Time          `bson:"expires_at"`
	Used      bool               `bson:"used"`
	CreatedAt time.Time          `bson:"created_at"`
}

func PasswordResetTokenFromDomain(token *domain.PasswordResetToken) *PasswordResetTokenDB {
	return &PasswordResetTokenDB{
		ID:        primitive.NewObjectID(),
		Email:     token.Email,
		RateLimit: int64(token.RateLimit),
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
	}
}

func PasswordResetTokenToDomain(token *PasswordResetTokenDB) *domain.PasswordResetToken {
	return &domain.PasswordResetToken{
		ID:        token.ID.Hex(),
		Email:     token.Email,
		RateLimit: int(token.RateLimit),
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
		CreatedAt: token.CreatedAt,
	}
}
