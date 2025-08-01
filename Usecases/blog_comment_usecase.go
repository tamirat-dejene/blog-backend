package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogCommentUsecase struct {
	commentRepo domain.BlogCommentRepository
	ctxtimeout  time.Duration
}

func (b *blogCommentUsecase) CreateComment(ctx context.Context, comment domain.BlogComment) (*domain.BlogComment, error) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.commentRepo.Create(c, comment)
}

func (b *blogCommentUsecase) DeleteComment(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.commentRepo.Delete(c, id)
}

// GetCommentByID implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, error) {
	panic("unimplemented")
}

// GetCommentsByBlogID implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) GetCommentsByBlogID(ctx context.Context, blogID string) ([]domain.BlogComment, error) {
	panic("unimplemented")
}

// UpdateComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) UpdateComment(ctx context.Context, id string, comment domain.BlogComment) (*domain.BlogComment, error) {
	panic("unimplemented")
}

func NewBlogCommentUsecase(commentRepo domain.BlogCommentRepository, timeout time.Duration) domain.BlogCommentUsecase {
	return &blogCommentUsecase{
		commentRepo: commentRepo,
		ctxtimeout:  timeout,
	}
}
