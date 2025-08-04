package usecases

import (
	"context"
	"errors"
	"fmt"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/ai"
	"time"
)

type blogAIUsecase struct {
	geminiConfig ai.GeminiConfig
	ctxtimeout   time.Duration
}

// GeneratePost implements domain.AIBlogPostUsecase.
func (b *blogAIUsecase) GeneratePost(ctx context.Context, req domain.AIBlogPostGenerate) (*domain.AIBlogPost, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate content using the AI service.
	content, err := b.geminiConfig.GenerateWithGemini(c, req.Topic, req.Keywords)
	if err != nil {
		fmt.Println("Error generating content:", err)
		return nil, &domain.DomainError{
			Err:  errors.New("failed to generate content"),
			Code: 500,
		}
	}

	// Parse the generated content from JSON.
	var generatedContent dto.AIBlogPostResponse
	if err1 := b.geminiConfig.ParseGeneratedContent(content, &generatedContent); err1 != nil {
		return nil, &domain.DomainError{
			Err:  err1,
			Code: 500,
		}
	}
	generatedContent.Topic = req.Topic
	generatedContent.Keywords = req.Keywords
	generatedContent.CreatedAt = time.Now().Format(time.RFC3339)

	return generatedContent.ToDomain(), nil
}

func NewBlogAIUsecase(geminiConfig ai.GeminiConfig, timeout time.Duration) domain.AIBlogPostUsecase {
	return &blogAIUsecase{
		geminiConfig: geminiConfig,
		ctxtimeout:   timeout,
	}
}
