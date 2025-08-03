package dto

import (
	domain "g6/blog-api/Domain"
	"time"
)

type BlogPostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	AuthorID string   `json:"author_id" binding:"required"` // string, will convert to ObjectID in domain
	Tags     []string `json:"tags"`
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
	PopularityScore int       `json:"popularity_score"` // computed popularity score
}

type BlogUserReactionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	BlogID string `json:"blog_id" binding:"required"`
	IsLike bool   `json:"is_like"`
}

type BlogUserReactionResponse struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	IsLike    bool      `json:"is_like"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogCommentRequest struct {
	ID       string `json:"id"`
	BlogID   string `json:"blog_id"`
	AuthorID string `json:"author_id"`
	Comment  string `json:"comment"`
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
	UserId string `form:"user_id" binding:"required"`
}

func ToDomainBlogPost(req BlogPostRequest) domain.BlogPost {
	return domain.BlogPost{
		Title:           req.Title,
		Content:         req.Content,
		AuthorID:        req.AuthorID,
		Tags:            req.Tags,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Likes:           0,
		Dislikes:        0,
		ViewCount:       0,
		CommentCount:    0,
		PopularityScore: 0,
	}
}
func FromDomainBlogPost(post *domain.BlogPost) BlogPostResponse {
	return BlogPostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuthorID:        post.AuthorID,
		AuthorName:      post.AuthorName,
		Tags:            post.Tags,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		Likes:           post.Likes,
		Dislikes:        post.Dislikes,
		ViewCount:       post.ViewCount,
		CommentCount:    post.CommentCount,
		PopularityScore: post.PopularityScore,
	}
}

func ToDomainBlogReaction(req BlogUserReactionRequest) domain.BlogUserReaction {
	return domain.BlogUserReaction{
		UserID:    req.UserID,
		BlogID:    req.BlogID,
		CreatedAt: time.Now(),
		IsLike:    req.IsLike,
	}
}

func FromDomainBlogReaction(response domain.BlogUserReaction) BlogUserReactionResponse {
	return BlogUserReactionResponse{
		ID:        response.ID,
		BlogID:    response.BlogID,
		UserID:    response.UserID,
		CreatedAt: response.CreatedAt,
		IsLike:    response.IsLike,
	}
}

func ToDomainBlogComment(req BlogCommentRequest) domain.BlogComment {
	return domain.BlogComment{
		ID:       req.ID,
		BlogID:   req.BlogID,
		AuthorID: req.AuthorID,
		Comment:  req.Comment,
	}
}

func FromDomainBlogComment(response domain.BlogComment) BlogCommentResponse {
	return BlogCommentResponse{
		ID:        response.ID,
		BlogID:    response.BlogID,
		AuthorID:  response.AuthorID,
		Comment:   response.Comment,
		CreatedAt: response.CreatedAt,
	}
}
