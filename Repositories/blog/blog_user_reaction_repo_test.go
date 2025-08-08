package repository

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	mongo_mocks "g6/blog-api/Infrastructure/database/mongo/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func TestNewUserReactionRepo(t *testing.T) {
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollections := &mongo.Collections{}

	repo := NewUserReactionRepo(mockDB, mockCollections)

	assert.NotNil(t, repo, "Expected non-nil repository")
}

func TestBlogUserReactionRepo_Create_NewReaction_Success(t *testing.T) {
	ctx := context.TODO()

	blogID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockReactionCollection := mongo_mocks.NewMockCollection(t)
	mockBlogCollection := mongo_mocks.NewMockCollection(t)
	mockSingleResult := mongo_mocks.NewMockSingleResult(t)

	mockDB.On("Collection", "blog_user_reactions").Return(mockReactionCollection)
	mockDB.On("Collection", "blog_posts").Return(mockBlogCollection)

	mockReactionCollection.On("FindOne", ctx, mock.Anything).Return(mockSingleResult)
	mockSingleResult.On("Decode", mock.Anything).Return(mongodriver.ErrNoDocuments)

	insertedID := primitive.NewObjectID()
	mockReactionCollection.On("InsertOne", ctx, mock.Anything).Return(&mongodriver.InsertOneResult{InsertedID: insertedID}, nil)

	mockBlogCollection.On("UpdateOne", ctx, bson.M{"_id": blogID}, bson.M{"$inc": bson.M{"likes": 1}}).Return(&mongodriver.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions", BlogPosts: "blog_posts"})

	reaction := &domain.BlogUserReaction{
		BlogID: blogID.Hex(),
		UserID: userID.Hex(),
		IsLike: true,
	}

	result, err := repo.Create(ctx, reaction)

	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected a result")
	assert.Equal(t, insertedID.Hex(), result.ID)
	mockDB.AssertExpectations(t)
	mockReactionCollection.AssertExpectations(t)
	mockBlogCollection.AssertExpectations(t)
}

func TestBlogUserReactionRepo_Create_InvalidID(t *testing.T) {
	ctx := context.TODO()

	mockDB := mongo_mocks.NewMockDatabase(t)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions"})

	invalidReaction := &domain.BlogUserReaction{
		BlogID: "invalid-id",
		UserID: primitive.NewObjectID().Hex(),
		IsLike: true,
	}

	result, err := repo.Create(ctx, invalidReaction)

	assert.Nil(t, result, "Expected nil result")
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "invalid reaction input: invalid blog ID: the provided hex string is not a valid ObjectID", err.Err.Error())
	mockDB.AssertNotCalled(t, "Collection", mock.Anything)
}

func TestBlogUserReactionRepo_Create_DBFailure(t *testing.T) {
	ctx := context.TODO()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockReactionCollection := mongo_mocks.NewMockCollection(t)
	mockSingleResult := mongo_mocks.NewMockSingleResult(t)

	mockDB.On("Collection", "blog_user_reactions").Return(mockReactionCollection)
	mockDB.On("Collection", "blog_posts").Return(mongo_mocks.NewMockCollection(t))

	dbErr := errors.New("database connection failed")
	mockReactionCollection.On("FindOne", ctx, mock.Anything).Return(mockSingleResult)
	mockSingleResult.On("Decode", mock.Anything).Return(dbErr)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions", BlogPosts: "blog_posts"})

	reaction := &domain.BlogUserReaction{
		BlogID: primitive.NewObjectID().Hex(),
		UserID: primitive.NewObjectID().Hex(),
		IsLike: true,
	}

	result, err := repo.Create(ctx, reaction)

	assert.Nil(t, result, "Expected nil result")
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Contains(t, err.Err.Error(), "failed to check existing reaction")
	mockDB.AssertExpectations(t)
	mockReactionCollection.AssertExpectations(t)
}

func TestBlogUserReactionRepo_Delete_Success(t *testing.T) {
	ctx := context.TODO()

	reactionID := primitive.NewObjectID()
	blogID := primitive.NewObjectID()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockReactionCollection := mongo_mocks.NewMockCollection(t)
	mockBlogCollection := mongo_mocks.NewMockCollection(t)
	mockSingleResult := mongo_mocks.NewMockSingleResult(t)

	mockDB.On("Collection", "blog_user_reactions").Return(mockReactionCollection)
	mockDB.On("Collection", "blog_posts").Return(mockBlogCollection)

	foundReaction := mapper.BlogUserReactionModel{
		ID:     reactionID,
		BlogID: blogID,
		IsLike: true,
	}
	mockReactionCollection.On("FindOne", ctx, mock.Anything).Return(mockSingleResult)
	mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*mapper.BlogUserReactionModel)
		*arg = foundReaction
	})

	mockReactionCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(1), nil)

	mockBlogCollection.On("UpdateOne", ctx, bson.M{"_id": blogID}, bson.M{"$inc": bson.M{"likes": -1}}).Return(&mongodriver.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions", BlogPosts: "blog_posts"})

	err := repo.Delete(ctx, reactionID.Hex())

	assert.Nil(t, err, "Expected no error")
	mockDB.AssertExpectations(t)
	mockReactionCollection.AssertExpectations(t)
	mockBlogCollection.AssertExpectations(t)
}

func TestBlogUserReactionRepo_Delete_InvalidID(t *testing.T) {
	ctx := context.TODO()

	mockDB := mongo_mocks.NewMockDatabase(t)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions"})

	err := repo.Delete(ctx, "invalid-id")

	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Contains(t, err.Err.Error(), "invalid ObjectID")
	mockDB.AssertNotCalled(t, "Collection", mock.Anything)
}

func TestBlogUserReactionRepo_Delete_NotFound(t *testing.T) {
	ctx := context.TODO()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockReactionCollection := mongo_mocks.NewMockCollection(t)
	mockSingleResult := mongo_mocks.NewMockSingleResult(t)

	mockDB.On("Collection", "blog_user_reactions").Return(mockReactionCollection)

	mockReactionCollection.On("FindOne", ctx, mock.Anything).Return(mockSingleResult)
	mockSingleResult.On("Decode", mock.Anything).Return(mongodriver.ErrNoDocuments)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions"})

	err := repo.Delete(ctx, primitive.NewObjectID().Hex())

	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Contains(t, err.Err.Error(), "no reaction found")
	mockDB.AssertExpectations(t)
	mockReactionCollection.AssertExpectations(t)
}

func TestBlogUserReactionRepo_Delete_DBFailure(t *testing.T) {
	ctx := context.TODO()

	reactionID := primitive.NewObjectID()
	blogID := primitive.NewObjectID()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockReactionCollection := mongo_mocks.NewMockCollection(t)
	mockSingleResult := mongo_mocks.NewMockSingleResult(t)

	mockDB.On("Collection", "blog_user_reactions").Return(mockReactionCollection)

	foundReaction := mapper.BlogUserReactionModel{
		ID:     reactionID,
		BlogID: blogID,
		IsLike: true,
	}
	mockReactionCollection.On("FindOne", ctx, mock.Anything).Return(mockSingleResult)
	mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*mapper.BlogUserReactionModel)
		*arg = foundReaction
	})

	dbErr := errors.New("database connection lost")
	mockReactionCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(0), dbErr)

	repo := NewUserReactionRepo(mockDB, &mongo.Collections{BlogUserReactions: "blog_user_reactions"})

	err := repo.Delete(ctx, reactionID.Hex())

	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Contains(t, err.Err.Error(), "failed to delete reaction")
	mockDB.AssertExpectations(t)
	mockReactionCollection.AssertExpectations(t)
}
