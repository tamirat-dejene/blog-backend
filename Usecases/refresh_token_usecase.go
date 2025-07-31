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
