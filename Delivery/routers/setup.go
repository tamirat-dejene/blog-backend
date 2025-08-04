package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Infrastructure/database/mongo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusPermanentRedirect, "/api") })

	api := router.Group("/api")
	{
		NewAuthRoutes(env, api, db)
		NewUserRoutes(env, api, db)
		NewBlogRoutes(env, api, db)
		NewBlogCommentRoutes(env, api, db)
		NewBlogUserReactionRoutes(env, api, db)
		NewBlogAIRoutes(env, api, db)
	}
}
