package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogCommentRepository struct {
	db          mongo.Database
	collections *collections
}

func (b *blogCommentRepository) Create(ctx context.Context, comment domain.BlogComment) (*domain.BlogComment, error) {
	blogComment, err := mapper.BlogCommentFromDomain(&comment)

	if err != nil {
		return &domain.BlogComment{}, err
	}
	_, err = b.db.Collection(b.collections.BlogComments).InsertOne(ctx, blogComment)

	if err != nil {
		return &domain.BlogComment{}, err
	}
	res := mapper.BlogCommentToDomain(blogComment)
	return res, nil
}

func (b *blogCommentRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = b.db.Collection(b.collections.BlogComments).DeleteOne(ctx, filter)

	return err
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
