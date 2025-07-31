package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
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
	collection := b.db.Collection(b.collection)

	query := utils.BuildBlogFilterQuery(filter)
	opts := utils.PaginationOpts(filter.Page, filter.PageSize, filter.Recency)

	cursor, err := collection.Find(ctx, query, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []mapper.BlogModel

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	blogs := make([]domain.Blog, 0, len(results))
	for i, bm := range results {
		blogs[i] = *mapper.BlogToDomain(&bm)

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