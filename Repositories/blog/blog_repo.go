package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"

	"go.mongodb.org/mongo-driver/bson"
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
	var pipeline []bson.D

	// Always start with a match stage
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: query}})

	// If popularity is enabled, add sorting by computed popularity_score
	if filter.Popular {
		pipeline = append(pipeline, utils.PopularityStages()...)
	} else {
		// Otherwise, use recency-based sort
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: utils.RecencySort(filter.Recency)}})
	}

	// Pagination stage
	skip := max((filter.Page - 1) * filter.PageSize, 0)
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: filter.PageSize}},
	)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []mapper.BlogModel
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	blogs := make([]domain.Blog, 0, len(results))
	for _, bm := range results {
		blogs = append(blogs, *mapper.BlogToDomain(&bm))
	}

	return blogs, nil
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
