package mongo

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type blogRepo struct {
	db        *mongo.Database
	collection string
}

// Create implements repositories.BlogRepository.
func (b blogRepo) Create(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	// example implement
	collection := b.db.Collection(b.collection)
	res, err := collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		blog.ID = oid
	} else {
		return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	return blog, nil
}

// Delete implements repositories.BlogRepository.
func (b blogRepo) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements repositories.BlogRepository.
func (b blogRepo) FindAll(ctx context.Context) ([]*domain.Blog, error) {
	panic("unimplemented")
}

// FindByTitle implements repositories.BlogRepository.
func (b blogRepo) FindByTitle(ctx context.Context, title string) ([]*domain.Blog, error) {
	panic("unimplemented")
}

// Update implements repositories.BlogRepository.
func (b blogRepo) Update(ctx context.Context, id string, blog *domain.Blog) (*domain.Blog, error) {
	panic("unimplemented")
}
