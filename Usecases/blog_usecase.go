// this package will be the interface betweeen the controller/s and the repository
package usecases

import (
	"context"
	domain "g6/blog-api/Domain"
	"time"
)

type blogUsecase struct {
	blogRepo   domain.BlogRepository
	ctxtimeout time.Duration
}

// CreateBlog implements domain.BlogUsecase.
func (b *blogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog) (*domain.Blog, error) {
	// example implement
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogRepo.Create(c, blog)
}

// DeleteBlog implements domain.BlogUsecase.
func (b *blogUsecase) DeleteBlog(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetBlogs implements domain.BlogUsecase.
func (b *blogUsecase) GetBlogs(ctx context.Context, filter *domain.BlogFilter) ([]domain.Blog, error) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogRepo.Get(c, filter)
}

// UpdateBlog implements domain.BlogUsecase.
func (b *blogUsecase) UpdateBlog(ctx context.Context, id string, blog domain.Blog) (domain.Blog, error) {
	c, cancel := context.WithTimeout(ctx, b.ctxtimeout)
	defer cancel()

	return b.blogRepo.Update(c, id, blog)
}

func NewBlogUsecase(blogRepo domain.BlogRepository, timeout time.Duration) domain.BlogUsecase {
	return &blogUsecase{
		blogRepo:   blogRepo,
		ctxtimeout: timeout,
	}
}
