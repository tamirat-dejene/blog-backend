package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	usercase "g6/blog-api/Usecases"
	"time"

	"g6/blog-api/Infrastructure/email"
	"g6/blog-api/Infrastructure/middleware"
	"g6/blog-api/Infrastructure/security"
	"g6/blog-api/Infrastructure/storage"

	"g6/blog-api/Infrastructure/database/mongo"

	"github.com/gin-gonic/gin"
)

func NewAuthRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
	// context time out
	ctxTimeout := time.Duration(env.CtxTSeconds) * time.Second

	// jwt services
	authService := security.NewJWTService(
		env.ATS,
		env.RTS,
		env.AccTEMinutes,
		env.RefTEHours,
	)

	// user repository
	userRepo := repositories.NewUserRepository(db, env.UserCollection)

	// email service
	emailService := email.NewGomailEmailService(
		env.SMTPHost,
		env.SMTPPort,
		env.SMTPFrom,
		env.SMTPUsername,
		env.SMTPPassword,
	)

	// reset password repository
	resetPasswordRepo := repositories.NewPasswordResetRepository(db, env.PasswordResetCollection)

	// password reset usecase
	passwordResetUsecase := usercase.NewPasswordResetUsecase(
		resetPasswordRepo,
		userRepo,
		emailService,
		time.Duration(env.PasswordResetExpiry)*time.Minute,
	)

	// storage services
	imageKitStorageService := storage.NewImageKitStorage(
		env.ImageKitPrivateKey,
		env.ImageKitPrivateKey,
		env.ImageKitEndpoint,
	)

	authController := controllers.AuthController{
		UserUsecase:          usercase.NewUserUsecase(userRepo, imageKitStorageService, ctxTimeout),
		AuthService:          authService,
		RefreshTokenUsecase:  usercase.NewRefreshTokenUsecase(repositories.NewRefreshTokenRepository(db, env.RefreshTokenCollection)),
		PasswordResetUsecase: passwordResetUsecase,
		Env:                  env,
	}

	auth := api.Group("/auth/")
	{
		auth.POST("/register", authController.RegisterRequest)
		auth.POST("/login", authController.LoginRequest)
		auth.POST("/logout", authController.LogoutRequest)
		auth.POST("/forgot-password", authController.ForgotPasswordRequest)
		auth.POST("/reset-password", authController.ResetPasswordRequest)
		auth.POST("/refresh", authController.RefreshToken)
		auth.PATCH("/change-role", middleware.AuthMiddleware(*env), authController.ChangeRoleRequest)
		auth.GET("/google/login",authController.GoogleLogin)
		auth.GET("/google/callback",authController.GoogleCallback)
	}
}