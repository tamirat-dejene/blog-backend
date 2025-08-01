package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"
	"strconv"

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

	return &domain.BlogPostFilter{
		Page:       pageInt,
		PageSize:   pageSizeInt,
		Recency:    domain.Recency(recency),
		Tags:       ctx.QueryArray("tags"), // in the url this looks like ?tags=tag1,tag2. if no tags are provided, it will be an empty slice
		AuthorName: ctx.Query("authorName"),
		Title:      ctx.Query("title"),
		Popular:    most_popular == "true", // convert string to bool
	}
}

func (b *BlogPostController) GetBlogPosts(ctx *gin.Context) {
	filter := b.parseBlogPostFilter(ctx)
	paginated_blogs, err := b.BlogPostUsecase.GetBlogs(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to retrieve blogs",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully retrieved blogs",
		Data:    gin.H{"TotalPages": len(paginated_blogs), "Pages": paginated_blogs},
	})
}

func (b *BlogPostController) GetBlogPostByID(ctx *gin.Context) {
	blog_id := ctx.Param("id")
	if blog_id == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Blog ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	blog, err := b.BlogPostUsecase.GetBlogByID(ctx, blog_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to retrieve blog",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Successfully retrieved blog",
		Data:    blog,
	})
}

func (b *BlogPostController) CreateBlog(ctx *gin.Context) {
	var req dto.BlogPostRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	blog := dto.ToDomainBlogPost(req)
	createdBlog, err := b.BlogPostUsecase.CreateBlog(ctx, &blog)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdBlog)
}

func (b *BlogPostController) UpdateBlog(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.BlogPostRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	blog := dto.ToDomainBlogPost(req)

	updatedBlog, err := b.BlogPostUsecase.UpdateBlog(ctx, id, blog)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedBlog)

}

func (b *BlogPostController) DeleteBlog(ctx *gin.Context) {
	id := ctx.Param("id")

	err := b.BlogPostUsecase.DeleteBlog(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
