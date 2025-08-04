// this package will be the interface betweeen the controller/s and the repository
package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogPostUsecase struct {
	blogPostRepo domain.BlogPostRepository
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

	return b.blogPostRepo.Delete(c, id)
}

// GetBlogs implements domain.BlogUsecase.
func (b *blogPostUsecase) GetBlogs(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPostsPage, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogPostRepo.Get(c, filter)
}

// GetBlogByID implements domain.BlogUsecase.
func (b *blogPostUsecase) GetBlogByID(ctx context.Context, id string) (*domain.BlogPost, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogPostRepo.GetBlogByID(c, id)
}

// UpdateBlog implements domain.BlogUsecase.
func (b *blogPostUsecase) UpdateBlog(ctx context.Context, id string, blog domain.BlogPost) (domain.BlogPost, *domain.DomainError) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogPostRepo.Update(c, id, blog)
}

func NewBlogPostUsecase(blogPostRepo domain.BlogPostRepository, timeout time.Duration) domain.BlogPostUsecase {
	return &blogPostUsecase{
		blogPostRepo: blogPostRepo,
		ctxtimeout:   timeout,
	}
}
