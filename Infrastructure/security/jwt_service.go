package security

import (
	domain "g6/blog-api/Domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  string
	refreshSecret string
	accessExp     time.Duration
	refreshExp    time.Duration
)

func InitJWT(ats, rts string, atMinutes, rtHours int) {
	accessSecret = ats
	refreshSecret = rts
	accessExp = time.Duration(atMinutes) * time.Minute
	refreshExp = time.Duration(rtHours) * time.Hour
}

func GenerateTokenPair(userID, role string) (domain.TokenPair, error) {
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(accessExp).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(refreshExp).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return domain.TokenPair{}, err
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func GetAccessSecret() string {
	return accessSecret
}
func GetRefreshSecret() string {
	return refreshSecret
}
