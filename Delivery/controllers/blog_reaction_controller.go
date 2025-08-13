package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogReactionController struct {
	BlogUserReactionUsecase domain.BlogUserReactionUsecase
	Env                     *bootstrap.Env
}

func (b *BlogReactionController) CreateReaction(ctx *gin.Context) {
	// Bind the request to the DTO
	var req dto.BlogUserReactionRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	reaction := req.ToDomain()
	reaction.UserID = ctx.GetString("user_id") // user_id is set in the context from middleware

	// Convert the DTO to the domain model
	createdReaction, err := b.BlogUserReactionUsecase.CreateReaction(ctx, reaction)

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// Convert the domain model back to the response DTO
	var response dto.BlogUserReactionResponse
	response.Parse(createdReaction)
	ctx.JSON(http.StatusCreated, response)
}

func (b *BlogReactionController) DeleteReaction(ctx *gin.Context) {
	// Extract the ID from the URL parameters
	if ctx.Param("id") == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "Missing reaction ID in request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Call the use case to delete the reaction
	err := b.BlogUserReactionUsecase.DeleteReaction(ctx, ctx.Param("id"))
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	// If successful, return a success message
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Reaction deleted successfully",
		Data:    nil,
	})
}

func (b *BlogReactionController) GetUserReaction(ctx *gin.Context) {
	// Bind the query parameters to the DTO
	var query dto.ReactionQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	// Call the use case to get the user reaction
	reaction, err := b.BlogUserReactionUsecase.GetUserReaction(ctx, query.BlogId, ctx.GetString("user_id"))

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}
	// Convert the domain model to the response DTO
	var response dto.BlogUserReactionResponse
	response.Parse(reaction)

	// Return the response
	ctx.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "User reaction retrieved successfully",
		Data:    response,
	})
}
