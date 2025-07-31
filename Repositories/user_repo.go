package repositories

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"g6/blog-api/Infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	DB         mongo.Database
	Collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.IUserRepository {
	return &UserRepository{
		DB:         db,
		Collection: collection,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	usermodel := mapper.UserFromDomain(user)
	_, err := repo.DB.Collection(repo.Collection).InsertOne(ctx, usermodel)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*mapper.UserModel
	cursor, err := repo.DB.Collection(repo.Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user *mapper.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return mapper.UserToDomainList(users), nil
}

func (repo *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userModel *mapper.UserModel
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"username": bson.M{"$regex": username, "$options": "i"}}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return mapper.UserToDomain(userModel), nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userModel *mapper.UserModel
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"email": bson.M{"$regex": email, "$options": "i"}}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments() {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return mapper.UserToDomain(userModel), nil
}

func (repo *UserRepository) FindByUsernameOrEmail(ctx context.Context, key string) (domain.User, error) {
	var userModel mapper.UserModel
	filter := bson.M{"$or": []bson.M{
		{"username": key},
		{"email": key},
	}}
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, filter).Decode(&userModel)
	if err != nil {
		return domain.User{}, err
	}
	return *mapper.UserToDomain(&userModel), nil
}

func (repo *UserRepository) InvalidateTokens(ctx context.Context, userID string) error {
	_, err := repo.DB.Collection(repo.Collection).UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"tokens": []string{}}})
	return err
}

func (repo *UserRepository) ChangeRole(ctx context.Context, targetUserID string, role string, username string) error {
	_, err := repo.DB.Collection(repo.Collection).UpdateOne(ctx, bson.M{"_id": targetUserID}, bson.M{
		"$set": bson.M{
			"role":       role,
			"updated_at": time.Now(),
			"username":   username,
		},
	})
	return err
}
