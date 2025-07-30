package domain

import (
	"context"
	"time"
)

type Blog struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     int
	Dislikes  int
	ViewCount int
}

type BlogComment struct {
	ID        string
	BlogID    string
	AuthorID  string
	Comment   string
	CreatedAt time.Time
}

// interface for blog data operations.
type BlogRepository interface {
	Create(ctx context.Context, blog *Blog) (*Blog, error)
	Update(ctx context.Context, id string, blog Blog) (Blog, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]Blog, error)

	// Find related blogs which title relates to the given title
	FindByTitle(ctx context.Context, title string) ([]Blog, error)

	//... more methods can be added based on the usecases
}

type BlogCommentRepository interface {
	Create(ctx context.Context, comment BlogComment) (Blog, error)
	Delete(ctx context.Context, id string) error
}

// interface for blog data operations
type BlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, blog Blog) (Blog, error)
	DeleteBlog(ctx context.Context, id string) error
	GetAllBlogs(ctx context.Context) ([]Blog, error)
	GetBlogsByTitle(ctx context.Context, title string) ([]Blog, error)
}

type BlogCommentUsecase interface {
	CreateComment(ctx context.Context, comment BlogComment) (Blog, error)
	DeleteComment(ctx context.Context, id string) error
}
