package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogAIUsecase struct {
	blogAIRepo domain.BlogAIRepository
	ctxtimeout time.Duration
}

func (b *blogAIUsecase) GenerateContent(ctx context.Context, req domain.BlogAIGenerate) (*domain.BlogAIContent, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	content := &domain.BlogAIContent{
		UserID:    req.UserID,
		Topic:     req.Topic,
		Keywords:  req.Keywords,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	err := b.blogAIRepo.StoreGeneratedContent(c, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (b *blogAIUsecase) GetGeneratedContentByID(ctx context.Context, id string) (*domain.BlogAIContent, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogAIRepo.GetGeneratedContentByID(c, id)
}

func (b *blogAIUsecase) GetUserPromptHistory(ctx context.Context, userID string) ([]domain.BlogAIPrompt, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogAIRepo.GetPromptsByUserID(c, userID)
}

func (b *blogAIUsecase) SubmitFeedback(ctx context.Context, feedback domain.BlogAIFeedback) *domain.DomainError {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogAIRepo.SaveFeedback(c, feedback)
}

func NewBlogAIUsecase(blogAIRepo domain.BlogAIRepository, timeout time.Duration) domain.BlogAIUsecase {
	return &blogAIUsecase{
		blogAIRepo: blogAIRepo,
		ctxtimeout: timeout,
	}
}
