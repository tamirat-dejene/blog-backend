package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogUserReactionUsecase struct {
	blogUserReactionRepo domain.BlogUserReactionRepository
	ctxtimeout           time.Duration
}

func (b *blogUserReactionUsecase) CreateReaction(ctx context.Context, reaction *domain.BlogUserReaction) (*domain.BlogUserReaction, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogUserReactionRepo.Create(c, reaction)
}

func (b *blogUserReactionUsecase) DeleteReaction(ctx context.Context, id string) *domain.DomainError {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogUserReactionRepo.Delete(c, id)
}

func (b *blogUserReactionUsecase) GetUserReaction(ctx context.Context, blogID string, userID string) (*domain.BlogUserReaction, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogUserReactionRepo.GetUserReaction(c, blogID, userID)
}

func NewBlogUserReactionUsecase(blogUserReactionRepo domain.BlogUserReactionRepository, timeout time.Duration) domain.BlogUserReactionUsecase {
	return &blogUserReactionUsecase{
		blogUserReactionRepo: blogUserReactionRepo,
		ctxtimeout:           timeout,
	}
}
