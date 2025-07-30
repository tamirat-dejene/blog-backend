package routers

import (
	configs "g6/blog-api/Configs"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	usercase "g6/blog-api/Usecases"

	"g6/blog-api/Infrastructure/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthRoutes(env *configs.Env, api *gin.RouterGroup, db *mongo.Database) {
	authService := security.NewJWTService(
		env.AccTS,
		env.RefTS,
		env.AccTE,
		env.RefTE,
	)

	authController := controllers.AuthController{
		UserUsecase:         usercase.NewUserUsecase(repositories.NewUserRepository(db, "users")),
		AuthService:         authService,
		RefreshTokenUsecase: usercase.NewRefreshTokenUsecase(repositories.NewRefreshTokenRepository(db, "refresh_tokens")),
	}

	api.POST("/register", authController.RegisterRequest)
	api.POST("/login", authController.LoginRequest)
}
