package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogUserReactionRepo struct {
	db          mongo.Database
	collections *collections
}

func NewUserReactionRepo(database mongo.Database, collections *collections) domain.BlogUserReactionRepository {
	return &BlogUserReactionRepo{
		db:          database,
		collections: collections,
	}
}

func (u *BlogUserReactionRepo) Create(ctx context.Context, reaction domain.BlogUserReaction) (domain.BlogUserReaction, error) {
	// Convert domain model to Mongo model
	blogReaction, err := mapper.BlogUserReactionFromDomain(&reaction)

	if err != nil {
		return domain.BlogUserReaction{}, err
	}
	filter := bson.M{
		"blog_id": blogReaction.BlogID,
		"user_id": blogReaction.UserID,
	}

	var existing mapper.BlogUserReactionModel

	// Check for existing reaction
	duplicateErr := u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, filter).Decode(&existing)

	blogCollection := u.db.Collection(u.collections.BlogPosts)
	blogFilter := bson.M{"_id": blogReaction.BlogID}

	if duplicateErr == mongo.ErrNoDocuments() {

		// No existing record, insert new
		reaction.CreatedAt = time.Now()
		_, err = u.db.Collection(u.collections.BlogUserReactions).InsertOne(ctx, blogReaction)

		if err != nil {
			return domain.BlogUserReaction{}, nil
		}

		//Increment like or dislike
		counterField := "likes"
		if !reaction.IsLike {
			counterField = "dislikes"
		}
		_, _ = blogCollection.UpdateOne(ctx, blogFilter, bson.M{
			"$inc": bson.M{counterField: 1},
		})

		response := mapper.BlogUserReactionToDomain(blogReaction)

		return *response, nil
	} else if duplicateErr != nil {

		// Other errors when fetching
		return domain.BlogUserReaction{}, err
	}

	// Reaction exists: if same reaction type, just return existing
	if existing.IsLike == reaction.IsLike {
		response := mapper.BlogUserReactionToDomain(&existing)

		return *response, nil
	}

	// Different reaction: update it
	update := bson.M{
		"$set": bson.M{
			"is_like":    reaction.IsLike,
			"created_at": time.Now(),
		},
	}

	_, err = u.db.Collection(u.collections.BlogUserReactions).UpdateOne(ctx, filter, update)

	if err != nil {
		return domain.BlogUserReaction{}, err
	}

	// Decrement the old increment the new
	inc := bson.M{
		"likes":    0,
		"dislikes": 0,
	}
	if reaction.IsLike {
		inc["likes"] = 1
		inc["dislikes"] = -1
	} else {
		inc["likes"] = -1
		inc["dislikes"] = 1
	}
	_, _ = blogCollection.UpdateOne(ctx, blogFilter, bson.M{
		"$inc": inc,
	})
	// Update local copy's timestamp for return
	reaction.CreatedAt = time.Now()
	return reaction, nil
}

func (u *BlogUserReactionRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var reaction mapper.BlogUserReactionModel
	err = u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, bson.M{"_id": oid}).Decode(&reaction)
	if err != nil {
		return err
	}
	_, err = u.db.Collection(u.collections.BlogUserReactions).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	//decrement the correct count
	decField := "likes"
	if !reaction.IsLike {
		decField = "dislikes"
	}
	_, err = u.db.Collection(u.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": reaction.BlogID}, bson.M{
		"$inc": bson.M{decField: -1},
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *BlogUserReactionRepo) GetUserReaction(ctx context.Context, blogID string, userID string) (domain.BlogUserReaction, error) {

	filter := bson.M{"blog_id": blogID, "user_id": userID}

	var reaction mapper.BlogUserReactionModel

	err := u.db.Collection(u.collections.BlogUserReactions).FindOne(ctx, filter).Decode(&reaction)

	if err != nil {
		return domain.BlogUserReaction{}, err
	}

	response := mapper.BlogUserReactionToDomain(&reaction)
	return *response, nil
}
