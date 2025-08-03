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
func (b *blogCommentUsecase) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	comment, err := b.commentRepo.GetCommentByID(c, id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err.Err,
			Code: err.Code,
		}
	}

	return comment, nil
}

// GetCommentsByBlogID implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	comments, err := b.commentRepo.GetCommentsByBlogID(c, blogID, limit)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err.Err,
			Code: err.Code,
		}
	}
	return comments, nil
}

// UpdateComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) UpdateComment(ctx context.Context, id string, comment domain.BlogComment) (*domain.BlogComment, error) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	updatedComment, err := b.commentRepo.Update(c, id, comment)

	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

func NewBlogCommentUsecase(commentRepo domain.BlogCommentRepository, timeout time.Duration) domain.BlogCommentUsecase {
	return &blogCommentUsecase{
		commentRepo: commentRepo,
		ctxtimeout:  timeout,
	}
}
