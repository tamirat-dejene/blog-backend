package repository

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUserReactionRepo struct {
	db          mongo.Database
	collections *mongo.Collections
}

func NewUserReactionRepo(database mongo.Database, collections *mongo.Collections) domain.BlogUserReactionRepository {
	return &BlogUserReactionRepo{
		db:          database,
		collections: collections,
	}
}

func (u *BlogUserReactionRepo) Create(ctx context.Context, reaction *domain.BlogUserReaction) (*domain.BlogUserReaction, *domain.DomainError) {
	// Convert domain model to MongoDB model
	blogReaction := &mapper.BlogUserReactionModel{}
	if err := blogReaction.Parse(reaction); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid reaction input: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	// Prepare filters
	reactionFilter := bson.M{
		"blog_id": blogReaction.BlogID,
		"user_id": blogReaction.UserID,
	}
	blogFilter := bson.M{"_id": blogReaction.BlogID}
	blogCollection := u.db.Collection(u.collections.BlogPosts)

	// Check if a reaction already exists
	var existing mapper.BlogUserReactionModel
	err := u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, reactionFilter).Decode(&existing)

	if err == mongo.ErrNoDocuments() {
		// No existing reaction — insert new
		reaction.CreatedAt = time.Now()
		blogReaction.CreatedAt = reaction.CreatedAt

		inserted, err := u.db.Collection(u.collections.BlogUserReactions).InsertOne(ctx, blogReaction)
		if err != nil {
			return nil, &domain.DomainError{
				Err:  fmt.Errorf("failed to insert reaction: %w", err),
				Code: http.StatusInternalServerError,
			}
		}

		// Safely assert inserted ID
		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			blogReaction.ID = oid
		}

		// Update blog's like/dislike counter
		counterField := "likes"
		if !reaction.IsLike {
			counterField = "dislikes"
		}
		if _, err := blogCollection.UpdateOne(ctx, blogFilter, bson.M{"$inc": bson.M{counterField: 1}}); err != nil {
			return nil, &domain.DomainError{
				Err:  fmt.Errorf("failed to increment %s: %w", counterField, err),
				Code: http.StatusInternalServerError,
			}
		}

		return blogReaction.ToDomain(), nil
	} else if err != nil {
		// Unexpected DB error
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to check existing reaction: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Existing reaction found
	if existing.IsLike == reaction.IsLike {
		// Same reaction type — no update needed
		return existing.ToDomain(), nil
	}

	// Reaction type changed — update document
	update := bson.M{
		"$set": bson.M{
			"is_like":    reaction.IsLike,
			"created_at": time.Now(),
		},
	}
	if _, err := u.db.Collection(u.collections.BlogUserReactions).UpdateOne(ctx, reactionFilter, update); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to update user reaction: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Adjust blog's like/dislike counters
	adjust := bson.M{
		"likes":    0,
		"dislikes": 0,
	}
	if reaction.IsLike {
		adjust["likes"] = 1
		adjust["dislikes"] = -1
	} else {
		adjust["likes"] = -1
		adjust["dislikes"] = 1
	}

	if _, err := blogCollection.UpdateOne(ctx, blogFilter, bson.M{"$inc": adjust}); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to adjust blog reaction counters: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Update popularity score, find the blog again
	var blogModel mapper.BlogPostModel
	if err := blogCollection.FindOne(ctx, blogFilter).Decode(&blogModel); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to find blog for popularity score update: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	ps := utils.CalculatePopularityScore(blogModel.Likes, blogModel.ViewCount, blogModel.CommentCount, blogModel.Dislikes)
	if _, err := blogCollection.UpdateOne(ctx, blogFilter, bson.M{"$set": bson.M{"popularity_score": ps}}); err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to update blog popularity score: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Update the return object with new timestamp
	reaction.CreatedAt = time.Now()
	return reaction, nil
}

func (u *BlogUserReactionRepo) Delete(ctx context.Context, id string) *domain.DomainError {
	// Convert string ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("invalid ObjectID: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	// Find the reaction to delete
	var reaction mapper.BlogUserReactionModel
	err = u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, bson.M{"_id": oid}).Decode(&reaction)
	if err == mongo.ErrNoDocuments() {
		return &domain.DomainError{
			Err:  fmt.Errorf("no reaction found with ID %s", id),
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to find reaction: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Delete the reaction
	_, err = u.db.Collection(u.collections.BlogUserReactions).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to delete reaction: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Decrement the correct count
	decField := "likes"
	if !reaction.IsLike {
		decField = "dislikes"
	}

	// Update the blog's like/dislike counter
	_, err = u.db.Collection(u.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": reaction.BlogID}, bson.M{
		"$inc": bson.M{decField: -1},
	})
	
	if err != nil {
		return &domain.DomainError{
			Err:  fmt.Errorf("failed to decrement %s count: %w", decField, err),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func (u *BlogUserReactionRepo) GetUserReaction(ctx context.Context, blogID string, userID string) (*domain.BlogUserReaction, *domain.DomainError) {
	// Convert string IDs to ObjectIDs
	blogIDObj, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: http.StatusBadRequest,
		}
	}
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid user ID: %w", err),
			Code: http.StatusBadRequest,
		}
	}
	// Prepare the filter and query the database
	filter := bson.M{"blog_id": blogIDObj, "user_id": userIDObj}
	var reaction mapper.BlogUserReactionModel
	err = u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, filter).Decode(&reaction)

	// Handle errors
	if err == mongo.ErrNoDocuments() {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("no reaction found for blog %s by user %s", blogID, userID),
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to find user reaction: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	// Convert the MongoDB model back to the domain model
	return reaction.ToDomain(), nil
}
