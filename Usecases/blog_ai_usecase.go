package usecases

import (
	"context"
	"errors"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/ai"
	"time"
)

type blogAIUsecase struct {
	blogAIRepo   domain.BlogAIRepository
	geminiConfig ai.GeminiConfig
	ctxtimeout   time.Duration
}

func (b *blogAIUsecase) GenerateContent(ctx context.Context, req domain.BlogAIGenerate) (*domain.BlogAIContent, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	// Generate content using the AI service.
	content, err := b.geminiConfig.GenerateWithGemini(c, req.Topic, req.Keywords)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to generate content"),
			Code: 500,
		}
	}

	// Parse the generated content from JSON.
	var generatedContent dto.BlogAIResponseDTO
	if err1 := b.geminiConfig.ParseGeneratedContent(content, &generatedContent); err1 != nil {
		return nil, &domain.DomainError{
			Err:  err1,
			Code: 500,
		}
	}
	generatedContent.Topic = req.Topic
	generatedContent.Keywords = req.Keywords

	// Store the generated content in the repository.
	ai_blog, err1 := b.blogAIRepo.StoreGeneratedContent(c, dto.BlogAIContentToDomain(&generatedContent, req.UserID, time.Now().UTC()))
	if err1 != nil {
		return nil, &domain.DomainError{
			Err:  err1.Err,
			Code: 500,
		}
	}

	return ai_blog, nil
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

func NewBlogAIUsecase(blogAIRepo domain.BlogAIRepository, geminiConfig ai.GeminiConfig, timeout time.Duration) domain.BlogAIUsecase {
	return &blogAIUsecase{
		blogAIRepo:   blogAIRepo,
		geminiConfig: geminiConfig,
		ctxtimeout:   timeout,
	}
}
