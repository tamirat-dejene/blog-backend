package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	usecases "g6/blog-api/Usecases"
	"time"

	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/middleware"
	"g6/blog-api/Infrastructure/storage"

	"github.com/gin-gonic/gin"
)

func NewUserRoutes(env *bootstrap.Env, group *gin.RouterGroup, db mongo.Database) {
	// context time out
	ctxTimeout := time.Duration(env.CtxTSeconds) * time.Second

	// storage services
	imageKitStorageService := storage.NewImageKitStorage(
		env.ImageKitPrivateKey,
		env.ImageKitPrivateKey,
		env.ImageKitEndpoint,
	)
	// repositories and usecases
	userRepo := repositories.NewUserRepository(db, env.UserCollection)
	userUsecase := usecases.NewUserUsecase(userRepo, imageKitStorageService, ctxTimeout)
	userController := controllers.NewUserController(userUsecase)

	group.PATCH("/users/update-profile", middleware.AuthMiddleware(*env), userController.UpdateProfile)

}
