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

	// otp usecase and otp repository
	otpRepo := repositories.NewOTPRepository(db, env.OtpCollection)
	otpUsecase := usercase.NewOTPUsecase(otpRepo, emailService, ctxTimeout, time.Duration(env.OtpExpireMinutes)*time.Minute, env.OtpMaximumAttempts, env.SecretSalt)

	// storage services
	imageKitStorageService := storage.NewImageKitStorage(
		env.ImageKitPrivateKey,
		env.ImageKitPrivateKey,
		env.ImageKitEndpoint,
	)

	authController := controllers.AuthController{
		UserUsecase:          usercase.NewUserUsecase(userRepo, imageKitStorageService, ctxTimeout),
		OTP:                  otpUsecase,
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

		auth.GET("/google/login", authController.GoogleLogin)
		auth.GET("/google/callback", authController.GoogleCallback)

	}
	authHead := auth
	authHead.Use(middleware.AuthMiddleware(*env))
	{
		authHead.POST("/verify-email", authController.VerifyEmailRequest)
		authHead.POST("/resend-otp", authController.ResendOTPRequest)
		authHead.PATCH("/verify-otp", authController.VerifyOTPRequest)
		authHead.PATCH("/change-role", authController.ChangeRoleRequest)
	}
}
