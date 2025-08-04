package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/ai"
	"g6/blog-api/Infrastructure/database/mongo"
	usecases "g6/blog-api/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

func NewBlogAIRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	blog_ai_controller := controllers.BlogAIController{
		BlogAIUsecase: usecases.NewBlogAIUsecase(
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
		ai.POST("/generate", blog_ai_controller.GenerateBlogContent) // Generate blog content from keywords
	}
}
