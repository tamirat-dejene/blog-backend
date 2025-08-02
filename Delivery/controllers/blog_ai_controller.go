package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"github.com/gin-gonic/gin"
)

type BlogAIController struct {
	BlogAIUsecase domain.BlogAIUsecase
	Env           *bootstrap.Env
}

func (b *BlogAIController) GenerateBlogContent(ctx *gin.Context) {
	var req dto.BlogAIGenerateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid request format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	user_id, ok := ctx.Get("user_id")
	if !ok || user_id == nil {
		ctx.JSON(401, domain.ErrorResponse{
			Message: "Unauthorized",
			Error:   "User ID not found",
			Code:    401,
		})
		return
	}

	content, err := b.BlogAIUsecase.GenerateContent(ctx, domain.BlogAIGenerate{
		UserID:   user_id.(string),
		Topic:    req.Topic,
		Keywords: req.Keywords,
	})

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Failed to generate content",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	resp, err1 := mapper.BlogAIContentFromDomain(content)
	if err1 != nil {
		ctx.JSON(500, domain.ErrorResponse{
			Message: "Failed to map generated content.",
			Error:   err1.Error(),
			Code:    500,
		})
		return
	}

	ctx.JSON(200, resp)
}

func (b *BlogAIController) GetGeneratedContent(ctx *gin.Context) {
	var uriParams struct {
		ID string `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid content ID format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	content, err := b.BlogAIUsecase.GetGeneratedContentByID(ctx, uriParams.ID)
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Error fetching generated content",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	resp, err1 := mapper.BlogAIContentFromDomain(content)
	if err1 != nil {
		ctx.JSON(500, domain.ErrorResponse{
			Message: "Failed to map generated content.",
			Error:   err1.Error(),
			Code:    500,
		})
		return
	}

	ctx.JSON(200, resp)
}

func (b *BlogAIController) GetPromptHistory(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok || userID == nil {
		ctx.JSON(401, domain.ErrorResponse{
			Message: "Unauthorized",
			Error:   "User ID not found",
			Code:    401,
		})
		return
	}

	prompts, err := b.BlogAIUsecase.GetUserPromptHistory(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Error fetching prompt history",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	ctx.JSON(200, prompts)
}

func (b *BlogAIController) SubmitFeedback(ctx *gin.Context) {
	var req dto.BlogAIFeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Message: "Invalid request format",
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	user_id, ok := ctx.Get("user_id")
	if !ok || user_id == nil {
		ctx.JSON(401, domain.ErrorResponse{
			Message: "Unauthorized",
			Error:   "User ID not found",
			Code:    401,
		})
		return
	}

	err := b.BlogAIUsecase.SubmitFeedback(ctx, domain.BlogAIFeedback{
		ContentID: req.ContentID,
		UserID:    user_id.(string),
		Rating:    req.Rating,
		Feedback:  req.Feedback,
	})

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Message: "Error submitting feedback",
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Feedback submitted successfully",
	})
}
