package DTO

import domain "g6/blog-api/Domain"

// change to refresh token to domain model
func (r *RefreshTokenRequest) FromRequestToDomainRefreshToken() *domain.RefreshToken {
	return &domain.RefreshToken{
		RefreshToken: r.RefreshToken,
	}
}

func FromDomainRefreshTokenToResponse(refreshToken *domain.RefreshToken) *domain.RefreshToken {
	if refreshToken == nil {
		return nil
	}
	return &domain.RefreshToken{
		AccessToken:  refreshToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}
}
