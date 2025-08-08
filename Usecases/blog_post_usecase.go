// this package will be the interface betweeen the controller/s and the repository
package usecases

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"g6/blog-api/Infrastructure/redis"
	"net/http"
	"time"
)

type blogPostUsecase struct {
	blogPostRepo domain.BlogPostRepository
	redisClient  redis.RedisClient
	ctxtimeout   time.Duration
}

// CreateBlog implements domain.BlogUsecase.
func (b *blogPostUsecase) CreateBlog(ctx context.Context, blog *domain.BlogPost) (*domain.BlogPost, *domain.DomainError) {
	// example implement
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogPostRepo.Create(c, blog)
}

// DeleteBlog implements domain.BlogUsecase.
func (b *blogPostUsecase) DeleteBlog(ctx context.Context, id string) *domain.DomainError {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	// Delete the blog post from the database
	if err := b.blogPostRepo.Delete(c, id); err != nil {
		return err
	}

	// Use the actual Redis service from the usecase to generate key
	redisKey := b.redisClient.Service().GenerateBlogPostKey(id)
	if err := b.redisClient.Delete(ctx, redisKey); err != nil {
		return &domain.DomainError{
			Err:  errors.New("failed to invalidate blog post cache"),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

// GetBlogs implements domain.BlogUsecase.
func (b *blogPostUsecase) GetBlogs(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPostsPage, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	// Generate the Redis key
	redis_key := b.redisClient.Service().GenerateRedisKey(filter)

	// Check the Redis cache first
	cachedPages, err := b.redisClient.Get(ctx, redis_key)
	if err == nil && cachedPages != "" {
		fmt.Println("Cache hit for key:", redis_key)

		pageModels, err := utils.DeserializeBlogPostsPage(cachedPages)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to deserialize blog posts page"),
				Code: http.StatusInternalServerError,
			}
		}

		return utils.PaginateBlogs(pageModels, filter.PageSize), nil
	}

	fmt.Println("Cache miss for key:", redis_key)

	// If not found in cache, query the database
	blogPosts, serialized, domErr := b.blogPostRepo.Get(c, filter)
	if domErr != nil {
		return nil, domErr
	}

	// Cache the serialized data
	if err := b.redisClient.Set(ctx, redis_key, serialized, b.redisClient.GetCacheExpiry()); err != nil {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to set blog posts page in cache"),
			Code: http.StatusInternalServerError,
		}
	}

	return blogPosts, nil
}

// GetBlogByID implements domain.BlogUsecase.
func (b *blogPostUsecase) GetBlogByID(ctx context.Context, user_id, blog_id string) (*domain.BlogPost, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	redisKey := b.redisClient.Service().GenerateBlogPostKey(blog_id)
	cachedBlog, err := b.redisClient.Get(ctx, redisKey)

	var blog *domain.BlogPost
	var err1 *domain.DomainError

	if err != nil || cachedBlog == "" {
		fmt.Println("Cache miss for key:", redisKey)
		// Fetch from DB
		blog, err1 = b.blogPostRepo.GetBlogByID(c, blog_id)
		if err1 != nil {
			return nil, err1
		}	

		// Serialize and cache it
		blogModel := &mapper.BlogPostModel{}
		if err := blogModel.Parse(blog); err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to parse blog post"),
				Code: http.StatusInternalServerError,
			}
		}

		serialized, err := utils.SerializeBlogPost(blogModel)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to serialize blog post"),
				Code: http.StatusInternalServerError,
			}
		}

		if err := b.redisClient.Set(ctx, redisKey, serialized, b.redisClient.GetCacheExpiry()); err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to set blog post in cache"),
				Code: http.StatusInternalServerError,
			}
		}
	} else {
		fmt.Println("Cache hit for key:", redisKey)
		// Deserialize the cached blog post
		blogModel, err := utils.DeserializeBlogPost(cachedBlog)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  errors.New("failed to deserialize blog post"),
				Code: http.StatusInternalServerError,
			}
		}

		// Convert the model to domain object
		blog = blogModel.ToDomain()
	}

	// Increment view count (can be async in future)
	err12 := b.IncrementViewCountWithLimit(c, user_id, blog_id)

	// Update popularity score
	if err12 != nil {
		return b.blogPostRepo.RefreshPopularityScore(c, blog_id)
	}

	return blog, nil

}

// UpdateBlog implements domain.BlogUsecase.
func (b *blogPostUsecase) UpdateBlog(ctx context.Context, id string, blog domain.BlogPost) (*domain.BlogPost, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	// Perform the update
	updated, err := b.blogPostRepo.Update(c, id, blog)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for this blog post
	redisKey := b.redisClient.Service().GenerateBlogPostKey(id)
	if delErr := b.redisClient.Delete(ctx, redisKey); delErr != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to invalidate blog post cache: %w", delErr),
			Code: http.StatusInternalServerError,
		}
	}

	return updated, nil
}

// Users can view a blog only once, within three hours of the last view to prevent excessive view count increments.
// track user: "userId+blogId:blogId" to allow user view multiple blogs
func (b *blogPostUsecase) IncrementViewCountWithLimit(ctx context.Context, user_id, blog_id string) (*domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	// Compose a unique key for user-blog view tracking
	viewKey := fmt.Sprintf("view:%s:%s", user_id, blog_id)

	// Check if the user has viewed this blog in the last 3 hours
	exists, err := b.redisClient.Exists(ctx, viewKey)
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to check view limit: %w", err),
			Code: http.StatusInternalServerError,
		}
	}
	if exists {
		return &domain.DomainError{
			Err:  errors.New("view already counted within the last 3 hours"),
			Code: http.StatusTooManyRequests,
		}
	}

	// Increment the view count in the database
	_, derr := b.blogPostRepo.IncrementViewCount(c, blog_id)
	if derr != nil {
		return derr
	}

	// Set the view key in Redis with a 3-hour expiry
	if err := b.redisClient.Set(ctx, viewKey, "1", 3*time.Hour); err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to set view limit: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	return nil
}

// NewBlogPostUsecase creates a new instance of blog post usecase.
func NewBlogPostUsecase(blogPostRepo domain.BlogPostRepository, redisClient redis.RedisClient, timeout time.Duration) domain.BlogPostUsecase {
	return &blogPostUsecase{
		blogPostRepo: blogPostRepo,
		redisClient:  redisClient,
		ctxtimeout:   timeout,
	}
}
