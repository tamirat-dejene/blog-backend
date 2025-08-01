package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type RefreshTokenUsecase struct {
	Repo domain.IRefreshTokenRepository
}

func NewRefreshTokenUsecase(repo domain.IRefreshTokenRepository) domain.IRefreshTokenUsecase {
	return &RefreshTokenUsecase{
		Repo: repo,
	}
}

func (uc *RefreshTokenUsecase) FindByToken(token string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.FindByToken(ctx, token)
}

func (uc *RefreshTokenUsecase) Save(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.Save(ctx, token)
}
func (uc *RefreshTokenUsecase) DeleteByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.DeleteByUserID(ctx, userID)
}

func (uc *RefreshTokenUsecase) ReplaceToken(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.ReplaceTokenByUserID(ctx, token)
}

// revoke token
func (uc *RefreshTokenUsecase) RevokedToken(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.RevokeToken(ctx, token.Token)
}

// find token by user id
func (uc *RefreshTokenUsecase) FindByUserID(userID string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.FindTokenByUserID(ctx, userID)
}
