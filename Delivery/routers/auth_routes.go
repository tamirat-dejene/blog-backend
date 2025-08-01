package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	usercase "g6/blog-api/Usecases"
	"time"

	"g6/blog-api/Infrastructure/security"

	"g6/blog-api/Infrastructure/database/mongo"
	"github.com/gin-gonic/gin"
	"g6/blog-api/Infrastructure/middleware"
)

func NewAuthRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	authService := security.NewJWTService(
		env.ATS,
		env.RTS,
		env.AccTEMinutes,
		env.RefTEHours,
	)

	authController := controllers.AuthController{
		UserUsecase:         usercase.NewUserUsecase(repositories.NewUserRepository(db, env.UserCollection), time.Duration(env.CtxTSeconds)*time.Second),
		AuthService:         authService,
		RefreshTokenUsecase: usercase.NewRefreshTokenUsecase(repositories.NewRefreshTokenRepository(db, env.RefreshTokenCollection)),
	}

	api.POST("/register", authController.RegisterRequest)
	api.POST("/login", authController.LoginRequest)
	//protected routes
	api.POST("/logout", middleware.AuthMiddleware(*env), authController.LogoutRequest)

}
