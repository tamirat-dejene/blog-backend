package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogCommentController struct {
	BlogCommentUsecase domain.BlogCommentUsecase
	Env         *bootstrap.Env
}


func (b *BlogCommentController) CreateComment(ctx *gin.Context) {
	// Implementation for creating a comment
}

func (b *BlogCommentController) DeleteComment(ctx *gin.Context) {
	// Implementation for deleting a comment
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
	// Implementation for getting comments by blog ID
	var uriParams dto.BlogIDParams
	var limit = ctx.DefaultQuery("limit", "30") // Default limit to 30 if not provided

	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid blog ID format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
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

	comments, domain_err := b.BlogCommentUsecase.GetCommentsByBlogID(ctx, uriParams.BlogID, limitInt)

	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Message: domain_err.Err.Error(),
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	ctx.JSON(200, comments)
}

func (b *BlogCommentController) UpdateComment(ctx *gin.Context) {
	// Implementation for updating a comment
}
