package repository

import (
	"context"
	"errors"
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
	// Validate and convert domain to Mongo model
	commentModel := &mapper.BlogCommentModel{}
	err := commentModel.Parse(comment)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid comment: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	// Validate blog ID
	blogOID, err := primitive.ObjectIDFromHex(comment.BlogID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	// Ensure blog exists
	var blogModel mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": blogOID}).Decode(&blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("blog not found: %w", err),
			Code: http.StatusNotFound,
		}
	}

	// Set created time and insert
	commentModel.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	result, err := b.db.Collection(b.collections.BlogComments).InsertOne(ctx, commentModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to create comment: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Set the inserted ID back to the domain object
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to retrieve inserted comment ID"),
			Code: http.StatusInternalServerError,
		}
	}

	comment.ID = insertedID.Hex()
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
		if errors.Is(err, mongo.ErrNoDocuments()) {
			return nil, &domain.DomainError{
				Err:  fmt.Errorf("comment not found: %w", err),
				Code: 404,
			}
		}
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to retrieve comment: %w", err),
			Code: 500,
		}
	}

	commentModel.ID = oid // Ensure ID is set correctly

	return commentModel.ToDomain(), nil
}

// GetCommentsByBlogID implements domain.BlogCommentRepository.
func (b *blogCommentRepository) GetCommentsByBlogID(ctx context.Context, blogID string, limit int) ([]domain.BlogComment, *domain.DomainError) {
	// 1. Validate the blog ID
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: 400,
		}
	}

	// 2. Query the comment collection
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // recent comments first

	cursor, err := b.db.Collection(b.collections.BlogComments).Find(ctx, bson.M{"blog_id": oid}, opts)
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

	// Check if any comments were found
	if len(domainComments) == 0 {
		return nil, &domain.DomainError{
			Err:  errors.New("no comments found for the specified blog"),
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
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid comment ID: %w", err),
			Code: 400,
		}
	}

	// Prepare the update
	update := bson.M{
		"$set": bson.M{
			"comment": comment.Comment,
		},
	}

	// Perform the update on the correct collection
	_, err = b.db.Collection(b.collections.BlogComments).UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to update comment: %w", err),
			Code: 500,
		}
	}

	// Fetch and return the updated comment
	updatedComment, domainErr := b.GetCommentByID(ctx, id)
	if domainErr != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to fetch updated comment: %w", domainErr.Err),
			Code: 500,
		}
	}

	return updatedComment, nil
}


func NewBlogCommentRepository(db mongo.Database, collections *mongo.Collections) domain.BlogCommentRepository {
	return &blogCommentRepository{
		db:          db,
		collections: collections,
	}
}
