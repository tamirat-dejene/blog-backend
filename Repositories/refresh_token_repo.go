package repositories

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"

	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

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
	tokenDb := mapper.FromRefreshTokenEntityToDB(token)
	if _, err := repo.DB.Collection(repo.Collection).InsertOne(ctx, tokenDb); err != nil {
		return err
	}
	return nil
}

func (repo *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var tokenDB mapper.RefreshTokenDB
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"token": token}).Decode(&tokenDB)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}
	return mapper.FromRefreshTokenDBToEntity(&tokenDB), nil
}

func (repo *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	deleteCount, err := repo.DB.Collection(repo.Collection).DeleteOne(ctx, bson.M{"user_id": userID})
	if deleteCount == 0 {
		return fmt.Errorf("no refresh token found for user ID: %s", userID)
	}
	return err
}

func (repo *RefreshTokenRepository) ReplaceTokenByUserID(ctx context.Context, token *domain.RefreshToken) error {
	tokenDB := mapper.FromRefreshTokenEntityToDB(token)
	_, err := repo.DB.Collection(repo.Collection).UpdateOne(
		ctx,
		bson.M{"user_id": token.UserID},
		bson.M{"$set": tokenDB},
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// revoke refresh token
func (repo *RefreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	_, err := repo.DB.Collection(repo.Collection).UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"revoked": true}},
	)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// find token by user id
func (repo *RefreshTokenRepository) FindTokenByUserID(ctx context.Context, userID string) (*domain.RefreshToken, error) {
	var tokenDB mapper.RefreshTokenDB
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"user_id": userID}).Decode(&tokenDB)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}
	return mapper.FromRefreshTokenDBToEntity(&tokenDB), nil
}
