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


import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogCommentUsecase struct {
	commentRepo domain.BlogCommentRepository
	ctxtimeout  time.Duration
}

// CreateComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) CreateComment(ctx context.Context, comment domain.BlogComment) (*domain.BlogComment, error) {
	panic("unimplemented")
}

// DeleteComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) DeleteComment(ctx context.Context, id string) error {
	panic("unimplemented")
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
