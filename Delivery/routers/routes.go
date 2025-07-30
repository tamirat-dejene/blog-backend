package routers

import (
	"g6/blog-api/Configs"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes initializes the routes for the application
func SetupRoutes(env *Configs.Env, db *mongo.Database) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		NewUserRoutes(env, api, db)
	}
	return router
}
