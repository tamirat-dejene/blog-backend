package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	mongo_mocks "g6/blog-api/Infrastructure/database/mongo/mocks"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func TestNewUserReactionRepo(t *testing.T) {
	t.Parallel()

	mock_db := mongo_mocks.NewMockDatabase(t)
	mock_collections := mongo.Collections{}

	repo := NewUserReactionRepo(mock_db, &mock_collections)

	assert.NotNil(t, repo, "Expected non-nil repository")
}

func TestUserReaction_Create_Success_First_Reaction(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock_db := new(mongo_mocks.MockDatabase)
	mock_reactions_collection := new(mongo_mocks.MockCollection)
	mock_posts_collection := new(mongo_mocks.MockCollection)

	BlogUserReactionCollection := "test_blog_user_reaction"
	BlogPostCollection := "test_blog_posts"

	mock_db.On("Collection", BlogUserReactionCollection).Return(mock_reactions_collection)
	mock_db.On("Collection", BlogPostCollection).Return(mock_posts_collection)

	mock_find_one_result := new(mongo_mocks.MockSingleResult)
	mock_find_one_result.On("Decode", mock.Anything).Return(mongodriver.ErrNoDocuments)

	userID, err := primitive.ObjectIDFromHex("6892fa4f6fb7d28124a96f68")
	assert.NoError(t, err)
	blogID, err := primitive.ObjectIDFromHex("6892fa4f6fb7d28124a96f74")
	assert.NoError(t, err)

	mock_reactions_collection.On("FindOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		if !ok {
			return false
		}
		return m["blog_id"] == blogID && m["user_id"] == userID
	})).Return(mock_find_one_result).Once()

	insertedID := primitive.NewObjectID()
	mock_reactions_collection.On("InsertOne", ctx, mock.Anything).Return(&mongodriver.InsertOneResult{
		InsertedID: insertedID,
	}, nil).Once()

	mock_posts_collection.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongodriver.UpdateResult{
		ModifiedCount: 1,
	}, nil).Once()

	repo := NewUserReactionRepo(mock_db, &mongo.Collections{
		BlogUserReactions: BlogUserReactionCollection,
		BlogPosts:         BlogPostCollection,
	})

	reaction := &domain.BlogUserReaction{
		UserID: "6892fa4f6fb7d28124a96f68",
		BlogID: "6892fa4f6fb7d28124a96f74",
		IsLike: true,
	}

	result, domainErr := repo.Create(ctx, reaction)
	assert.Nil(t, domainErr, "no domain error expected on create")
	assert.NotNil(t, reaction)
	assert.Equal(t, insertedID.Hex(), result.ID)
	assert.NotNil(t, result)

	mock_db.AssertExpectations(t)
	mock_reactions_collection.AssertExpectations(t)
	mock_posts_collection.AssertExpectations(t)
}

func TestUserReaction_Create_ExistingSameType(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock_db := new(mongo_mocks.MockDatabase)
	mock_reactions_collection := new(mongo_mocks.MockCollection)
	mock_posts_collection := new(mongo_mocks.MockCollection)

	BlogUserReactionCollection := "test_blog_user_reaction"
	BlogPostCollection := "test_blog_posts"

	mock_db.On("Collection", BlogUserReactionCollection).Return(mock_reactions_collection).Once()
	mock_db.On("Collection", BlogPostCollection).Return(mock_posts_collection).Once()

	userIDStr, _ := primitive.ObjectIDFromHex("6892fa4f6fb7d28124a96f68")
	blogIDStr, _ := primitive.ObjectIDFromHex("6892fa4f6fb7d28124a96f74")
	existingReactionID := primitive.NewObjectID()
	existingReaction := mapper.BlogUserReactionModel{
		ID:        existingReactionID,
		UserID:    userIDStr,
		BlogID:    blogIDStr,
		IsLike:    true,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	mock_find_one_result := new(mongo_mocks.MockSingleResult)
	mock_find_one_result.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*mapper.BlogUserReactionModel)
		*arg = existingReaction
	}).Return(nil).Once()

	mock_reactions_collection.On("FindOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		return ok && m["blog_id"] == blogIDStr && m["user_id"] == userIDStr
	})).Return(mock_find_one_result).Once()

	repo := NewUserReactionRepo(mock_db, &mongo.Collections{
		BlogUserReactions: BlogUserReactionCollection,
		BlogPosts:         BlogPostCollection,
	})

	newReactionAttempt := &domain.BlogUserReaction{
		UserID: userIDStr.Hex(),
		BlogID: blogIDStr.Hex(),
		IsLike: true,
	}

	result, domainErr := repo.Create(ctx, newReactionAttempt)
	assert.Nil(t, domainErr, "no domain error expected for existing same type reaction")
	assert.NotNil(t, result)
	assert.Equal(t, existingReactionID.Hex(), result.ID, "Returned ID should match existing reaction ID")

	mock_db.AssertExpectations(t)
	mock_reactions_collection.AssertExpectations(t)
	mock_posts_collection.AssertExpectations(t)
}

func TestUserReaction_Create_ExistingDifferentType(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock_db := new(mongo_mocks.MockDatabase)
	mock_reactions_collection := new(mongo_mocks.MockCollection)
	mock_posts_collection := new(mongo_mocks.MockCollection)

	BlogUserReactionCollection := "test_blog_user_reaction"
	BlogPostCollection := "test_blog_posts"

	mock_db.On("Collection", BlogUserReactionCollection).Return(mock_reactions_collection).Twice()
	mock_db.On("Collection", BlogPostCollection).Return(mock_posts_collection).Once()

	userIDStr := "6892fa4f6fb7d28124a96f68"
	blogIDStr := "6892fa4f6fb7d28124a96f74"

	userIDObjID, err := primitive.ObjectIDFromHex(userIDStr)
	assert.NoError(t, err)
	blogIDObjID, err := primitive.ObjectIDFromHex(blogIDStr)
	assert.NoError(t, err)

	blogObjectID, err := primitive.ObjectIDFromHex(blogIDStr)
	assert.NoError(t, err)

	existingReactionID := primitive.NewObjectID()
	existingReaction := mapper.BlogUserReactionModel{
		ID:        existingReactionID,
		UserID:    userIDObjID,
		BlogID:    blogIDObjID,
		IsLike:    true,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	initialBlogPost := mapper.BlogPostModel{
		ID:           blogObjectID,
		Likes:        5,
		Dislikes:     2,
		ViewCount:    100,
		CommentCount: 10,
	}

	mock_find_one_reaction_result := new(mongo_mocks.MockSingleResult)
	mock_find_one_reaction_result.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*mapper.BlogUserReactionModel)
		*arg = existingReaction
	}).Return(nil).Once()

	mock_reactions_collection.On("FindOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		if !ok {
			return false
		}
		return m["blog_id"] == blogIDObjID && m["user_id"] == userIDObjID
	})).Return(mock_find_one_reaction_result).Once()

	mock_reactions_collection.On("UpdateOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		return ok && m["blog_id"] == blogIDObjID && m["user_id"] == userIDObjID
	}), mock.MatchedBy(func(update interface{}) bool {
		u, ok := update.(bson.M)
		if !ok {
			return false
		}
		set, ok := u["$set"].(bson.M)
		return ok && set["is_like"] == false
	})).Return(&mongodriver.UpdateResult{ModifiedCount: 1}, nil).Once()

	mock_posts_collection.On("UpdateOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		return ok && m["_id"] == blogObjectID
	}), mock.MatchedBy(func(update interface{}) bool {
		u, ok := update.(bson.M)
		if !ok {
			return false
		}
		inc, ok := u["$inc"].(bson.M)
		return ok && inc["likes"] == -1 && inc["dislikes"] == 1
	})).Return(&mongodriver.UpdateResult{ModifiedCount: 1}, nil).Once()

	mock_find_one_blog_result := new(mongo_mocks.MockSingleResult)
	mock_find_one_blog_result.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*mapper.BlogPostModel)
		*arg = mapper.BlogPostModel{
			ID:           blogObjectID,
			Likes:        initialBlogPost.Likes - 1,
			Dislikes:     initialBlogPost.Dislikes + 1,
			ViewCount:    initialBlogPost.ViewCount,
			CommentCount: initialBlogPost.CommentCount,
		}
	}).Return(nil).Once()

	mock_posts_collection.On("FindOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		return ok && m["_id"] == blogObjectID
	})).Return(mock_find_one_blog_result).Once()

	mock_posts_collection.On("UpdateOne", ctx, mock.MatchedBy(func(filter interface{}) bool {
		m, ok := filter.(bson.M)
		return ok && m["_id"] == blogObjectID
	}), mock.MatchedBy(func(update interface{}) bool {
		u, ok := update.(bson.M)
		if !ok {
			return false
		}
		set, ok := u["$set"].(bson.M)
		expectedPS := utils.CalculatePopularityScore(initialBlogPost.Likes-1, initialBlogPost.ViewCount, initialBlogPost.CommentCount, initialBlogPost.Dislikes+1)
		return ok && set["popularity_score"] == expectedPS
	})).Return(&mongodriver.UpdateResult{ModifiedCount: 1}, nil).Once()

	repo := NewUserReactionRepo(mock_db, &mongo.Collections{
		BlogUserReactions: BlogUserReactionCollection,
		BlogPosts:         BlogPostCollection,
	})

	newReactionAttempt := &domain.BlogUserReaction{
		UserID: userIDStr,
		BlogID: blogIDStr,
		IsLike: false,
	}

	result, domainErr := repo.Create(ctx, newReactionAttempt)
	assert.Nil(t, domainErr, "no domain error expected for existing different type reaction")
	assert.NotNil(t, result)
	assert.Equal(t, existingReactionID.Hex(), result.ID, "Returned ID should match existing reaction ID")
	assert.Equal(t, false, result.IsLike, "IsLike should be updated to false in the returned result")

	mock_db.AssertExpectations(t)
	mock_reactions_collection.AssertExpectations(t)
	mock_posts_collection.AssertExpectations(t)
}
