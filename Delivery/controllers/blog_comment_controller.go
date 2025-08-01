package controllers

import (
	"g6/blog-api/Delivery/bootstrap"
	domain "g6/blog-api/Domain"

	"github.com/gin-gonic/gin"
)

type BlogCommentController struct {
	CommentRepo domain.BlogCommentRepository
	Env         *bootstrap.Env
}

func NewBlogCommentController(commentRepo domain.BlogCommentRepository, env *bootstrap.Env) *BlogCommentController {
	return &BlogCommentController{
		CommentRepo: commentRepo,
		Env:         env,
	}
}

func (b *BlogCommentController) CreateComment(ctx *gin.Context) {
	// Implementation for creating a comment
}

func (b *BlogCommentController) DeleteComment(ctx *gin.Context) {
	// Implementation for deleting a comment
}

func (b *BlogCommentController) GetCommentByID(ctx *gin.Context) {
	// Implementation for getting a comment by ID
}

func (b *BlogCommentController) GetCommentsByBlogID(ctx *gin.Context) {
	// Implementation for getting comments by blog ID
}

func (b *BlogCommentController) UpdateComment(ctx *gin.Context) {
	// Implementation for updating a comment
}