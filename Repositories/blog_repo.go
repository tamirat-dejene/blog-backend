package repositories

import (
	"context"
	domain "g6/blog-api/Domain"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *domain.Blog) (*domain.Blog, error)
	Update(ctx context.Context, id string, blog *domain.Blog) (*domain.Blog, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*domain.Blog, error)

	// Find related blogs which title relates to the given title
	FindByTitle(ctx context.Context, title string) ([]*domain.Blog, error)

	//... more methods can be added based on the usecases
}

type BlogCommentRepository interface {
	Create(ctx context.Context, comment *domain.BlogComment) (*domain.Blog, error)
	// Comment update will not be a usecase hopefully
	Delete(ctx context.Context, id string) error
}
