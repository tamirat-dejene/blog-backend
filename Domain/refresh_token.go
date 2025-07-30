package Domain

import "github.com/golang-jwt/jwt/v4"

type RefreshToken struct {
	AccessToken  string
	RefreshToken string
}

type IRefreshTokenUsecase interface {
	GenerateTokens(user User) (RefreshToken, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(token string) (jwt.MapClaims, error)
	GetByUsername(string) (*User, error)
}

type IRefreshTokenRepository interface {
	GenerateTokens(user User) (RefreshToken, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ValidateRefreshToken(token string) (jwt.MapClaims, error)
}
