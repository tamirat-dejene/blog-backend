package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" binding:"required"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	AuthorID  primitive.ObjectID `bson:"author_id" json:"author_id" binding:"required"`
	Tags      []string           `bson:"tags,omitempty" json:"tags"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Likes     int                `bson:"likes" json:"likes"`
	Dislikes  int                `bson:"dislikes" json:"dislikes"`
	ViewCount int                `bson:"view_count" json:"view_count"`
}

type BlogComment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id" binding:"required"`
	AuthorID  primitive.ObjectID `bson:"author_id" json:"author_id" binding:"required"`
	Comment   string             `bson:"comment" json:"comment" binding:"required"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
