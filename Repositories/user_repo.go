package repositories

import (
	"context"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB         *mongo.Database
	Collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.IUserRepository {
	return &UserRepository{
		DB:         db,
		Collection: collection,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	fmt.Println("Creating user:", user)
	userDb := database.FromUserEntityToDB(user)
	if _, err := repo.DB.Collection(repo.Collection).InsertOne(ctx, userDb); err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*database.UserDB
	cursor, err := repo.DB.Collection(repo.Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user database.UserDB
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return database.FromUserDBListToEntityList(users), nil
}

func (repo *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userDB database.UserDB
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"username": bson.M{"$regex": username, "$options": "i"}}).Decode(&userDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return database.FromUserDBToEntity(&userDB), nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userDB database.UserDB
	err := repo.DB.Collection(repo.Collection).FindOne(ctx, bson.M{"email": bson.M{"$regex": email, "$options": "i"}}).Decode(&userDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return database.FromUserDBToEntity(&userDB), nil
}
