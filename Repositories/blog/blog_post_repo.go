package repository

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogPostRepo struct {
	db          mongo.Database
	collections *mongo.Collections
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
func (b *blogPostRepo) Get(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPostsPage, error) {
	collection := b.db.Collection(b.collections.BlogPosts)
	pipeline := utils.BuildBlogRetrievalAggregationPipeline(filter)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate blog posts failed: %w", err)
	}
	defer cursor.Close(ctx)

	var dbResults []mapper.BlogPostModel
	if err := cursor.All(ctx, &dbResults); err != nil {
		return nil, fmt.Errorf("decoding blog posts failed: %w", err)
	}

	return utils.PaginateBlogs(dbResults, filter.PageSize), nil
}

// GetBlogByID implements domain.BlogRepository.
func (b *blogPostRepo) GetBlogByID(ctx context.Context, id string)(*domain.BlogPost, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %w", err)
	}

	var blogModel mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
	if err != nil {
		return nil, fmt.Errorf("failed to find blog by ID: %w", err)
	}

	blog := mapper.BlogToDomain(&blogModel)
	
	return blog, nil
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

func NewBlogPostRepo(database mongo.Database, collections *mongo.Collections) domain.BlogPostRepository {
	return &blogPostRepo{
		db:          database,
		collections: collections,
	}
}
