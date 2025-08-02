package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
)

type blogAIRepository struct {
	db          mongo.Database
	collections *Collections
}

// GetGeneratedContentByID implements domain.BlogAIRepository.
func (b *blogAIRepository) GetGeneratedContentByID(ctx context.Context, id string) (*domain.BlogAIContent, error) {
	panic("unimplemented")
}

// GetPromptsByUserID implements domain.BlogAIRepository.
func (b *blogAIRepository) GetPromptsByUserID(ctx context.Context, userID string) ([]domain.BlogAIPrompt, error) {
	panic("unimplemented")
}

// SaveFeedback implements domain.BlogAIRepository.
func (b *blogAIRepository) SaveFeedback(ctx context.Context, feedback domain.BlogAIFeedback) error {
	panic("unimplemented")
}

// StoreGeneratedContent implements domain.BlogAIRepository.
func (b *blogAIRepository) StoreGeneratedContent(ctx context.Context, content *domain.BlogAIContent) error {
	panic("unimplemented")
}

func NewBlogAIRepository(db mongo.Database, collections *Collections) domain.BlogAIRepository {
	return &blogAIRepository{
		db:          db,
		collections: collections,
	}
}
