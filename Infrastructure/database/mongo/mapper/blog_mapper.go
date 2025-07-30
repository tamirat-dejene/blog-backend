package mapper

import (
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	AuthorID  primitive.ObjectID `bson:"author_id"`
	Tags      []string           `bson:"tags,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Likes     int                `bson:"likes"`
	Dislikes  int                `bson:"dislikes"`
	ViewCount int                `bson:"view_count"`
}

type BlogCommentModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id"`
	AuthorID  primitive.ObjectID `bson:"author_id"`
	Comment   string             `bson:"comment"`
	CreatedAt time.Time          `bson:"created_at"`
}

// Convert to domain
func ToDomain(blog *BlogModel) *domain.Blog {
	return &domain.Blog{
		ID:        blog.ID.Hex(),
		Title:     blog.Title,
		Content:   blog.Content,
		AuthorID:  blog.AuthorID.Hex(),
		Tags:      blog.Tags,
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
		Likes:     blog.Likes,
		Dislikes:  blog.Dislikes,
		ViewCount: blog.ViewCount,
	}
}

// Convert from domain
func FromDomain(blog *domain.Blog) (*BlogModel, error) {
	authorID, err := primitive.ObjectIDFromHex(blog.AuthorID)
	if err != nil {
		return nil, err
	}

	var objectID primitive.ObjectID
	if blog.ID != "" {
		objectID, err = primitive.ObjectIDFromHex(blog.ID)
		if err != nil {
			return nil, err
		}
	} else {
		objectID = primitive.NewObjectID()
	}

	return &BlogModel{
		ID:        objectID,
		Title:     blog.Title,
		Content:   blog.Content,
		AuthorID:  authorID,
		Tags:      blog.Tags,
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
		Likes:     blog.Likes,
		Dislikes:  blog.Dislikes,
		ViewCount: blog.ViewCount,
	}, nil
}
