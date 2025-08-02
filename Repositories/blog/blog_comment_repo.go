package repository

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogCommentRepository struct {
	db          mongo.Database
	collections *mongo.Collections
}

func (b *blogCommentRepository) Create(ctx context.Context, comment domain.BlogComment) (*domain.BlogComment, error) {
	blogComment, err := mapper.BlogCommentFromDomain(&comment)

	if err != nil {
		return &domain.BlogComment{}, err
	}
	_, err = b.db.Collection(b.collections.BlogComments).InsertOne(ctx, blogComment)

	if err != nil {
		return &domain.BlogComment{}, err
	}
	res := mapper.BlogCommentToDomain(blogComment)
	return res, nil
}

func (b *blogCommentRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = b.db.Collection(b.collections.BlogComments).DeleteOne(ctx, filter)

	return err
}

// GetCommentByID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid comment ID: %w", err),
			Code: 400,
		}
	}

	var commentModel mapper.BlogCommentModel
	err = b.db.Collection(b.collections.BlogComments).FindOne(ctx, bson.M{"_id": oid}).Decode(&commentModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("comment not found: %w", err),
			Code: 404,
		}
	}

	return mapper.BlogCommentToDomain(&commentModel), nil
}

// GetCommentsByBlogID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]domain.BlogComment, *domain.DomainError) {
	fmt.Println("Fetching comments for blog ID:", blogID, b.collections.BlogPosts)
	// 1. Validate if the blog exists
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: 400,
		}
	}
	
	var blogModel mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("blog not found: %w", err),
			Code: 404,
		}
	}

	// 2. Query the comment collection
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // recent comments first

	cursor, err := b.db.Collection(b.collections.BlogComments).Find(ctx, bson.M{"blog_id": blogID}, opts)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to retrieve comments: %w", err),
			Code: 500,
		}
	}
	defer cursor.Close(ctx)

	// 3. Parse results
	var comments []mapper.BlogCommentModel
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to decode comments: %w", err),
			Code: 500,
		}
	}
	// 4. Convert to domain model
	var domainComments []domain.BlogComment
	for _, comment := range comments {
		domainComments = append(domainComments, *mapper.BlogCommentToDomain(&comment))
	}

	if len(domainComments) == 0 {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("no comments found for blog ID: %s", blogID),
			Code: 404,
		}
	}
	return domainComments, nil
}

// Update implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Update(ctx context.Context, id string, comment domain.BlogComment) (*domain.BlogComment, error) {
	panic("unimplemented")
}
func NewBlogCommentRepository(db mongo.Database, collections *mongo.Collections) domain.BlogCommentRepository {
	return &blogCommentRepository{
		db:          db,
		collections: collections,
	}
}
