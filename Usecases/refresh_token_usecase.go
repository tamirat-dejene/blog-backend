package usecases

import (
	"context"
	"g6/blog-api/Domain"
	"time"
)

type RefreshTokenUsecase struct {
	Repo Domain.IRefreshTokenRepository
}

func NewRefreshTokenUsecase(repo Domain.IRefreshTokenRepository) Domain.IRefreshTokenUsecase {
	return &RefreshTokenUsecase{
		Repo: repo,
	}
}

func (uc *RefreshTokenUsecase) FindByToken(token string) (*Domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.FindByToken(ctx, token)
}

func (uc *RefreshTokenUsecase) Save(token *Domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.Save(ctx, token)
}
func (uc *RefreshTokenUsecase) DeleteByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return uc.Repo.DeleteByUserID(ctx, userID)
}
