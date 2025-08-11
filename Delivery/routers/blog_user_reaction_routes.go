package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/middleware"
	repository "g6/blog-api/Repositories/blog"
	usecases "g6/blog-api/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

func NewBlogUserReactionRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	blogUserReactionGroup := api.Group("/blog/reactions", middleware.AuthMiddleware(*env))

	// Initialize the blog user reaction repository, usecase, and controller
	blog_user_reaction_controller := controllers.BlogReactionController{
		BlogUserReactionUsecase: usecases.NewBlogUserReactionUsecase(
			repository.NewUserReactionRepo(db, &mongo.Collections{
			BlogPosts:         env.BlogPostCollection,
			BlogComments:      env.BlogCommentCollection,
			BlogUserReactions: env.BlogUserReactionCollection,
		}), time.Duration(env.CtxTSeconds)*time.Second),
		Env: env,
	}

	// Define routes for blog user reactions
	blogUserReactionGroup.POST("/", blog_user_reaction_controller.CreateReaction)      // Create a new reaction
	blogUserReactionGroup.GET("/", blog_user_reaction_controller.GetUserReaction)      // Get user reaction
	blogUserReactionGroup.DELETE("/:id", blog_user_reaction_controller.DeleteReaction) // Delete user reaction
}
