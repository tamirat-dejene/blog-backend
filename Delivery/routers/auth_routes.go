package routers

import (
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/controllers"
	repositories "g6/blog-api/Repositories"
	usercase "g6/blog-api/Usecases"
	"time"

	"g6/blog-api/Infrastructure/email"
	"g6/blog-api/Infrastructure/security"

	"g6/blog-api/Infrastructure/database/mongo"

	"g6/blog-api/Infrastructure/middleware"
	"github.com/gin-gonic/gin"
)

func NewAuthRoutes(env *bootstrap.Env, api *gin.RouterGroup, db mongo.Database) {
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

	authController := controllers.AuthController{
		UserUsecase:          usercase.NewUserUsecase(repositories.NewUserRepository(db, env.UserCollection), time.Duration(env.CtxTSeconds)*time.Second),
		AuthService:          authService,
		RefreshTokenUsecase:  usercase.NewRefreshTokenUsecase(repositories.NewRefreshTokenRepository(db, env.RefreshTokenCollection)),
		PasswordResetUsecase: passwordResetUsecase,
	}

	auth := api.Group("/auth/")
	{
		auth.POST("/register", authController.RegisterRequest)
		auth.POST("/login", authController.LoginRequest)
		//protected routes
		auth.POST("/logout", middleware.AuthMiddleware(*env), authController.LogoutRequest)
		auth.POST("/forgot-password", authController.ForgotPasswordRequest)
		auth.POST("/reset-password", authController.ResetPasswordRequest)
		auth.POST("/refresh", authController.RefreshToken)
		auth.PATCH("/change-role", middleware.AuthMiddleware(*env), authController.ChangeRoleRequest)
	}

}
