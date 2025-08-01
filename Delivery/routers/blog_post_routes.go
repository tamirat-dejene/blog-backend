package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	"g6/blog-api/Infrastructure/database/mongo"
	repository "g6/blog-api/Repositories/blog"
	usecases "g6/blog-api/Usecases"
	"time"

	"github.com/gin-gonic/gin"
)

func NewBlogRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	blogGroup := api.Group("/blogs")

	blog_post_controller := controllers.BlogPostController{
		BlogPostUsecase: usecases.NewBlogPostUsecase(repository.NewBlogPostRepo(db, repository.NewCollections(env.BlogPostCollection, env.BlogCommentCollection, env.BlogUserReactionCollection)), time.Duration(env.CtxTSeconds)*time.Second),
		Env:             env,
	}

	blogGroup.GET("/", blog_post_controller.GetBlogPosts) // Get all blogs with optional filters
	// blogGroup.GET("/:id", blog_controller.GetBlogByID) // Get a single blog by ID
	// blogGroup.POST("/", blog_controller.CreateBlog) // Create a new blog
	// blogGroup.PUT("/:id", blog_controller.UpdateBlog) // Update an existing blog
	// blogGroup.DELETE("/:id", blog_controller.DeleteBlog) // Delete a blog by ID
}
