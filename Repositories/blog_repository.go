package repositories

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//	type Blog struct {
//		ID        primitive.ObjectID `bson:"_id,omitempty"`
//		Title     string             `bson:"title"`
//		Content   string             `bson:"content"`
//		AuthorID  primitive.ObjectID `bson:"author_id"`
//		Tags      []string           `bson:"tags"`
//		CreatedAt time.Time          `bson:"created_at"`
//		UpdatedAt time.Time          `bson:"updated_at"`
//		Likes     int                `bson:"likes"`
//		Dislikes  int                `bson:"dislikes"`
//		ViewCount int                `bson:"view_count"`
//	}
type BlogRepository struct {
	posts    *mongo.Collection
	comments *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) domain.BlogRepository {
	return &BlogRepository{
		posts:    db.Collection("posts"),
		comments: db.Collection("comments"),
	}
}

func (b *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	now := time.Now()
	blog.CreatedAt = now
	blog.UpdatedAt = now
	authorID, err := primitive.ObjectIDFromHex(blog.AuthorID)
	if err != nil {
		return errors.New("failed to convert author ID to ObjectID")
	}
	doc := bson.M{
		"title":      blog.Title,
		"content":    blog.Content,
		"author_id":  authorID,
		"tags":       blog.Tags,
		"created_at": blog.CreatedAt,
		"updated_at": blog.UpdatedAt,
		"likes":      0,
		"dislikes":   0,
		"view_count": 0,
	}
	res, err := b.posts.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to convert inserted ID to ObjectID")
	}
	blog.ID = oid.Hex()
	return nil
}
func (b *BlogRepository) Update(ctx context.Context, blog *domain.Blog) error
func (b *BlogRepository) Delete(ctx context.Context, id string) error
