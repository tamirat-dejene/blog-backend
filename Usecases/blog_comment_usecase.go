package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogCommentUsecase struct {
	blogCommentRepo domain.BlogCommentRepository
	ctxtimeout      time.Duration
}

func (b *blogCommentUsecase) CreateComment(ctx context.Context, comment domain.BlogComment) (domain.BlogComment, error) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogCommentRepo.Create(c, comment)
}

func (b *blogCommentUsecase) DeleteComment(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogCommentRepo.Delete(c, id)
}

func NewBlogCommentUseCase(blogCommentRepo domain.BlogCommentRepository, timeout time.Duration) domain.BlogCommentUsecase {
	return &blogCommentUsecase{
		blogCommentRepo: blogCommentRepo,
		ctxtimeout:      timeout,
	}

}
