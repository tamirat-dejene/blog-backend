package dto

import (
	domain "g6/blog-api/Domain"
	"time"
)

type BlogPostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

type BlogPostResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	AuthorID        string    `json:"author_id"`
	AuthorName      string    `json:"author_name"` // for easy access to author's name: first_name + last_name
	Tags            []string  `json:"tags,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Likes           int       `json:"likes"`
	Dislikes        int       `json:"dislikes"`
	ViewCount       int       `json:"view_count"`
	CommentCount    int       `json:"comment_count"`    // for easy access to comment count
	PopularityScore float64   `json:"popularity_score"` // computed popularity score
}

type BlogPostsPageResponse struct {
	Blogs      []BlogPostResponse `json:"blogs"`
	PageSize   int                `json:"page_size"`
	PageNumber int                `json:"page_number"`
}

type BlogUserReactionRequest struct {
	BlogID string `json:"blog_id" binding:"required"`
	IsLike bool   `json:"is_like" binding:"required"`
}

type BlogUserReactionResponse struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	IsLike    bool      `json:"is_like"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogCommentRequest struct {
	BlogID  string `json:"blog_id"`
	Comment string `json:"comment"`
}

type BlogCommentResponse struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	AuthorID  string    `json:"author_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type ReactionQuery struct {
	BlogId string `form:"blog_id" binding:"required"`
}

func (b *BlogPostRequest) ToDomain() *domain.BlogPost {
	return &domain.BlogPost{
		Title:           b.Title,
		Content:         b.Content,
		Tags:            b.Tags,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Likes:           0,
		Dislikes:        0,
		ViewCount:       0,
		CommentCount:    0,
		PopularityScore: 0,
	}
}

func (b *BlogPostResponse) Parse(blog *domain.BlogPost) {
	b.ID = blog.ID
	b.Title = blog.Title
	b.Content = blog.Content
	b.AuthorID = blog.AuthorID
	b.AuthorName = blog.AuthorName
	b.Tags = blog.Tags
	b.CreatedAt = blog.CreatedAt
	b.UpdatedAt = blog.UpdatedAt
	b.Likes = blog.Likes
	b.Dislikes = blog.Dislikes
	b.ViewCount = blog.ViewCount
	b.CommentCount = blog.CommentCount
	b.PopularityScore = blog.PopularityScore
}

func (r *BlogUserReactionRequest) ToDomain() *domain.BlogUserReaction {
	return &domain.BlogUserReaction{
		BlogID:    r.BlogID,
		IsLike:    r.IsLike,
		CreatedAt: time.Now(),
	}
}

func (r *BlogUserReactionResponse) Parse(reaction *domain.BlogUserReaction) {
	r.ID = reaction.ID
	r.BlogID = reaction.BlogID
	r.UserID = reaction.UserID
	r.IsLike = reaction.IsLike
	r.CreatedAt = reaction.CreatedAt
}

func (b *BlogCommentRequest) ToDomain() *domain.BlogComment {
	return &domain.BlogComment{
		BlogID:    b.BlogID,
		Comment:   b.Comment,
		CreatedAt: time.Now(),
	}
}

func (b *BlogCommentResponse) Parse(comment *domain.BlogComment) {
	b.ID = comment.ID
	b.BlogID = comment.BlogID
	b.AuthorID = comment.AuthorID
	b.Comment = comment.Comment
	b.CreatedAt = comment.CreatedAt
}

func (pr *BlogPostsPageResponse) Parse(page *domain.BlogPostsPage) {
	pr.PageSize = page.PageSize
	pr.PageNumber = page.PageNumber
	pr.Blogs = make([]BlogPostResponse, len(page.Blogs))
	for i, blog := range page.Blogs {
		pr.Blogs[i] = BlogPostResponse{}
		pr.Blogs[i].Parse(&blog)
	}
}
