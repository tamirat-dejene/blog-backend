package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogCommentController struct {
	BlogCommentUsecase domain.BlogCommentUsecase
	Env                *bootstrap.Env
}

func (b *BlogCommentController) CreateComment(ctx *gin.Context) {
	var req dto.BlogCommentRequest

	err := ctx.ShouldBind(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment := dto.ToDomainBlogComment(req)

	createdComment, err := b.BlogCommentUsecase.CreateComment(ctx, comment)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": createdComment})
}

func (b *BlogCommentController) DeleteComment(ctx *gin.Context) {
	id := ctx.Param("id")

	err := b.BlogCommentUsecase.DeleteComment(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func (b *BlogCommentController) GetCommentByID(ctx *gin.Context) {
	var uriParams struct {
		id string `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid comment ID format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	comment, domain_err := b.BlogCommentUsecase.GetCommentByID(ctx, uriParams.id)
	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Message: domain_err.Err.Error(),
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	ctx.JSON(200, comment)
}

func (b *BlogCommentController) GetCommentsByBlogID(ctx *gin.Context) {
	var limit = ctx.DefaultQuery("limit", "30") // Default limit to 30 if not provided
	// Convert limit to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid limit format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	comments, domain_err := b.BlogCommentUsecase.GetCommentsByBlogID(ctx, ctx.Param("id"), limitInt)

	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Message: "Error fetching comments",
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	ctx.JSON(200, comments)
}

func (b *BlogCommentController) UpdateComment(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.BlogCommentRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	comment := dto.ToDomainBlogComment(req)
	updatedComment, err := b.BlogCommentUsecase.UpdateComment(ctx, id, comment)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Updated Comment": updatedComment})
}
