package routers

import (
	config "g6/blog-api/Configs"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	userUsecase "g6/blog-api/Usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRoutes(env *config.Env, group *gin.RouterGroup, db *mongo.Database) {
	ur := repositories.NewUserRepository(db, "users")
	uu := userUsecase.NewUserUsecase(ur)
	uc := controllers.NewUserController(uu)

	group.GET("/users/:id", uc.GetAllUsers)
	// group.POST("/users", uc.CreateUser)

}
