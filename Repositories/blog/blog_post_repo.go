package repository

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"g6/blog-api/Infrastructure/redis"
	"net/http"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogPostRepo struct {
	db          mongo.Database
	redisClient redis.RedisClient
	collections *mongo.Collections
}

// Create implements domain.BlogRepository.
func (b *blogPostRepo) Create(ctx context.Context, blog *domain.BlogPost) (*domain.BlogPost, *domain.DomainError) {
	// Map the domain model to the DB model
	blogModel := &mapper.BlogPostModel{}
	if err := blogModel.Parse(blog); err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Insert the blog into the collection
	result, err := b.db.Collection(b.collections.BlogPosts).InsertOne(ctx, blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Extract the inserted ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to cast inserted ID to ObjectID"),
			Code: http.StatusInternalServerError,
		}
	}

	// Set the generated ID back to domain model
	blog.ID = objectID.Hex()
	return blog, nil
}

// Delete implements domain.BlogRepository.
func (b *blogPostRepo) Delete(ctx context.Context, id string) *domain.DomainError {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	_, err = b.db.Collection(b.collections.BlogPosts).DeleteOne(ctx, bson.M{"_id": oid})
	if err == mongo.ErrNoDocuments() {
		return &domain.DomainError{
			Err:  fmt.Errorf("blog post with ID %s not found", id),
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Invalidate the cache for this blog post
	redis_service := &redis.RedisService{}
	redis_key := redis_service.GenerateBlogPostKey(id)
	if err := b.redisClient.Delete(ctx, redis_key); err != nil {
		return &domain.DomainError{
			Err:  errors.New("failed to delete blog post from cache"),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

// Get implements domain.BlogRepository.
func (b *blogPostRepo) Get(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPostsPage, *domain.DomainError) {
	// generate the Redis key
	redis_service := &redis.RedisService{}
	redis_key := redis_service.GenerateRedisKey(filter)

	// Check the redis cache first
	pages, err := b.redisClient.Get(ctx, redis_key)
	if err == nil && pages != "" {
		fmt.Println("Cache hit for key:", redis_key)
		page_models, err := utils.DeserializeBlogPostsPage(pages)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to deserialize blog posts page"),
				Code: http.StatusInternalServerError,
			}
		}

		return utils.PaginateBlogs(page_models, filter.PageSize), nil
	}
	fmt.Println("Cache miss for key:", redis_key)

	collection := b.db.Collection(b.collections.BlogPosts)
	pipeline := utils.BuildBlogRetrievalAggregationPipeline(filter)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer cursor.Close(ctx)

	var dbResults []mapper.BlogPostModel
	if err := cursor.All(ctx, &dbResults); err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	serialized, err := utils.SerializeBlogPostsPage(dbResults)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to serialize blog posts page"),
			Code: http.StatusInternalServerError,
		}
	}

	b.redisClient.Set(ctx, redis_key, serialized, b.redisClient.GetCacheExpiry())

	return utils.PaginateBlogs(dbResults, filter.PageSize), nil
}

// GetBlogByID implements domain.BlogRepository.
func (b *blogPostRepo) GetBlogByID(ctx context.Context, id string) (*domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	// Check the Redis cache first
	var blogModel *mapper.BlogPostModel
	redis_service := &redis.RedisService{}
	redis_key := redis_service.GenerateBlogPostKey(id)

	value, err := b.redisClient.Get(ctx, redis_key)
	if err == nil && value != "" {
		fmt.Println("Cache hit for key:", redis_key)
		blogModel, err = utils.DeserializeBlogPost(value)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to deserialize blog post"),
				Code: http.StatusInternalServerError,
			}
		}
	} else {
		fmt.Println("Cache miss for key:", redis_key)
		err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
		if err == mongo.ErrNoDocuments() {
			return nil, &domain.DomainError{
				Err:  err,
				Code: http.StatusNotFound,
			}
		} else if err != nil {
			return nil, &domain.DomainError{
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
	}

	// Increment the view count
	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$inc": bson.M{"view_count": 1},
	})
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	blogModel.ViewCount += 1

	// Convert the model to domain
	if blogModel.ID.IsZero() {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusNotFound,
		}
	}

	// Calculate the popularity score
	ps := utils.CalculatePopularityScore(blogModel.Likes, blogModel.ViewCount, blogModel.CommentCount, blogModel.Dislikes)
	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"popularity_score": ps},
	})
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	blogModel.PopularityScore = int(ps)

	blog := blogModel.ToDomain()
	if blog == nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to convert blog model to domain for ID: %s", id),
			Code: http.StatusInternalServerError,
		}
	}

	// Cache the blog post
	serialized, err := utils.SerializeBlogPost(blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to serialize blog post"),
			Code: http.StatusInternalServerError,
		}
	}
	if err := b.redisClient.Set(ctx, redis_key, serialized, b.redisClient.GetCacheExpiry()); err != nil {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to cache blog post"),
			Code: http.StatusInternalServerError,
		}
	}

	// Set the ID back to the domain model
	blog.ID = blogModel.ID.Hex()
	return blog, nil
}

// Update implements domain.BlogRepository.
func (b *blogPostRepo) Update(ctx context.Context, id string, blog domain.BlogPost) (domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(blog.ID)

	if err != nil {
		return domain.BlogPost{}, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
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
	if err != nil {
		return domain.BlogPost{}, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Reset the cache for this blog post
	redis_service := &redis.RedisService{}
	redis_key := redis_service.GenerateBlogPostKey(id)
	if err := b.redisClient.Delete(ctx, redis_key); err != nil {
		return domain.BlogPost{}, &domain.DomainError{
			Err:  errors.New("failed to delete blog post from cache"),
			Code: http.StatusInternalServerError,
		}
	}

	// Return the updated blog
	blog.ID = oid.Hex()
	return blog, nil
}

func NewBlogPostRepo(database mongo.Database, redisClient redis.RedisClient, collections *mongo.Collections) domain.BlogPostRepository {
	return &blogPostRepo{
		db:          database,
		collections: collections,
		redisClient: redisClient,
	}
}
