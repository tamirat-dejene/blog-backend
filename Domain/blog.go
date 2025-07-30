package domain

import "time"

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
