package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
)

type blogAIRepository struct {
	db          mongo.Database
	collections *mongo.Collections
}

// GetGeneratedContentByID implements domain.BlogAIRepository.
func (b *blogAIRepository) GetGeneratedContentByID(ctx context.Context, id string) (*domain.BlogAIContent, *domain.DomainError) {
	panic("unimplemented")
}

// SaveFeedback implements domain.BlogAIRepository.
func (b *blogAIRepository) SaveFeedback(ctx context.Context, feedback domain.BlogAIFeedback) *domain.DomainError {
	panic("unimplemented")
}

// StoreGeneratedContent implements domain.BlogAIRepository.
func (b *blogAIRepository) StoreGeneratedContent(ctx context.Context, content *domain.BlogAIContent) (*domain.BlogAIContent, *domain.DomainError) {
	panic("unimplemented")
}

func NewBlogAIRepository(db mongo.Database, collections *mongo.Collections) domain.BlogAIRepository {
	return &blogAIRepository{
		db:          db,
		collections: collections,
	}
}
