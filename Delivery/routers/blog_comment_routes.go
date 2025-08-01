package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/database/mongo"
	repository "g6/blog-api/Repositories/blog"

	"github.com/gin-gonic/gin"
)

func NewBlogCommentRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	comments := api.Group("/blog/:blogID/comments")

	comment_controller := controllers.BlogCommentController{
		CommentRepo: repository.NewBlogCommentRepository(db, repository.NewCollections(env.BlogCommentCollection, env.BlogPostCollection, env.BlogUserReactionCollection)),
		Env:         env,
	}

	{
		comments.POST("/", comment_controller.CreateComment)
		comments.DELETE("/:id", comment_controller.DeleteComment)
		comments.GET("/:id", comment_controller.GetCommentByID)
		comments.GET("/blog/:blogID", comment_controller.GetCommentsByBlogID)
		comments.PUT("/:id", comment_controller.UpdateComment)
	}
}
