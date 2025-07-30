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
func BlogToDomain(blog *BlogModel) *domain.Blog {
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
func BlogFromDomain(blog *domain.Blog) (*BlogModel, error) {
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

// Convert BlogCommentModel to domain.BlogComment
func BlogCommentToDomain(comment *BlogCommentModel) *domain.BlogComment {
	return &domain.BlogComment{
		ID:        comment.ID.Hex(),
		BlogID:    comment.BlogID.Hex(),
		AuthorID:  comment.AuthorID.Hex(),
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
	}
}

// Convert domain.BlogComment to BlogCommentModel
func BlogCommentFromDomain(comment *domain.BlogComment) (*BlogCommentModel, error) {
	blogID, err := primitive.ObjectIDFromHex(comment.BlogID)
	if err != nil {
		return nil, err
	}
	authorID, err := primitive.ObjectIDFromHex(comment.AuthorID)
	if err != nil {
		return nil, err
	}

	var objectID primitive.ObjectID
	if comment.ID != "" {
		objectID, err = primitive.ObjectIDFromHex(comment.ID)
		if err != nil {
			return nil, err
		}
	} else {
		objectID = primitive.NewObjectID()
	}

	return &BlogCommentModel{
		ID:        objectID,
		BlogID:    blogID,
		AuthorID:  authorID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
	}, nil
}
