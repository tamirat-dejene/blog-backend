package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
  
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogPostRepo struct {
	db         mongo.Database
	collections *collections
}

// Create implements domain.BlogRepository.
func (b *blogPostRepo) Create(ctx context.Context, blog *domain.BlogPost) (*domain.BlogPost, error) {
	blog_model, err := mapper.BlogFromDomain(blog)
	if err != nil {
		return nil, err
	}
	insertedID, err := b.db.Collection(b.collections.BlogPosts).InsertOne(ctx, blog_model)
	if err != nil {
		return nil, err
	}
	blog.ID = insertedID.(string)
	return blog, nil
}

// Delete implements domain.BlogRepository.
func (b *blogPostRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = b.db.Collection(b.collections.BlogPosts).DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// Get implements domain.BlogRepository.
func (b *blogPostRepo) Get(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPost, error) {
	collection := b.db.Collection(b.collections.BlogPosts)

	query := utils.BuildBlogPostFilterQuery(filter)
	var pipeline []bson.D

	// Always start with a match stage
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: query}})

	// Sorting logic
	switch {
	case filter.Popular:
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "popularity_score", Value: -1}}}})
	case filter.Recency != "":
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: utils.RecencySort(filter.Recency)}})
	default:
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}})
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

	var results []mapper.BlogPostModel
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	blogs := make([]domain.BlogPost, 0, len(results))
	for _, bm := range results {
		blogs = append(blogs, *mapper.BlogToDomain(&bm))
	}

	return blogs, nil
}

// Update implements domain.BlogRepository.
func (b *blogPostRepo) Update(ctx context.Context, id string, blog domain.BlogPost) (domain.BlogPost, error) {
	oid, err := primitive.ObjectIDFromHex(blog.ID)

	if err != nil {
		return domain.BlogPost{}, err
	}
	blog.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":      blog.Title,
			"content":    blog.Content,
			"tags":       blog.Tags,
			"updated_at": blog.UpdatedAt,
		},
	}

	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, oid, update)
	return domain.BlogPost{}, err
}

func NewBlogPostRepo(database mongo.Database, collections *collections) domain.BlogPostRepository {
	return &blogPostRepo{
		db:         database,
		collections: collections,
	}
}
