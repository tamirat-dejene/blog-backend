package domain

import (
	"context"
	"time"
)

type Blog struct {
	ID              string
	Title           string
	Content         string
	AuthorID        string
	AuthorName      string // for easy access to author's name: first_name + last_name
	Tags            []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Likes           int
	Dislikes        int
	ViewCount       int
	CommentCount    int // for easy access to comment count
	PopularityScore int // computed popularity score : score =  Normalized((likes * 3) + (views * 2) + (comments * 1.5) - (dislikes * 2.5))
}

type BlogComment struct {
	ID        string
	BlogID    string
	AuthorID  string
	Comment   string
	CreatedAt time.Time
}

type BlogUserReaction struct {
	ID        string
	BlogID    string
	UserID    string
	IsLike    bool
	CreatedAt time.Time
}

// BlogFilter defines filtering and pagination options for querying blogs.
type Recency string

const (
	RecencyNewest Recency = "newest"
	RecencyOldest Recency = "oldest"
)

type BlogFilter struct {
	Page       int
	PageSize   int
	Recency    Recency
	Tags       []string
	AuthorName string
	Title      string
	Popular    bool // indicates if the filter is for most popular blogs
}

// Repository Interfaces provide an abstraction layer for data access operations related to blogs, comments, and user reactions.
type BlogRepository interface {
	Create(ctx context.Context, blog *Blog) (*Blog, error)
	Update(ctx context.Context, id string, blog Blog) (Blog, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, filter *BlogFilter) ([]Blog, error)

	//... more methods can be added based on the usecases
}

type BlogCommentRepository interface {
	Create(ctx context.Context, comment BlogComment) (Blog, error)
	Delete(ctx context.Context, id string) error
}

type BlogUserReactionRepository interface {
	Create(ctx context.Context, reaction BlogUserReaction) (BlogUserReaction, error)
	Delete(ctx context.Context, id string) error
	GetUserReaction(ctx context.Context, blogID, userID string) (BlogUserReaction, error)
}

// Usecase Interfaces define the business logic for handling blogs, comments, and user reactions.
type BlogUsecase interface {
	GetBlogs(ctx context.Context, filter *BlogFilter) ([]Blog, error)
	CreateBlog(ctx context.Context, blog *Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, blog Blog) (Blog, error)
	DeleteBlog(ctx context.Context, id string) error
}

type BlogCommentUsecase interface {
	CreateComment(ctx context.Context, comment BlogComment) (Blog, error)
	DeleteComment(ctx context.Context, id string) error
}

type BlogUserReactionUsecase interface {
	CreateReaction(ctx context.Context, reaction BlogUserReaction) (BlogUserReaction, error)
	DeleteReaction(ctx context.Context, id string) error
	GetUserReaction(ctx context.Context, blogID, userID string) (BlogUserReaction, error)
}
