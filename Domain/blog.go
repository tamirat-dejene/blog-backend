package domain

import (
	"context"
	"time"
)

type BlogPost struct {
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
	CommentCount    int     // for easy access to comment count
	PopularityScore float64 // computed popularity score : score =  Normalized((likes * 3) + (views * 2) + (comments * 1.5) - (dislikes * 2.5))
}

type BlogPostsPage struct {
	Blogs      []BlogPost
	PageNumber int
	PageSize   int
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

type BlogPostFilter struct {
	Page       int
	PageSize   int
	Recency    Recency
	Tags       []string
	AuthorName string
	Title      string
	Popular    bool // indicates if the filter is for most popular blogs
}

// Repository Interfaces provide an abstraction layer for data access operations related to blogs, comments, and user reactions.
type BlogPostRepository interface {
	Create(ctx context.Context, blog *BlogPost) (*BlogPost, *DomainError)
	Update(ctx context.Context, id string, blog BlogPost) (*BlogPost, *DomainError)
	Delete(ctx context.Context, id string) *DomainError
	Get(ctx context.Context, filter *BlogPostFilter) ([]BlogPostsPage, *string, *DomainError) // pages, serialized string for caching, error
	GetBlogByID(ctx context.Context, id string) (*BlogPost, *DomainError)
	RefreshPopularityScore(ctx context.Context, id string) (*BlogPost, *DomainError)
	IncrementViewCount(ctx context.Context, id string) (*BlogPost, *DomainError)
	UpdateCommentCount(ctx context.Context, id string, increment bool) (*BlogPost, *DomainError)
	UpdateReactionCount(ctx context.Context, is_like bool, id string, increment bool) (*BlogPost, *DomainError)

	//... more methods can be added based on the usecases
}

type BlogCommentRepository interface {
	Create(ctx context.Context, comment *BlogComment) (*BlogComment, *DomainError)
	Delete(ctx context.Context, id string) *DomainError
	Update(ctx context.Context, id string, comment *BlogComment) (*BlogComment, *DomainError)
	GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]BlogComment, *DomainError)
	GetCommentByID(ctx context.Context, id string) (*BlogComment, *DomainError)
}

type BlogUserReactionRepository interface {
	Create(ctx context.Context, reaction *BlogUserReaction) (*BlogUserReaction, *DomainError)
	Delete(ctx context.Context, id string) *DomainError
	GetUserReaction(ctx context.Context, blogID, userID string) (*BlogUserReaction, *DomainError)
}

// Usecase Interfaces define the business logic for handling blogs, comments, and user reactions.
type BlogPostUsecase interface {
	GetBlogs(ctx context.Context, filter *BlogPostFilter) ([]BlogPostsPage, *DomainError)
	GetBlogByID(ctx context.Context, user_id, blog_id string) (*BlogPost, *DomainError)
	CreateBlog(ctx context.Context, blog *BlogPost) (*BlogPost, *DomainError)
	UpdateBlog(ctx context.Context, id string, blog BlogPost) (*BlogPost, *DomainError)
	DeleteBlog(ctx context.Context, id string) *DomainError
	IncrementViewCountWithLimit(ctx context.Context, user_id, blog_id string) (*DomainError)
}

type BlogCommentUsecase interface {
	CreateComment(ctx context.Context, comment *BlogComment) (*BlogComment, *DomainError)
	DeleteComment(ctx context.Context, id string) *DomainError
	GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]BlogComment, *DomainError)
	GetCommentByID(ctx context.Context, id string) (*BlogComment, *DomainError)
	UpdateComment(ctx context.Context, id string, comment *BlogComment) (*BlogComment, *DomainError)
}

type BlogUserReactionUsecase interface {
	CreateReaction(ctx context.Context, reaction *BlogUserReaction) (*BlogUserReaction, *DomainError)
	DeleteReaction(ctx context.Context, id string) *DomainError
	GetUserReaction(ctx context.Context, blogID, userID string) (*BlogUserReaction, *DomainError)
}
