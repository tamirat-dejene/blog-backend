package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	dto "g6/blog-api/Delivery/dto"
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

type BlogPostControllerSuite struct {
	suite.Suite
	BlogPostUsecase *domain_mocks.MockBlogPostUsecase
	Controller      *BlogPostController
	Context         *gin.Context
	Env             *bootstrap.Env
}

func (suite *BlogPostControllerSuite) SetupTest() {
	suite.BlogPostUsecase = domain_mocks.NewMockBlogPostUsecase(suite.T())
	suite.Env = &bootstrap.Env{}

	suite.Controller = &BlogPostController{
		BlogPostUsecase: suite.BlogPostUsecase,
		Env:             suite.Env,
	}

	suite.Context, _ = gin.CreateTestContext(httptest.NewRecorder())

}

func TestBlogControllerSuite(t *testing.T) {
	suite.Run(t, new(BlogPostControllerSuite))
}

func (s *BlogPostControllerSuite) TestGetBlogPosts() {
	// blogFilter := domain.BlogPostFilter{} // whatever your filter struct is

	// s.Controller.ParseBlogPostFilter = func(ctx *gin.Context) *domain.BlogPostFilter {
	// 	return &blogFilter
	// }

	// s.Run("Usecase returns error", func() {
	// 	expectedErr := &domain.DomainError{
	// 		Code: http.StatusInternalServerError,
	// 		Err:  fmt.Errorf("db query failed"),
	// 	}

	// 	res := httptest.NewRecorder()
	// 	ctx, _ := gin.CreateTestContext(res)
	// 	ctx.Request = httptest.NewRequest("GET", "/api/blogs", nil)

	// 	s.BlogPostUsecase.EXPECT().
	// 		GetBlogs(ctx, blogFilter).
	// 		Return(nil, expectedErr)

	// 	s.Controller.GetBlogPosts(ctx)

	// 	var got domain.ErrorResponse
	// 	s.NoError(json.Unmarshal(res.Body.Bytes(), &got))
	// 	s.Equal(http.StatusInternalServerError, res.Code)
	// 	s.Equal("Failed to retrieve blogs", got.Message)
	// 	s.Equal("db query failed", got.Error)
	// 	s.Equal(http.StatusInternalServerError, got.Code)
	// })

	// s.Run("Successful retrieval", func() {
	// 	mockBlogs := []domain.BlogPost{
	// 		{ID: "blog1", Title: "Blog 1"},
	// 		{ID: "blog2", Title: "Blog 2"},
	// 	}

	// 	res := httptest.NewRecorder()
	// 	ctx, _ := gin.CreateTestContext(res)
	// 	ctx.Request = httptest.NewRequest("GET", "/api/blogs", nil)

	// 	s.BlogPostUsecase.EXPECT().
	// 		GetBlogs(ctx, blogFilter).
	// 		Return(&mockBlogs, nil)

	// 	s.Controller.GetBlogPosts(ctx)

	// 	var got domain.SuccessResponse
	// 	s.NoError(json.Unmarshal(res.Body.Bytes(), &got))
	// 	s.Equal(http.StatusOK, res.Code)
	// 	s.Equal("Successfully retrieved blogs", got.Message)

	// 	// Check response data
	// 	data := got.Data.(map[string]interface{})
	// 	s.Equal(float64(len(mockBlogs)), data["total_pages"])
	// 	s.Len(data["pages"], len(mockBlogs))
	// })
}

func (s *BlogPostControllerSuite) TestGetBlogByID() {
	userId := "455"
	blogID := "123"

	s.Run("Missing blog ID", func() {
		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("GET", "/api/blogs", nil)
		s.Controller.GetBlogPostByID(ctx)

		s.Equal(http.StatusBadRequest, res.Code)
		s.Contains(res.Body.String(), "Blog ID is required")
	})

	s.Run("Usecase returns error", func() {
		expectedErr := &domain.DomainError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("something went wrong"),
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("GET", "/api/blogs", nil)
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		ctx.Set("user_id", userId)

		s.BlogPostUsecase.EXPECT().
			GetBlogByID(ctx, userId, blogID).
			Return(nil, expectedErr)

		s.Controller.GetBlogPostByID(ctx)

		s.Equal(http.StatusInternalServerError, res.Code)
		s.Contains(res.Body.String(), "Failed to retrieve blog")
		s.Contains(res.Body.String(), "something went wrong")
	})

	s.Run("Successful response", func() {
		expectedBlog := &domain.BlogPost{
			ID:    blogID,
			Title: "My Blog",
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("GET", "/api/blogs/", nil)
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		ctx.Set("user_id", userId)

		s.BlogPostUsecase.EXPECT().
			GetBlogByID(ctx, userId, blogID).
			Return(expectedBlog, nil)

		s.Controller.GetBlogPostByID(ctx)

		s.Equal(http.StatusOK, res.Code)
		s.Contains(res.Body.String(), "Successfully retrieved blog")
		s.Contains(res.Body.String(), "My Blog")
	})

	s.Run("Missing user_id in context but valid blog_id", func() {
		expectedBlog := &domain.BlogPost{
			ID:    blogID,
			Title: "Blog without user id",
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("GET", "/api/blogs/", nil)
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		// user_id not set

		s.BlogPostUsecase.EXPECT().
			GetBlogByID(ctx, "", blogID).
			Return(expectedBlog, nil)

		s.Controller.GetBlogPostByID(ctx)

		s.Equal(http.StatusOK, res.Code)
		s.Contains(res.Body.String(), "Blog without user id")
	})
}

func (s *BlogPostControllerSuite) TestCreateBlog() {
	s.Run("Invalid request payload", func() {
		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		// Create a request with invalid JSON body (empty string)
		ctx.Request = httptest.NewRequest("POST", "/api/blogs", nil)

		s.Controller.CreateBlog(ctx)

		s.Equal(http.StatusBadRequest, res.Code)
		s.Contains(res.Body.String(), "Invalid request payload")
	})

	s.Run("Usecase returns error", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "Test Blog",
			Content: "Some content",
			Tags:    []string{"test", "tags"},
		}
		blogDomain := blogReq.ToDomain()
		userId := "user-123"

		expectedErr := &domain.DomainError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("db error"),
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		// Prepare valid JSON request body
		body := strings.NewReader(`{"title":"Test Blog","content":"Some content","tags":["test","tags"]}`)
		ctx.Request = httptest.NewRequest("POST", "/api/blogs", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Set("user_id", userId)

		s.BlogPostUsecase.EXPECT().
			CreateBlog(ctx, mock.MatchedBy(func(blog *domain.BlogPost) bool {
				return blog.Title == blogDomain.Title && blog.Content == blogDomain.Content && blog.AuthorID == userId
			})).
			Return(nil, expectedErr)

		s.Controller.CreateBlog(ctx)

		s.Equal(http.StatusInternalServerError, res.Code)
		s.Contains(res.Body.String(), "Failed to create blog")
		s.Contains(res.Body.String(), "db error")
	})

	s.Run("Successful creation", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "Success Blog",
			Content: "Some content here",
		}
		userId := "user-321"

		createdBlog := blogReq.ToDomain()
		createdBlog.ID = "blog-999"
		createdBlog.AuthorID = userId

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		body := strings.NewReader(`{"title":"Success Blog","content":"Some content here"}`)
		ctx.Request = httptest.NewRequest("POST", "/api/blogs", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Set("user_id", userId)

		s.BlogPostUsecase.EXPECT().
			CreateBlog(ctx, mock.MatchedBy(func(blog *domain.BlogPost) bool {
				return blog.Title == createdBlog.Title && blog.Content == createdBlog.Content && blog.AuthorID == userId
			})).
			Return(createdBlog, nil)

		s.Controller.CreateBlog(ctx)

		s.Equal(http.StatusCreated, res.Code)
		s.Contains(res.Body.String(), "Successfully created blog")
		s.Contains(res.Body.String(), createdBlog.Title)
	})

	s.Run("Missing user_id in context", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "No UserID Blog",
			Content: "Content without user id",
		}
		createdBlog := blogReq.ToDomain()
		createdBlog.ID = "blog-100"
		createdBlog.AuthorID = "" // since user_id missing in context

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		body := strings.NewReader(`{"title":"No UserID Blog","content":"Content without user id"}`)
		ctx.Request = httptest.NewRequest("POST", "/api/blogs", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		// No user_id set in context here

		s.BlogPostUsecase.EXPECT().
			CreateBlog(ctx, mock.MatchedBy(func(blog *domain.BlogPost) bool {
				return blog.Title == createdBlog.Title && blog.Content == createdBlog.Content && blog.AuthorID == ""
			})).
			Return(createdBlog, nil)

		s.Controller.CreateBlog(ctx)

		s.Equal(http.StatusCreated, res.Code)
		s.Contains(res.Body.String(), "Successfully created blog")
		s.Contains(res.Body.String(), createdBlog.Title)
	})
}

func (s *BlogPostControllerSuite) TestUpdateBlog() {
	blogID := "blog-123"
	userID := "user-456"

	s.Run("Invalid request payload", func() {
		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		// Send invalid JSON
		body := strings.NewReader(`{"title": 123}`) // title should be string
		ctx.Request = httptest.NewRequest("PUT", "/api/blogs/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}

		s.Controller.UpdateBlog(ctx)

		s.Equal(http.StatusBadRequest, res.Code)
		s.Contains(res.Body.String(), "Invalid request payload")
	})

	s.Run("Usecase returns error", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "Update Fail",
			Content: "Broken content",
		}

		expectedErr := &domain.DomainError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("db error"),
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		body := strings.NewReader(`{"title":"Update Fail","content":"Broken content"}`)
		ctx.Request = httptest.NewRequest("PUT", "/api/blogs/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		ctx.Set("user_id", userID)

		// Loose matcher: only check fields we care about
		s.BlogPostUsecase.EXPECT().
			UpdateBlog(ctx, blogID, mock.MatchedBy(func(b domain.BlogPost) bool {
				return b.Title == blogReq.Title &&
					b.Content == blogReq.Content &&
					b.AuthorID == userID &&
					b.ID == blogID
			})).
			Return(nil, expectedErr)

		s.Controller.UpdateBlog(ctx)

		s.Equal(http.StatusInternalServerError, res.Code)
		s.Contains(res.Body.String(), "Failed to update blog")
		s.Contains(res.Body.String(), "db error")
	})

	s.Run("Successful update", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "Updated Blog",
			Content: "Updated content here",
		}
		updatedDomain := blogReq.ToDomain()
		updatedDomain.ID = blogID
		updatedDomain.AuthorID = userID

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		body := strings.NewReader(`{"title":"Updated Blog","content":"Updated content here"}`)
		ctx.Request = httptest.NewRequest("PUT", "/api/blogs/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		ctx.Set("user_id", userID)

		s.BlogPostUsecase.EXPECT().
			UpdateBlog(ctx, blogID, mock.MatchedBy(func(b domain.BlogPost) bool {
				return b.Title == updatedDomain.Title &&
					b.Content == updatedDomain.Content &&
					b.AuthorID == userID &&
					b.ID == blogID
			})).
			Return(updatedDomain, nil)

		s.Controller.UpdateBlog(ctx)

		s.Equal(http.StatusOK, res.Code)
		s.Contains(res.Body.String(), "Successfully updated blog")
		s.Contains(res.Body.String(), updatedDomain.Title)
	})

	s.Run("Missing user_id in context", func() {
		blogReq := dto.BlogPostRequest{
			Title:   "No User",
			Content: "Updated without user id",
		}
		updatedDomain := blogReq.ToDomain()
		updatedDomain.ID = blogID
		updatedDomain.AuthorID = ""

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)

		body := strings.NewReader(`{"title":"No User","content":"Updated without user id"}`)
		ctx.Request = httptest.NewRequest("PUT", "/api/blogs/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}
		// No user_id set

		s.BlogPostUsecase.EXPECT().
			UpdateBlog(ctx, blogID, mock.MatchedBy(func(b domain.BlogPost) bool {
				return b.Title == updatedDomain.Title &&
					b.Content == updatedDomain.Content &&
					b.AuthorID == "" &&
					b.ID == blogID
			})).
			Return(updatedDomain, nil)

		s.Controller.UpdateBlog(ctx)

		s.Equal(http.StatusOK, res.Code)
		s.Contains(res.Body.String(), "Successfully updated blog")
		s.Contains(res.Body.String(), updatedDomain.Title)
	})
}

func (s *BlogPostControllerSuite) TestDeleteBlog() {
	blogID := "blog-123"

	s.Run("Usecase returns error", func() {
		expectedErr := &domain.DomainError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("db delete error"),
		}

		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("DELETE", "/api/blogs/"+blogID, nil)
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}

		s.BlogPostUsecase.EXPECT().
			DeleteBlog(ctx, blogID).
			Return(expectedErr)

		s.Controller.DeleteBlog(ctx)

		s.Equal(http.StatusInternalServerError, res.Code)
		s.Contains(res.Body.String(), "Failed to delete blog")
		s.Contains(res.Body.String(), "db delete error")
	})

	s.Run("Successful deletion", func() {
		res := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(res)
		ctx.Request = httptest.NewRequest("DELETE", "/api/blogs/"+blogID, nil)
		ctx.Params = gin.Params{{Key: "id", Value: blogID}}

		s.BlogPostUsecase.EXPECT().
			DeleteBlog(ctx, blogID).
			Return(nil)

		s.Controller.DeleteBlog(ctx)

		s.Equal(http.StatusOK, res.Code)
		s.Contains(res.Body.String(), "Successfully deleted blog")
	})
}
