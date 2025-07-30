package database

import (
	"g6/blog-api/Domain"
	"time"
)

type RefreshTokenDB struct {
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

func FromRefreshTokenEntityToDB(token *Domain.RefreshToken) *RefreshTokenDB {
	return &RefreshTokenDB{
		Token:     token.Token,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: time.Now(),
	}
}

func FromRefreshTokenDBToEntity(tokenDB *RefreshTokenDB) *Domain.RefreshToken {
	return &Domain.RefreshToken{
		Token:     tokenDB.Token,
		UserID:    tokenDB.UserID,
		ExpiresAt: tokenDB.ExpiresAt,
		CreatedAt: tokenDB.CreatedAt,
	}
}
