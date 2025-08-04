package repository

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogCommentRepository struct {
	db          mongo.Database
	collections *mongo.Collections
}

// Create implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Create(ctx context.Context, comment *domain.BlogComment) (*domain.BlogComment, *domain.DomainError) {
	// Validate the comment
	comment_mondel := &mapper.BlogCommentModel{}
	err := comment_mondel.Parse(comment)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Set CreatedAt field and insert the comment
	comment_mondel.CreatedAt = time.Now()
	inserted, err := b.db.Collection(b.collections.BlogComments).InsertOne(ctx, comment_mondel)

	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to create comment: %w", err),
			Code: 500,
		}
	}

	// Convert the inserted ID to string and return the comment
	comment.ID = inserted.InsertedID.(primitive.ObjectID).Hex()
	return comment, nil
}

// Delete implements domain.BlogCommentRepository.
func (b *blogCommentRepository) Delete(ctx context.Context, id string) *domain.DomainError {
	// Validate the comment ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("invalid comment ID: %w", err),
			Code: 400,
		}
	}

	// Delete the comment from the database
	filter := bson.M{"_id": oid}
	del_cnt, err := b.db.Collection(b.collections.BlogComments).DeleteOne(ctx, filter)

	// Handle errors and check if the comment was found
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to delete comment: %w", err),
			Code: 500,
		}
	}
	if del_cnt == 0 {
		return &domain.DomainError{
			Err:  fmt.Errorf("comment not found with ID: %s", id),
			Code: 404,
		}
	}

	return nil
}

// GetCommentByID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentByID(ctx context.Context, id string) (*domain.BlogComment, *domain.DomainError) {
	// Validate the comment ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid comment ID: %w", err),
			Code: 400,
		}
	}

	// Query the comment by ID
	var commentModel mapper.BlogCommentModel
	err = b.db.Collection(b.collections.BlogComments).FindOne(ctx, bson.M{"_id": oid}).Decode(&commentModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("comment not found: %w", err),
			Code: 404,
		}
	}

	// Convert the comment model to domain model
	return commentModel.ToDomain(), nil
}

// GetCommentsByBlogID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]domain.BlogComment, *domain.DomainError) {
	// Validate the blog ID
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: 400,
		}
	}

	// Check if the blog exists
	var blogModel mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("blog not found: %w", err),
			Code: 404,
		}
	}

	// Query the comment collection
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
		domainComments = append(domainComments, *comment.ToDomain())
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
func (b *blogCommentRepository) Update(ctx context.Context, id string, comment *domain.BlogComment) (*domain.BlogComment, *domain.DomainError) {
	// Validate the comment ID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.BlogComment{}, &domain.DomainError{
			Err:  fmt.Errorf("invalid comment ID: %w", err),
			Code: 400,
		}
	}

	// Validate the comment content
	update := bson.M{
		"$set": bson.M{
			"comment": comment.Comment,
		},
	}
	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, oid, update)
	if err != nil {
		return &domain.BlogComment{}, &domain.DomainError{
			Err:  fmt.Errorf("failed to update comment: %w", err),
			Code: 500,
		}
	}

	// Fetch the updated comment
	updatedComment, err1 := b.GetCommentByID(ctx, id)
	if err1 != nil {
		return &domain.BlogComment{}, &domain.DomainError{
			Err:  fmt.Errorf("failed to fetch updated comment: %w", err1.Err),
			Code: 500,
		}
	}

	// Return the updated comment
	return updatedComment, nil
}

func NewBlogCommentRepository(db mongo.Database, collections *mongo.Collections) domain.BlogCommentRepository {
	return &blogCommentRepository{
		db:          db,
		collections: collections,
	}
}
