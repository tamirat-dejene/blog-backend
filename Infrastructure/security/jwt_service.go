package security

import (
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtService struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessExpiry, refreshExpiry int) domain.IRefreshTokenRepository {
	return &JwtService{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		AccessExpiry:  time.Duration(accessExpiry) * time.Hour,
		RefreshExpiry: time.Duration(refreshExpiry) * time.Hour,
	}
}

func (s *JwtService) GenerateTokens(user domain.User) (domain.RefreshToken, error) {
	fmt.Println(s, user)
	accessClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(s.AccessExpiry).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.AccessSecret))
	if err != nil {
		return domain.RefreshToken{}, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.RefreshExpiry).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(s.RefreshSecret))
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return domain.RefreshToken{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *JwtService) ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.AccessSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *JwtService) ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.RefreshSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
