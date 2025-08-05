package repository

import (
	"context"
	"errors"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	mongo_mocks "g6/blog-api/Infrastructure/database/mongo/mocks"
	"net/http"
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

// Test for Create method in blog post repository
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

	assert.Nil(t, err, "Expected no error on create")
	assert.NotNil(t, result)
	assert.Equal(t, insertedID.Hex(), result.ID, "Expected inserted ID to match")
	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}

// Passing invalid author id in blog post creation
func TestBlogPostRepo_Create_Error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	// Prepare invalid blog post with malformed ObjectID string
	invalidBlog := &domain.BlogPost{
		Title:    "Invalid Author ID",
		Content:  "Content",
		AuthorID: "1", // This should trigger Parse() error
	}

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	// Act
	result, err := repo.Create(ctx, invalidBlog)

	// Assert
	assert.Nil(t, result, "Expected result to be nil")
	assert.NotNil(t, err, "Expected error due to invalid AuthorID")
	assert.Contains(t, err.Err.Error(), "invalid author ID") // or your custom error message

	// Ensure InsertOne is NOT called
	mockCollection.AssertNotCalled(t, "InsertOne", mock.Anything, mock.Anything)
	mockDB.AssertNotCalled(t, "Collection", blogPostsCollection)
}

// Passing valid author id in delete blog post
func TestBlogPostRepo_Delete_Success(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	mockDB.On("Collection", blogPostsCollection).Return(mockCollection)

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	blogID := primitive.NewObjectID().Hex()

	mockCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(1), nil)

	err := repo.Delete(ctx, blogID)

	assert.Nil(t, err, "Expected no error on delete")
	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}

// Passing invalid blog id in delete blog post
func TestBlogPostRepo_Delete_Error(t *testing.T) {
	t.Parallel()

	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	invalidBlogID := "1"
	_, expectedErr := primitive.ObjectIDFromHex(invalidBlogID)

	// Act
	err := repo.Delete(context.TODO(), invalidBlogID)

	// Assert
	assert.NotNil(t, err, "Expected error due to invalid ObjectID")
	assert.Equal(t, http.StatusBadRequest, err.Code, "Expected error code to be Bad Request")
	assert.Contains(t, err.Err.Error(), expectedErr.Error(), "Expected error message to indicate invalid ObjectID")

	// Ensure DeleteOne is NOT called
	mockCollection.AssertNotCalled(t, "DeleteOne", mock.Anything, mock.Anything)
	mockDB.AssertNotCalled(t, "Collection", blogPostsCollection)
}

// Test for Delete method in blog post repository with not found scenario
func TestBlogPostRepo_Delete_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	mockDB.On("Collection", blogPostsCollection).Return(mockCollection)

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	blogID := primitive.NewObjectID().Hex()

	mockCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(0), mongodriver.ErrNoDocuments)

	err := repo.Delete(ctx, blogID)

	assert.NotNil(t, err, "Expected error on delete")
	assert.Equal(t, http.StatusNotFound, err.Code, "Expected error code to be Not Found")
	assert.Contains(t, err.Err.Error(), "not found", "Expected error message to indicate not found")

	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}

// Test for Delete method in blog post repository with internal server error
func TestBlogPostRepo_Delete_InternalServerError(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogPostsCollection := "test_blog_posts"

	mockDB.On("Collection", blogPostsCollection).Return(mockCollection)

	repo := NewBlogPostRepo(mockDB, &mongo.Collections{BlogPosts: blogPostsCollection})

	blogID := primitive.NewObjectID().Hex()

	mockCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(0), errors.New("internal server error"))

	err := repo.Delete(ctx, blogID)

	assert.NotNil(t, err, "Expected error on delete")
	assert.Equal(t, http.StatusInternalServerError, err.Code, "Expected error code to be Internal Server Error")
	assert.Contains(t, err.Err.Error(), "internal server error", "Expected error message to indicate internal server error")

	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}

// Test for Get method in blog post repository with success scenario
func TestBlogPostRepo_Get_Success(t *testing.T) {
	
}