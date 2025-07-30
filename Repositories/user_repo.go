package repositories

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByUsernameOrEmail(ctx context.Context, key string) (domain.User, error)
	UpdateTokens(ctx context.Context, userID string, tp domain.TokenPair) error
	InvalidateTokens(ctx context.Context, userID string) error
	ChangeRole(ctx context.Context, userID string, newRole domain.Role) error
	//... more methods can be added based on the usecases

}


type userRepo struct {
	col *mongo.Collection
}


func NewUserRepo(col *mongo.Collection) UserRepository {
	return &userRepo{col: col}
}


func (r *userRepo) Create(ctx context.Context, user domain.User) error {
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	model, _ := mapper.UserFromDomain(&user)
	_, err := r.col.InsertOne(ctx, model)
	return err
}

func (r *userRepo) FindByUsernameOrEmail(ctx context.Context, key string) (domain.User, error) {
	var userModel mapper.UserModel
	filter := bson.M{"$or": []bson.M{
		{"username": key},
		{"email": key},
	}}
	err := r.col.FindOne(ctx, filter).Decode(&userModel)
	if err != nil{
		return domain.User{}, err
	}
	return *mapper.UserToDomain(&userModel), nil
}

func (r *userRepo) UpdateTokens(ctx context.Context, userID string, tp domain.TokenPair) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{	
			"tokens.access_token":  tp.AccessToken,
			"tokens.refresh_token": tp.RefreshToken,
		},
	})
	return err
}

func (r *userRepo) InvalidateTokens(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$unset": bson.M{
			"tokens.access_token":  "",
			"tokens.refresh_token": "",
		},
	})
	return err
}

func (r *userRepo) ChangeRole(ctx context.Context, userID string, role domain.Role) error {
	objectID, _ := primitive.ObjectIDFromHex(userID)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{		
			"role": string(role),
		},
	})		
	return err
}