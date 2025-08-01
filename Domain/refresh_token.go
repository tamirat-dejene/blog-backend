package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PasswordResetToken struct {
	ID        string
	Email     string
	TokenHash string
	RateLimit int
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type IAuthService interface {
	GenerateTokens(user User) (RefreshTokenResponse, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(token string) (jwt.MapClaims, error)
}

type IRefreshTokenUsecase interface {
	FindByToken(token string) (*RefreshToken, error)
	Save(token *RefreshToken) error
	DeleteByUserID(userID string) error
	ReplaceToken(token *RefreshToken) error
	RevokedToken(token *RefreshToken) error
	FindByUserID(id string) (*RefreshToken, error)
}

type IRefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID string) error
	ReplaceTokenByUserID(ctx context.Context, token *RefreshToken) error
	RevokeToken(ctx context.Context, token string) error
	FindTokenByUserID(ctx context.Context, token string) (*RefreshToken, error)
}
