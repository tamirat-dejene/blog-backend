package repositories

import (
	"context"
	"errors"
	"testing"

	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	mocks "g6/blog-api/Infrastructure/database/mongo/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

func newRepoWithMocks() (*UserRepository, *mocks.MockDatabase, *mocks.MockCollection) {
	mockDB := new(mocks.MockDatabase)
	mockColl := new(mocks.MockCollection)
	mockDB.On("Collection", "users").Return(mockColl)
	repo := &UserRepository{
		DB:         mockDB,
		Collection: "users",
	}
	return repo, mockDB, mockColl
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, mockDB, mockColl := newRepoWithMocks()
		mockColl.
			On("InsertOne", ctx, mock.AnythingOfType("*mapper.UserModel")).
			Return(&mongo.InsertOneResult{}, nil)

		user := &domain.User{Username: "test"}
		err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockColl.AssertExpectations(t)
	})

	t.Run("insert error", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockColl.
			On("InsertOne", ctx, mock.AnythingOfType("*mapper.UserModel")).
			Return(nil, errors.New("insert failed"))

		err := repo.CreateUser(ctx, &domain.User{})
		assert.EqualError(t, err, "insert failed")
	})
}

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockCursor := new(mocks.MockCursor)

		mockColl.On("Find", ctx, mock.Anything).Return(mockCursor, nil)
		mockCursor.On("Next", ctx).Return(true).Once()
		mockCursor.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Run(func(args mock.Arguments) {
			um := &mapper.UserModel{Username: "u1"}
			*(args.Get(0).(**mapper.UserModel)) = um
		}).Return(nil).Once()
		mockCursor.On("Next", ctx).Return(false)
		mockCursor.On("Err").Return(nil)
		mockCursor.On("Close", ctx).Return(nil)

		users, err := repo.GetAllUsers(ctx)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "u1", users[0].Username)
	})

	t.Run("find error", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockColl.On("Find", ctx, mock.Anything).Return(nil, errors.New("find failed"))
		users, err := repo.GetAllUsers(ctx)
		assert.Nil(t, users)
		assert.EqualError(t, err, "find failed")
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	repo, _, mockColl := newRepoWithMocks()

	mockColl.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	err := repo.UpdateUser(ctx, "60c72b2f9b1d8b3a0c8b4567", &domain.User{Username: "updated"})
	assert.NoError(t, err)
}

func TestFindUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockSR := new(mocks.MockSingleResult)
		mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)
		mockSR.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Run(func(args mock.Arguments) {
			um := &mapper.UserModel{Username: "u1"}
			*(args.Get(0).(**mapper.UserModel)) = um
		}).Return(nil)

		user, err := repo.FindUserByID(ctx, "60c72b2f9b1d8b3a0c8b4567")
		assert.NoError(t, err)
		assert.Equal(t, "u1", user.Username)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo, _, _ := newRepoWithMocks()
		user, err := repo.FindUserByID(ctx, "invalid")
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockSR := new(mocks.MockSingleResult)
		mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)
		mockSR.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Return(mongo.ErrNoDocuments)

		user, err := repo.FindUserByID(ctx, "60c72b2f9b1d8b3a0c8b4567")
		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})
}

func TestGetUserByUsername(t *testing.T) {
	ctx := context.Background()
	repo, _, mockColl := newRepoWithMocks()
	mockSR := new(mocks.MockSingleResult)
	mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)

	// This Run func sets the value of the pointer passed to Decode
	mockSR.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Run(func(args mock.Arguments) {
		// Get the pointer to *mapper.UserModel passed to Decode
		ptr := args.Get(0).(**mapper.UserModel)
		// Set it to a valid struct
		*ptr = &mapper.UserModel{Username: "user1"}
	}).Return(nil)

	user, err := repo.GetUserByUsername(ctx, "user1")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Username)
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockSR := new(mocks.MockSingleResult)
		mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)
		mockSR.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Run(func(args mock.Arguments) {
			um := &mapper.UserModel{Email: "email@example.com"}
			*(args.Get(0).(**mapper.UserModel)) = um
		}).Return(nil)
		_, err := repo.GetUserByEmail(ctx, "email@example.com")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockSR := new(mocks.MockSingleResult)
		mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)
		mockSR.On("Decode", mock.AnythingOfType("**mapper.UserModel")).Return(mongo.ErrNoDocuments)
		user, err := repo.GetUserByEmail(ctx, "missing@example.com")
		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})
}

func TestFindByUsernameOrEmail(t *testing.T) {
	ctx := context.Background()
	repo, _, mockColl := newRepoWithMocks()
	mockSR := new(mocks.MockSingleResult)
	mockColl.On("FindOne", ctx, mock.Anything).Return(mockSR)
	mockSR.On("Decode", mock.AnythingOfType("*mapper.UserModel")).Return(nil)

	_, err := repo.FindByUsernameOrEmail(ctx, "test")
	assert.NoError(t, err)
}

func TestInvalidateTokens(t *testing.T) {
	ctx := context.Background()
	repo, _, mockColl := newRepoWithMocks()
	mockColl.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	err := repo.InvalidateTokens(ctx, "some-id")
	assert.NoError(t, err)
}

func TestChangeRole(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, _, mockColl := newRepoWithMocks()
		mockColl.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
		err := repo.ChangeRole(ctx, "60c72b2f9b1d8b3a0c8b4567", "admin", "newusername")
		assert.NoError(t, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo, _, _ := newRepoWithMocks()
		err := repo.ChangeRole(ctx, "invalid", "admin", "newusername")
		assert.Error(t, err)
	})
}
