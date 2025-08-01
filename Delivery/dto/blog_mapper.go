package dto

import "time"

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
