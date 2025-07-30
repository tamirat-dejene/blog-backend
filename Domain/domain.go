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

type BlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	Update(ctx context.Context, blog *Blog) error
	Delete(ctx context.Context, id string) error
	IncrementViewCount(ctx context.Context, id string) error
	Like(ctx context.Context, id string) error
	Dislike(ctx context.Context, id string) error
}
