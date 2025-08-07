package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/middleware"
	"g6/blog-api/Infrastructure/redis"
	repository "g6/blog-api/Repositories/blog"
	usecases "g6/blog-api/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

func NewBlogCommentRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	comment_controller := controllers.BlogCommentController{
		BlogCommentUsecase: usecases.NewBlogCommentUsecase(
			repository.NewBlogCommentRepository(
				db,
				&mongo.Collections{
					BlogPosts:         env.BlogPostCollection,
					BlogComments:      env.BlogCommentCollection,
					BlogUserReactions: env.BlogUserReactionCollection,
				},
			),
			redis.NewRedisClient(env, &redis.RedisService{}),
			time.Duration(env.CtxTSeconds)*time.Second,
		),
		Env: env,
	}

	// Routes for managing comments on a specific blog
	blog_comments := api.Group("/blogs/:id/comments")
	{
		blog_comments.POST("/", middleware.AuthMiddleware(*env), comment_controller.CreateComment) // Create a comment for a blog
		blog_comments.GET("/", comment_controller.GetCommentsByBlogID)                             // Get all comments for a blog
	}

	// General comment routes (independent of blog)
	comments := api.Group("/comments")
	{
		comments.GET("/:id", comment_controller.GetCommentByID)   // Get comment by ID
		comments.PUT("/:id", comment_controller.UpdateComment)    // Update a comment by ID
		comments.DELETE("/:id", comment_controller.DeleteComment) // Delete a comment by ID
	}
}
