package repositories

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
)

type PasswordResetRepository struct {
	DB         mongo.Database
	Collection string
}

func NewPasswordResetRepository(db mongo.Database, col string) domain.IPasswordResetRepository {
	return &PasswordResetRepository{
		DB:         db,
		Collection: col,
	}
}

func (r *PasswordResetRepository) SaveResetToken(ctx context.Context, token *domain.PasswordResetToken) error {
	tokenModel := mapper.PasswordResetTokenFromDomain(token)
	_, err := r.DB.Collection(r.Collection).InsertOne(ctx, tokenModel)
	if err != nil {
		return err
	}
	return nil
}

func (r *PasswordResetRepository) FindByEmail(ctx context.Context, email string) (*domain.PasswordResetToken, error) {
	var tokenModel mapper.PasswordResetToken
	err := r.DB.Collection(r.Collection).FindOne(ctx, bson.M{"email": email}).Decode(&tokenModel)
	if err != nil {
		return nil, err
	}
	return mapper.PasswordResetTokenToDomain(&tokenModel), nil
}

func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, token *domain.PasswordResetToken) error {
	_, err := r.DB.Collection(r.Collection).UpdateOne(ctx, bson.M{"email": token.Email, "token_hash": token.TokenHash}, bson.M{"$set": bson.M{"used": true}})
	if err != nil {
		return err
	}
	return nil
}
func (r *PasswordResetRepository) DeleteByEmail(ctx context.Context, email string) error {
	_, err := r.DB.Collection(r.Collection).DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	return nil
}
