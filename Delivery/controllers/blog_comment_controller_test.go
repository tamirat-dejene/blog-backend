package controllers

import (
	"errors"
	"g6/blog-api/Delivery/bootstrap"
	domain "g6/blog-api/Domain"
	domain_mocks "g6/blog-api/Domain/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BlogCommentControllerSuite struct {
	suite.Suite
	BlogCommentUsecase *domain_mocks.MockBlogCommentUsecase
	Controller         *BlogCommentController
	Context            *gin.Context
	Env                *bootstrap.Env
}

func (suite *BlogCommentControllerSuite) SetupTest() {
	suite.BlogCommentUsecase = domain_mocks.NewMockBlogCommentUsecase(suite.T())
	suite.Env = &bootstrap.Env{}

	suite.Controller = &BlogCommentController{
		BlogCommentUsecase: suite.BlogCommentUsecase,
		Env:                suite.Env,
	}

	suite.Context, _ = gin.CreateTestContext(httptest.NewRecorder())

}

func TestBlogCommentControllerSuite(t *testing.T) {
	suite.Run(t, new(BlogCommentControllerSuite))
}

func (s *BlogCommentControllerSuite) TestCreateComment_InvalidJSONRequest() {
	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/123/comments", strings.NewReader(`{"blog_id":123, "comment":}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "123"}}
	c.Set("user_id", "user123")

	s.Context = c

	// Act
	s.Controller.CreateComment(c)

	// Assert
	s.Contains(w.Body.String(), "Invalid request")
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BlogCommentControllerSuite) TestCreateComment_UsecaseReturnsError() {
	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"blog_id":"1", "comment":"Nice post!"}`
	c.Request = httptest.NewRequest("POST", "/api/123/comments", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user123")

	s.Context = c
	s.BlogCommentUsecase.
		On("CreateComment", s.Context, mock.MatchedBy(func(c *domain.BlogComment) bool {
			return c.BlogID == "1" &&
				c.Comment == "Nice post!" &&
				c.AuthorID == "user123"
		})).
		Return(nil, &domain.DomainError{
			Code: 500,
			Err:  errors.New("database error"),
		})

	// Act
	s.Controller.CreateComment(c)

	// Assert
	s.Contains(w.Body.String(), "Error creating comment")
	s.Equal(http.StatusInternalServerError, w.Code)
	s.BlogCommentUsecase.AssertExpectations(s.T())
}

// func (s *BlogCommentControllerSuite) TestCreateComment_Success() {
// 	// Arrange
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	body := `{"blog_id":"1", "comment":"Nice post!"}`
// 	c.Request = httptest.NewRequest("POST", "/api/123/comments", strings.NewReader(body))
// 	c.Request.Header.Set("Content-Type", "application/json")
// 	c.Params = gin.Params{{Key: "id", Value: "123"}}
// 	c.Set("user_id", "user123")

// 	s.Context = c
// 	s.BlogCommentUsecase.
// 		On("CreateComment", mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(comment *domain.BlogComment) bool {
// 			return comment.BlogID == "1" &&
// 				comment.Comment == "Nice post!" &&
// 				comment.AuthorID == "user123"
// 		})).
// 		Return(nil)

// 	// Act
// 	s.Controller.CreateComment(c)

// 	// Assert
// 	s.Contains(w.Body.String(), "Comment created successfully")
// 	s.Equal(http.StatusCreated, w.Code)
// 	s.BlogCommentUsecase.AssertExpectations(s.T())
// }
func (s *BlogCommentControllerSuite) TestDeleteComment_InvalidIDFormat() {
	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/api/comments/", nil)
	c.Request = req

	// Act
	s.Controller.DeleteComment(c)

	// Assert
	s.Contains(w.Body.String(), "Invalid comment ID format")
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *BlogCommentControllerSuite) TestDeleteComment_UsecaseReturnsError() {
	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/api/comments/comment123", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "comment123"}}

	s.Context = c
	s.BlogCommentUsecase.
		On("DeleteComment", s.Context, "comment123").
		Return(&domain.DomainError{
			Code: 500,
			Err:  errors.New("database error"),
		})

	// Act
	s.Controller.DeleteComment(c)

	// Assert
	s.Contains(w.Body.String(), "Error deleting comment")
	s.Equal(http.StatusInternalServerError, w.Code)
	s.BlogCommentUsecase.AssertExpectations(s.T())
}

func (s *BlogCommentControllerSuite) TestDeleteComment_Success() {
	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/api/comments/comment123", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "comment123"}}

	s.Context = c
	s.BlogCommentUsecase.
		On("DeleteComment", s.Context, "comment123").
		Return(nil)

	// Act
	s.Controller.DeleteComment(c)

	// Assert
	s.Contains(w.Body.String(), "Comment deleted successfully")
	s.Equal(http.StatusOK, w.Code)
	s.BlogCommentUsecase.AssertExpectations(s.T())
}
