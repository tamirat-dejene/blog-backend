package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Infrastructure/database/mongo"

	"github.com/gin-gonic/gin"
)

func NewBlogCommentRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	// blogCommentGroup := api.Group("/blog/comments")

	// // Initialize the blog comment repository, usecase, and controller
	// blogCommentRepo := repository.NewBlogCommentRepo(db, env.BlogCommentCollection)
	// blogCommentUsecase := usecases.NewBlogCommentUsecase(blogCommentRepo, time.Duration(env.CtxTSeconds)*time.Second)
	// blogCommentController := controllers.BlogCommentController{
	// 	BlogCommentUsecase: blogCommentUsecase,
	// 	Env:                env,
	// }

	// // Define routes for blog comments
	// blogCommentGroup.GET("/", blogCommentController.GetBlogComments) // Get all comments for a blog
	// // Add more blog comment-related routes here as needed
}
