package repositories

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database"

	"g6/blog-api/Infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type RefreshTokenRepository struct {
	DB         mongo.Database
	Collection string
}

func NewRefreshTokenRepository(db mongo.Database, collection string) domain.IRefreshTokenRepository {
	return &RefreshTokenRepository{
		DB:         db,
		Collection: collection,
	}
}

func (repo *RefreshTokenRepository) Save(ctx context.Context, token *domain.RefreshToken) error {
	tokenDb := database.FromRefreshTokenEntityToDB(token)
	if _, err := repo.DB.Collection(repo.Collection).InsertOne(ctx, tokenDb); err != nil {
		return err
	}
	return nil
}

func (repo *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var tokenDB database.RefreshTokenDB
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"token": token}).Decode(&tokenDB)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}
	return database.FromRefreshTokenDBToEntity(&tokenDB), nil
}

func (repo *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := repo.DB.Collection(repo.Collection).DeleteOne(ctx, bson.M{"_id": userID})
	return err
}
