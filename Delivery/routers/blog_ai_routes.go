package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/ai"
	"g6/blog-api/Infrastructure/database/mongo"
	repository "g6/blog-api/Repositories/blog"
	usecases "g6/blog-api/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

func NewBlogAIRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	blog_ai_controller := controllers.BlogAIController{
		BlogAIUsecase: usecases.NewBlogAIUsecase(
			repository.NewBlogAIRepository(db, &repository.Collections{
				BlogPosts:         env.BlogPostCollection,
				BlogComments:      env.BlogCommentCollection,
				BlogUserReactions: env.BlogUserReactionCollection,
			}),
			ai.GeminiConfig{
				APIKey:    env.GeminiAPIKey,
				ModelName: env.GeminiModelName,
			},
			time.Duration(env.CtxTSeconds)*time.Second,
		),
		Env: env,
	}

	ai := api.Group("/ai/blog")
	{
		ai.POST("/generate", blog_ai_controller.GenerateBlogContent)     // Generate blog content from keywords
		ai.GET("/generate/:id", blog_ai_controller.GetGeneratedContent)  // Fetch a specific AI-generated result
		ai.GET("/keywords/history", blog_ai_controller.GetPromptHistory) // List previous user prompts
		ai.POST("/feedback", blog_ai_controller.SubmitFeedback)          // Submit feedback on generated content
	}
}
