package usecases

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"g6/blog-api/Infrastructure/redis"
	"time"
)

type blogCommentUsecase struct {
	commentRepo domain.BlogCommentRepository
	redisClient redis.RedisClient
	ctxtimeout  time.Duration
}

// CreateComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) CreateComment(ctx context.Context, comment *domain.BlogComment) (*domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.commentRepo.Create(c, comment)
}

// DeleteComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) DeleteComment(ctx context.Context, id string) *domain.DomainError {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	if err := b.commentRepo.Delete(c, id); err != nil {
		return err
	}

	// Invalidate Redis cache for this comment
	redisService := b.redisClient.Service()
	redisKey := redisService.GenerateBlogCommentKey(id)
	if err := b.redisClient.Delete(c, redisKey); err != nil {
		return &domain.DomainError{
			Err:  errors.New("failed to delete comment from cache"),
			Code: 500,
		}
	}
	return nil
}

// GetCommentByID implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	redisService := b.redisClient.Service()
	redisKey := redisService.GenerateBlogCommentKey(id)

	// Try fetching from Redis cache
	cachedComment, err := b.redisClient.Get(c, redisKey)
	if err == nil && cachedComment != "" {
		comment, err := utils.DeserializeBlogComment(cachedComment)
		if err == nil {
			return comment.ToDomain(), nil
		}
	}

	// Fallback to DB
	comment, domainErr := b.commentRepo.GetCommentByID(c, id)
	if domainErr != nil {
		return nil, domainErr
	}

	// Cache the result (ignore cache write error)
	var model mapper.BlogCommentModel
	model.Parse(comment)
	serialized, err := utils.SerializeBlogComment(model)
	if err == nil {
		_ = b.redisClient.Set(c, redisKey, serialized, b.redisClient.GetCacheExpiry())
	}

	return comment, nil
}

// GetCommentsByBlogID implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.commentRepo.GetCommentsByBlogID(c, blogID, limit)
}

// UpdateComment implements domain.BlogCommentUsecase.
func (b *blogCommentUsecase) UpdateComment(ctx context.Context, id string, comment *domain.BlogComment) (*domain.BlogComment, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.commentRepo.Update(c, id, comment)
}

func NewBlogCommentUsecase(commentRepo domain.BlogCommentRepository, redisClient redis.RedisClient, timeout time.Duration) domain.BlogCommentUsecase {
	return &blogCommentUsecase{
		commentRepo: commentRepo,
		redisClient: redisClient,
		ctxtimeout:  timeout,
	}
}
