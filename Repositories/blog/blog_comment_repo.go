package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
)

type blogCommentRepository struct {
	db          mongo.Database
	collections *collections
}

// Create implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Create(ctx context.Context, comment domain.BlogComment) (*domain.BlogComment, error) {
	panic("unimplemented")
}

// Delete implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetCommentByID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, error) {
	panic("unimplemented")
}

// GetCommentsByBlogID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentsByBlogID(ctx context.Context, blogID string) ([]domain.BlogComment, error) {
	panic("unimplemented")
}

// Update implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Update(ctx context.Context, id string, comment domain.BlogComment) (*domain.BlogComment, error) {
	panic("unimplemented")
}

func NewBlogCommentRepository(db mongo.Database, collections *collections) domain.BlogCommentRepository {
	return &blogCommentRepository{
		db:          db,
		collections: collections,
	}
}
