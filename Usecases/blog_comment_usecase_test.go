package usecases

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	redis_mocks "g6/blog-api/Infrastructure/redis/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogCommentUsecaseSuite struct {
	suite.Suite
	blogCommentUsecase domain.BlogCommentUsecase
	Repo               *domain_mocks.MockBlogCommentRepository
	Redis              *redis_mocks.MockRedisClient
	Ctx                context.Context
	Comment            *domain.BlogComment
}

func (s *BlogCommentUsecaseSuite) SetupTest() {
	Comment := domain.BlogComment{
		AuthorID:  primitive.NewObjectID().Hex(),
		BlogID:    primitive.NewObjectID().Hex(),
		Comment:   "this is a comment",
		CreatedAt: time.Now(),
	}
	s.Comment = &Comment
	s.Repo = new(domain_mocks.MockBlogCommentRepository)
	s.Redis = new(redis_mocks.MockRedisClient)
	s.Ctx = context.Background()
	s.blogCommentUsecase = NewBlogCommentUsecase(s.Repo, s.Redis, time.Second*2)

}

func (s *BlogCommentUsecaseSuite) TestCommentUsecase_Create_Success() {

	s.Repo.On("Create", mock.Anything, s.Comment).Return(s.Comment, nil)

	result, err := s.blogCommentUsecase.CreateComment(s.Ctx, s.Comment)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(s.Comment, result)
	s.Repo.AssertExpectations(s.T())
}

func (s *BlogCommentUsecaseSuite) TestCommentUsecase_Create_Error() {

	expectedError := &domain.DomainError{
		Err:  fmt.Errorf("invalid comment: %w", errors.New("some db error")),
		Code: http.StatusBadRequest,
	}
	s.Repo.On("Create", mock.Anything, s.Comment).Return(nil, expectedError)

	result, err := s.blogCommentUsecase.CreateComment(s.Ctx, s.Comment)

	s.Nil(result)
	s.NotNil(err)
	s.Equal(err, expectedError)

	s.Repo.AssertExpectations(s.T())
}

func (s *BlogCommentUsecaseSuite) TestCommentUsecase_Update_Success() {
	s.Repo.On("Update", mock.Anything, "id", s.Comment).Return(s.Comment, nil)

	result, err := s.blogCommentUsecase.UpdateComment(s.Ctx, "id", s.Comment)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(result, s.Comment)

	s.Repo.AssertExpectations(s.T())
}

func (s *BlogCommentUsecaseSuite) TestBlogCommentUsecase_Update_Error() {
	expectedError := &domain.DomainError{
		Err:  fmt.Errorf("invalid id: %w", errors.New("bad request")),
		Code: http.StatusBadRequest,
	}
	s.Repo.On("Update", mock.Anything, "", s.Comment).Return(nil, expectedError)

	result, err := s.blogCommentUsecase.UpdateComment(s.Ctx, "", s.Comment)

	s.Nil(result)
	s.NotNil(err)
	s.Equal(err, expectedError)

	s.Repo.AssertExpectations(s.T())
}

func TestBlogCommentUsecaseSuite(t *testing.T) {
	suite.Run(t, new(BlogCommentUsecaseSuite))
}
