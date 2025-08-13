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

type BlogCommentController struct {
	BlogCommentUsecase domain.BlogCommentUsecase
	Env                *bootstrap.Env
}

func (b *BlogCommentController) CreateComment(ctx *gin.Context) {
	var req dto.BlogCommentRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Convert dto to domain model, and set the AuthorID if needed
	var comment = req.ToDomain()
	comment.AuthorID = ctx.GetString("user_id") // Assuming user_id is set in context after authentication
	comment.BlogID = req.BlogID                 // Set the BlogID from the request

	// Create the comment using the usecase
	createdComment, domain_err := b.BlogCommentUsecase.CreateComment(ctx, comment)
	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	// Convert the created comment to response DTO
	var response dto.BlogCommentResponse
	response.Parse(createdComment)

	// Return the created comment response
	ctx.JSON(http.StatusCreated, domain.SuccessResponse{
		Message: "Comment created successfully",
		Data:    response,
	})
}

func (b *BlogCommentController) DeleteComment(ctx *gin.Context) {
	// Extract the comment ID from the URL parameters
	if err := ctx.ShouldBindUri(&struct {
		ID string `uri:"id" binding:"required"`
	}{}); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	// Call the usecase to delete the comment
	id := ctx.Param("id")
	err := b.BlogCommentUsecase.DeleteComment(ctx, id)

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Comment deleted successfully",
		Data:    nil,
	})
}

func (b *BlogCommentController) GetCommentByID(ctx *gin.Context) {
	var uriParams struct {
		ID string `uri:"id" binding:"required"`
	}

	// Extract the comment ID from the URL parameters
	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    400,
		})
		return
	}
	fmt.Println("Fetching comment by ID:", uriParams.ID)

	// Call the usecase to get the comment by ID
	comment, domain_err := b.BlogCommentUsecase.GetCommentByID(ctx, uriParams.ID)
	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	// Convert the comment to response DTO
	var commentResponse dto.BlogCommentResponse
	commentResponse.Parse(comment)

	// Return the comment response
	ctx.JSON(200, domain.SuccessResponse{
		Message: "Comment fetched successfully",
		Data:    commentResponse,
	})
}

func (b *BlogCommentController) GetCommentsByBlogID(ctx *gin.Context) {
	var limit = ctx.DefaultQuery("limit", "30") // Default limit to 30 if not provided
	// Convert limit to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		ctx.JSON(400, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	// Call the usecase to get comments by blog ID
	comments, domain_err := b.BlogCommentUsecase.GetCommentsByBlogID(ctx, ctx.Param("id"), limitInt)
	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	// Convert the comments to response DTOs
	var commentResponses []dto.BlogCommentResponse
	for _, comment := range comments {
		var commentResponse dto.BlogCommentResponse
		commentResponse.Parse(&comment)
		commentResponses = append(commentResponses, commentResponse)
	}

	// Return the comments response
	ctx.JSON(200, domain.SuccessResponse{
		Message: "Comments fetched successfully",
		Data:    commentResponses,
	})
}

func (b *BlogCommentController) UpdateComment(ctx *gin.Context) {
	// Extract the comment ID from the URL parameters
	if err := ctx.ShouldBindUri(&struct {
		ID string `uri:"id"`
	}{}); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	id := ctx.Param("id")

	// Bind the request body to the BlogCommentRequest DTO
	var req dto.BlogCommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Convert the request DTO to domain model and call the usecase to update the comment
	comment := req.ToDomain()
	comment.ID = id                             // Set the ID for the comment to update
	comment.AuthorID = ctx.GetString("user_id") // Assuming user_id is set in context

	updatedComment, domain_err := b.BlogCommentUsecase.UpdateComment(ctx, id, comment)

	if domain_err != nil {
		ctx.JSON(domain_err.Code, domain.ErrorResponse{
			Error:   domain_err.Err.Error(),
			Code:    domain_err.Code,
		})
		return
	}

	// Convert the updated comment to response DTO
	var response dto.BlogCommentResponse
	response.Parse(updatedComment)

	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Comment updated successfully",
		Data:    response,
	})
}
