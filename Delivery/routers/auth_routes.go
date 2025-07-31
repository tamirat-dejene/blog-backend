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
)

func NewAuthRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	authService := security.NewJWTService(
		env.ATS,
		env.RTS,
		env.AccTEMinutes,
		env.RefTEHours,
	)

	// reset password repository
	resetPasswordRepo := repositories.NewPasswordResetRepository(db, env.PasswordResetCollection)

	authController := controllers.AuthController{
		UserUsecase:          usercase.NewUserUsecase(repositories.NewUserRepository(db, env.UserCollection), time.Duration(env.CtxTSeconds)*time.Second),
		AuthService:          authService,
		RefreshTokenUsecase:  usercase.NewRefreshTokenUsecase(repositories.NewRefreshTokenRepository(db, env.RefreshTokenCollection)),
		PasswordResetUsecase: usercase.NewPasswordResetUsecase(resetPasswordRepo, time.Duration(env.PasswordResetExpiry)*time.Minute),
	}
	auth := api.Group("/auth/")
	{
		auth.POST("/register", authController.RegisterRequest)
		auth.POST("/login", authController.LoginRequest)
		auth.POST("/logout", authController.LogoutRequest)
		auth.POST("/forgot-password", authController.ForgotPasswordRequest)
		auth.POST("/refresh", authController.RefreshToken)
	}

}
