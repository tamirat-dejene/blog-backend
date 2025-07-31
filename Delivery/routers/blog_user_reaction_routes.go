package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Infrastructure/database/mongo"

	"github.com/gin-gonic/gin"
)

func NewBlogUserReactionRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	// blogUserReactionGroup := api.Group("/blog/reactions")

	// Initialize the blog user reaction repository, usecase, and controller
	// blogUserReactionRepo := repository.NewBlogUserReactionRepo(db, env.BlogUserReactionCollection)
	// blogUserReactionUsecase := usecases.NewBlogUserReactionUsecase(blogUserReactionRepo, time.Duration(env.CtxTSeconds)*time.Second)
	// blogUserReactionController := controllers.BlogUserReactionController{
	// 	BlogUserReactionUsecase: blogUserReactionUsecase,
	// 	Env:                     env,
	// }

	// Define routes for blog user reactions
	// blogUserReactionGroup.POST("/", blogUserReactionController.CreateBlogUserReaction) // Create a new reaction
	// Add more blog user reaction-related routes here as needed
}
