package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
)

type blogRepo struct {
	db         mongo.Database
	collection string
}

// Create implements domain.BlogRepository.
func (b *blogRepo) Create(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	blog_model, err := mapper.BlogFromDomain(blog)
	if err != nil {
		return nil, err
	}
	insertedID, err := b.db.Collection(b.collection).InsertOne(ctx, blog_model)
	if err != nil {
		return nil, err
	}
	blog.ID = insertedID.(string)
	return blog, nil
}

// Delete implements domain.BlogRepository.
func (b *blogRepo) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// Get implements domain.BlogRepository.
func (b *blogRepo) Get(ctx context.Context, filter *domain.BlogFilter) ([]domain.Blog, error) {
	var blogs []domain.Blog
	var err error

	if filter != nil {
		// Apply filtering logic here
		// For example, you can use filter.Page, filter.PageSize, etc.
	}

	return blogs, err
}

// Update implements domain.BlogRepository.
func (b *blogRepo) Update(ctx context.Context, id string, blog domain.Blog) (domain.Blog, error) {
	panic("unimplemented")
}

func NewBlogRepo(database mongo.Database, collection string) domain.BlogRepository {
	return &blogRepo{
		db:         database,
		collection: collection,
	}
}
