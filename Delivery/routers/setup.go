package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/database/mongo"
	repository "g6/blog-api/Repositories/blog"
	usecases "g6/blog-api/Usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, router *gin.Engine) {
	// homepage
	router.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusPermanentRedirect, "/api") })

	// authGroup := router.Group("/api/auth")
	{
		// Define auth-related routes here, e.g.:
		// authGroup.POST("/login", authHandler.Login)
		// authGroup.POST("/register", authHandler.Register)
	}
	
	
	blogGroup := router.Group("/api/blogs") // Blog-related routes

	blog_repo := repository.NewBlogRepo(db, env.BlogCollection)
	blog_usecase := usecases.NewBlogUsecase(blog_repo, timeout)
	blog_controller := controllers.BlogController{
		BlogUsecase: blog_usecase,
		Env:         env,
	}
	blogGroup.GET("/", blog_controller.GetBlogs) // Get all blogs with optional filters
	// ....

	
}
