package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogUserReactionRepo struct {
	db         mongo.Database
	collection string
}

func NewUserReactionRepo(database mongo.Database, collection string) domain.BlogUserReactionRepository {
	return &BlogUserReactionRepo{
		db:         database,
		collection: collection,
	}
}

func (u *BlogUserReactionRepo) Create(ctx context.Context, reaction domain.BlogUserReaction) (domain.BlogUserReaction, error) {
	blogReaction, err := mapper.BlogUserReactionFromDomain(&reaction)

	if err != nil {
		return domain.BlogUserReaction{}, err
	}
	filter := bson.M{
		"blog_id": blogReaction.BlogID,
		"user_id": blogReaction.UserID,
	}
	update := bson.M{
		"$set": bson.M{
			"is_like":    reaction.IsLike,
			"created_at": time.Now(),
		},
	}
	opt := options.Update().SetUpsert(true)

	_, err = u.db.Collection(u.collection).UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return domain.BlogUserReaction{}, err
	}
	res := mapper.BlogUserReactionToDomain(blogReaction)

	return *res, nil

}

func (u *BlogUserReactionRepo) Delete(ctx context.Context, id string) error {
	// to be implemented
	return nil
}
func (u *BlogUserReactionRepo) GetUserReaction(ctx context.Context, blogID, userID string) (domain.BlogUserReaction, error) {
	// to be implemented
	return domain.BlogUserReaction{}, nil
}
