package controllers

import (
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	domain "g6/blog-api/Domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	BlogUsecase domain.BlogPostUsecase
	Env         *bootstrap.Env
}

func (b *BlogController) parseBlogFilter(ctx *gin.Context) *domain.BlogPostFilter {
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

func (b *BlogController) GetBlogs(ctx *gin.Context) {
	filter := b.parseBlogFilter(ctx)
	blogs, err := b.BlogUsecase.GetBlogs(ctx, filter)
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
		Data:    blogs,
	})
}