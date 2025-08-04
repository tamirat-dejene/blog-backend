package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BlogPostController struct {
	BlogPostUsecase domain.BlogPostUsecase
	Env             *bootstrap.Env
}

func (b *BlogPostController) parseBlogPostFilter(ctx *gin.Context) *domain.BlogPostFilter {
	page := ctx.DefaultQuery("page", fmt.Sprint(b.Env.Page))
	page_size := ctx.DefaultQuery("pageSize", fmt.Sprint(b.Env.PageSize))
	recency := ctx.DefaultQuery("recency", b.Env.Recency)
	most_popular := ctx.DefaultQuery("mostPopular", "false")

	// check if the page and pageSize are valid numbers
	if _, err := strconv.Atoi(page); err != nil {
		page = fmt.Sprint(b.Env.Page) // if not a number, default to env.Page
	}
	if _, err := strconv.Atoi(page_size); err != nil {
		page_size = fmt.Sprint(b.Env.PageSize) // if not a number, default to env.PageSize
	}

	// check if recency is either "newest" or "oldest"
	if recency != string(domain.RecencyNewest) && recency != string(domain.RecencyOldest) {
		recency = b.Env.Recency
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(page_size)
	tgs := ctx.Query("tags")
	var tags []string
	if tgs != "" && tgs != "null" {
		tags = strings.Split(strings.TrimSpace(tgs), ",")
	}

	return &domain.BlogPostFilter{
		Page:       pageInt,
		PageSize:   pageSizeInt,
		Recency:    domain.Recency(recency),
		Tags:       tags,
		AuthorName: ctx.Query("authorName"),
		Title:      ctx.Query("title"),
		Popular:    most_popular == "true", // convert string to bool
	}
}

func (b *BlogPostController) GetBlogPosts(ctx *gin.Context) {
	filter := b.parseBlogPostFilter(ctx)
	paginated_blogs, err := b.BlogPostUsecase.GetBlogs(ctx, filter)
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to retrieve blogs",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// Convert paginated blogs to response DTO
	var response = make([]dto.BlogPostsPageResponse, len(paginated_blogs))
	for idx, page := range paginated_blogs {
		page_response := dto.BlogPostsPageResponse{}
		page_response.Parse(&page)
		response[idx] = page_response
	}

	// Return the response
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully retrieved blogs",
		Data:    gin.H{"total_pages": len(paginated_blogs), "pages": response},
	})
}

func (b *BlogPostController) GetBlogPostByID(ctx *gin.Context) {
	// Get the blog ID from the URL parameters
	blog_id := ctx.Param("id")
	if blog_id == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Blog ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Attempt to retrieve the blog post by ID
	blog, err := b.BlogPostUsecase.GetBlogByID(ctx, blog_id)
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to retrieve blog",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// Convert the blog post to a response DTO
	var blog_response dto.BlogPostResponse
	blog_response.Parse(blog)

	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully retrieved blog",
		Data:    blog_response,
	})
}

func (b *BlogPostController) CreateBlog(ctx *gin.Context) {
	var req dto.BlogPostRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid request payload",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	blog_post := req.ToDomain()
	blog_post.AuthorID = ctx.GetString("user_id") // user_id is set in the context from middleware

	createdBlog, err := b.BlogPostUsecase.CreateBlog(ctx, blog_post)

	// Send create error response if any
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to create blog",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// Convert the created blog to response DTO
	var response dto.BlogPostResponse
	response.Parse(createdBlog)

	// Return the response
	ctx.JSON(http.StatusCreated, domain.SuccessResponse{
		Message: "Successfully created blog",
		Data:    response,
	})
}

func (b *BlogPostController) UpdateBlog(ctx *gin.Context) {
	id := ctx.Param("id")

	// Bind request body to DTO
	var req dto.BlogPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid request payload",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Convert DTO to domain model
	blog := req.ToDomain()
	blog.ID = id
	blog.AuthorID = ctx.GetString("user_id") // Extracted from auth middleware

	// Attempt to update the blog post
	updatedBlog, err := b.BlogPostUsecase.UpdateBlog(ctx, id, *blog)
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to update blog",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	var response dto.BlogPostResponse
	response.Parse(&updatedBlog)

	// Success response
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully updated blog",
		Data:    response,
	})
}

func (b *BlogPostController) DeleteBlog(ctx *gin.Context) {
	id := ctx.Param("id")

	err := b.BlogPostUsecase.DeleteBlog(ctx, id)

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to delete blog",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully deleted blog",
		Data:    nil,
	})
}
