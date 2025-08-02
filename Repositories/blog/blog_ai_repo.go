package repository

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogAIRepository struct {
	db          mongo.Database
	collections *mongo.Collections
}

// GetGeneratedContentByID implements domain.BlogAIRepository.
func (b *blogAIRepository) GetGeneratedContentByID(ctx context.Context, id string) (*domain.BlogAIContent, *domain.DomainError) {
	coll := b.db.Collection(b.collections.AIBlogPosts)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid ID format: %w", err),
			Code: 400,
		}
	}

	var contentModel mapper.BlogAIContentModel
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&contentModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments()) {
			return nil, &domain.DomainError{
				Err:  errors.New("AI-generated blog not found"),
				Code: 404,
			}
		}
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("error fetching AI content: %w", err),
			Code: 500,
		}
	}

	return mapper.BlogAIContentToDomain(&contentModel), nil
}

// SaveFeedback implements domain.BlogAIRepository.
func (b *blogAIRepository) SaveFeedback(ctx context.Context, feedback domain.BlogAIFeedback) (*domain.BlogAIFeedback, *domain.DomainError) {
	feedbackModel, err := mapper.BlogAIFeedbackFromDomain(&feedback)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("error mapping feedback: %w", err),
			Code: 500,
		}
	}

	coll := b.db.Collection(b.collections.AIBlogPostsFeedback)
	_, err = coll.InsertOne(ctx, feedbackModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("error saving feedback: %w", err),
			Code: 500,
		}
	}

	return mapper.BlogAIFeedbackToDomain(feedbackModel), nil
}


// StoreGeneratedContent implements domain.BlogAIRepository.
func (b *blogAIRepository) StoreGeneratedContent(ctx context.Context, content *domain.BlogAIContent) (*domain.BlogAIContent, *domain.DomainError) {
	contentModel, err := mapper.BlogAIContentFromDomain(content)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("error mapping content: %w", err),
			Code: 500,
		}
	}

	coll := b.db.Collection(b.collections.AIBlogPosts)
	_, err = coll.InsertOne(ctx, contentModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("error storing generated content: %w", err),
			Code: 500,
		}
	}

	return mapper.BlogAIContentToDomain(contentModel), nil
}

func NewBlogAIRepository(db mongo.Database, collections *mongo.Collections) domain.BlogAIRepository {
	return &blogAIRepository{
		db:          db,
		collections: collections,
	}
}
