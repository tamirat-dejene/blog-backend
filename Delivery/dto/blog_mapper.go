package dto

import "time"

type BlogRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	AuthorID string   `json:"author_id" binding:"required"` // string, will convert to ObjectID in domain
	Tags     []string `json:"tags"`
}

type BlogResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	ViewCount int       `json:"view_count"`
}
