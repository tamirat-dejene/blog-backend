package repository

import (
	"context"
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

func TestNewCommentRepo(t *testing.T) {
	t.Parallel()

	mock_db := mongo_mocks.NewMockDatabase(t)
	mock_collections := mongo.Collections{}

	repo := NewBlogCommentRepository(mock_db, &mock_collections)

	assert.NotNil(t, repo, "Expected non-nil repository")
}

func TestComment_Create_Success(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogCommentCollection := "test_blog_comment"

	repo := NewBlogCommentRepository(mockDB, &mongo.Collections{BlogComments: blogCommentCollection})
	mockDB.On("Collection", blogCommentCollection).Return(mockCollection)

	comment := &domain.BlogComment{
		BlogID:   primitive.NewObjectID().Hex(),
		AuthorID: primitive.NewObjectID().Hex(),
		Comment:  "this is a comment",
	}
	insertedID := primitive.NewObjectID()

	mockCollection.On("InsertOne", ctx, mock.Anything).Return(&mongodriver.InsertOneResult{
		InsertedID: insertedID,
	}, nil)

	result, err := repo.Create(ctx, comment)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.ID, insertedID.Hex())

	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)

}

func TestCommentRepo_Create_Error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogCommentCollection := "test_blog_comment"

	repo := NewBlogCommentRepository(mockDB, &mongo.Collections{BlogComments: blogCommentCollection})

	comment := &domain.BlogComment{
		BlogID:   primitive.NewObjectID().Hex(),
		AuthorID: "12",
		Comment:  "this is a comment",
	}

	result, err := repo.Create(ctx, comment)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Err.Error(), "invalid author ID")

	mockCollection.AssertExpectations(t)

	mockDB.AssertNotCalled(t, "Collection", blogCommentCollection)
	mockCollection.AssertNotCalled(t, "InsertOne", mock.Anything, mock.Anything)
}

func TestCommentRepo_Delete_Success(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogCommentCollection := "test_blog_comment"

	repo := NewBlogCommentRepository(mockDB, &mongo.Collections{BlogComments: blogCommentCollection})
	mockDB.On("Collection", blogCommentCollection).Return(mockCollection)

	id := primitive.NewObjectID().Hex()

	mockCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(1), nil)

	err := repo.Delete(ctx, id)

	assert.Nil(t, err)

	mockCollection.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestCommentRepo_Delete_Error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogCommentCollection := "test_blog_comment"

	invalidId := "1"

	repo := NewBlogCommentRepository(mockDB, &mongo.Collections{BlogComments: blogCommentCollection})

	_, ExpectedErr := primitive.ObjectIDFromHex(invalidId)

	err := repo.Delete(ctx, invalidId)

	assert.NotNil(t, err, "Expected error due to invalid ObjectID")
	assert.Equal(t, http.StatusBadRequest, err.Code, "Expected error code to be Bad Request")
	assert.Contains(t, err.Err.Error(), ExpectedErr.Error(), "Expected error message to indicate invalid ObjectID")

	mockCollection.AssertNotCalled(t, "DeleteOne", mock.Anything, mock.Anything)
	mockDB.AssertNotCalled(t, "Collection", blogCommentCollection)
}

func TestCommentRepo_Delete_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mockDB := mongo_mocks.NewMockDatabase(t)
	mockCollection := mongo_mocks.NewMockCollection(t)
	blogCommentCollection := "test_blog_comments"

	mockDB.On("Collection", blogCommentCollection).Return(mockCollection)

	repo := NewBlogCommentRepository(mockDB, &mongo.Collections{BlogComments: blogCommentCollection})

	id := primitive.NewObjectID().Hex()

	mockCollection.On("DeleteOne", ctx, mock.Anything).Return(int64(0), nil)

	err := repo.Delete(ctx, id)

	assert.NotNil(t, err, "Expected error on delete")
	assert.Equal(t, http.StatusNotFound, err.Code, "Expected error code to be Not Found")
	assert.Contains(t, err.Err.Error(), "comment not found", "Expected error message to indicate not found")

	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}
