package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	dto "g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"

	"github.com/gin-gonic/gin"
)

type BlogAIController struct {
	BlogAIUsecase domain.AIBlogPostUsecase
	Env           *bootstrap.Env
}

func (b *BlogAIController) GenerateBlogContent(ctx *gin.Context) {
	var req dto.AIBlogPostGenerateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, domain.ErrorResponse{
			Error:   err.Error(),
			Code:    400,
		})
		return
	}

	generated_blog, err := b.BlogAIUsecase.GeneratePost(ctx, domain.AIBlogPostGenerate{
		Topic:    req.Topic,
		Keywords: req.Keywords,
	})

	if err != nil {
		ctx.JSON(err.Code, domain.ErrorResponse{
			Error:   err.Err.Error(),
			Code:    err.Code,
		})
		return
	}
	
	response := dto.AIBlogPostResponse{}
	response.FromDomain(generated_blog)
	
	ctx.JSON(200, response)
}
