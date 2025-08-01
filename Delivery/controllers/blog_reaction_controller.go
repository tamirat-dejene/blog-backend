package controllers

import (
	"g6/blog-api/Delivery/dto"
	domain "g6/blog-api/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogReactionController struct {
	BlogUserReactionUsecase domain.BlogUserReactionUsecase
}

func (b *BlogReactionController) CreateReaction(ctx *gin.Context) {
	var req dto.BlogUserReactionRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	reaction := dto.ToDomainBlogReaction(req)

	createdReaction, err := b.BlogUserReactionUsecase.CreateReaction(ctx, reaction)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusCreated, createdReaction)
}

func (b *BlogReactionController) DeleteReaction(ctx *gin.Context) {
	id := ctx.Param("id")

	err := b.BlogUserReactionUsecase.DeleteReaction(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

func (b *BlogReactionController) GetUserReaction(ctx *gin.Context) {
	//GET("/blogs/:blog_id/users/:user_id/reaction"
	blogId := ctx.Param("blog_id")
	userId := ctx.Param("user_id")

	reaction, err := b.BlogUserReactionUsecase.GetUserReaction(ctx, blogId, userId)

	if err != nil {
		if err.Error() == "no reaction found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := dto.FromDomainBlogReaction(reaction)

	ctx.JSON(http.StatusOK, response)
}
