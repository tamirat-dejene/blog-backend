package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	mongo_mocks "g6/blog-api/Infrastructure/database/mongo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func TestNewBlogPostRepo(t *testing.T) {
	t.Parallel()

	mock_database := mongo_mocks.NewMockDatabase(t)
	mock_collections := mongo.Collections{}

	repo := NewBlogPostRepo(mock_database, &mock_collections)
	assert.NotNil(t, repo, "Expected non-nil repository")
}

// --- mock BlogPostModel's Parse ---
// If Parse is a method on BlogPostModel that needs mocking, you'd have to refactor
// it for testability (e.g., dependency injection, or separating parsing logic).

func TestBlogPostRepo_Create_Success(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	mockDB.On("Collection", blogPostsCollection).Return(mockCollection)

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	blog := &domain.BlogPost{
		Title:    "Test Title",
		Content:  "Test Content",
		AuthorID: primitive.NewObjectID().Hex(),
	}

	insertedID := primitive.NewObjectID()

	mockCollection.On("InsertOne", ctx, mock.Anything).Return(&mongodriver.InsertOneResult{
		InsertedID: insertedID,
	}, nil)

	result, err := repo.Create(ctx, blog)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, insertedID.Hex(), result.ID)
	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}
