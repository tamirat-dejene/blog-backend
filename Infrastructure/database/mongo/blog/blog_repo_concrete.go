package mongo

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	repositories "g6/blog-api/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type blogRepo struct {
	db         mongo.Database
	collection string
}

// Create implements repositories.BlogRepository.
func (b blogRepo) Create(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	// example implement

	blog_model, err := mapper.FromDomain(blog)
	if err != nil {
		return nil, fmt.Errorf("failed to convert blog to model: %w", err)
	}

	coll := b.db.Collection(b.collection)
	res, err := coll.InsertOne(ctx, blog_model)
	if err != nil {
		return nil, fmt.Errorf("failed to insert blog: %w", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to get inserted ID")
	}
	blog.ID = oid.Hex()
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

func NewBlogRepo(database mongo.Database, collection string) repositories.BlogRepository {
	return &blogRepo{
		db:         database,
		collection: collection,
	}
}
