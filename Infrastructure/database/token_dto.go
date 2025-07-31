package database

import (
	domain "g6/blog-api/Domain"
	"time"
)

type RefreshTokenDB struct {
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

func FromRefreshTokenEntityToDB(token *domain.RefreshToken) *RefreshTokenDB {
	return &RefreshTokenDB{
		Token:     token.Token,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: time.Now(),
	}
}

func FromRefreshTokenDBToEntity(tokenDB *RefreshTokenDB) *domain.RefreshToken {
	return &domain.RefreshToken{
		Token:     tokenDB.Token,
		UserID:    tokenDB.UserID,
		ExpiresAt: tokenDB.ExpiresAt,
		CreatedAt: tokenDB.CreatedAt,
	}
}
