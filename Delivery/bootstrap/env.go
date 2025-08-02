package bootstrap

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Port         string `mapstructure:"PORT"`
	AppEnv       string `mapstructure:"APP_ENV"`
	DB_Uri       string `mapstructure:"DB_URI"`
	DB_Name      string `mapstructure:"DB_NAME"`
	RTS          string `mapstructure:"REFRESH_TOKEN_SECRET"`
	ATS          string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefTEHours   int    `mapstructure:"REFRESH_TOKEN_EXPIRE_HOURS"`
	AccTEMinutes int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	CtxTSeconds  int    `mapstructure:"CONTEXT_TIMEOUT_SECONDS"`

	// blog post defaults
	Page           int    `mapstructure:"PAGE"`
	PageSize       int    `mapstructure:"PAGE_SIZE"`
	Recency        string `mapstructure:"RECENCY"`
	BlogPostCollection string `mapstructure:"BLOG_POST_COLLECTION"`

	// blog comment defaults
	BlogCommentCollection string `mapstructure:"BLOG_COMMENT_COLLECTION"`
	// blog user reaction defaults
	BlogUserReactionCollection string `mapstructure:"BLOG_USER_REACTION_COLLECTION"`

	// user collection
	UserCollection string `mapstructure:"USER_COLLECTION"`

	// user refresh token collection
	RefreshTokenCollection string `mapstructure:"REFRESH_TOKEN_COLLECTION"`

	// password reset token collection
	PasswordResetCollection string `mapstructure:"PASSWORD_RESET_TOKEN_COLLECTION"`
	// password reset token expiry
	PasswordResetExpiry int `mapstructure:"PASSWORD_RESET_TOKEN_EXPIRE_MINUTES"` // in minutes

	// email configuration
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"` // App Password for Gmail
	ResetURL     string `mapstructure:"RESET_URL"`

	// Gemini AI configuration
	GeminiAPIKey    string `mapstructure:"GEMINI_API_KEY"`
	GeminiModelName string `mapstructure:"GEMINI_MODEL_NAME"`

	// AI Blog Post Collection
	AIBlogPostsCollection string `mapstructure:"AI_BLOG_POSTS_COLLECTION"`
	// AI Blog Post Feedback Collection
	AIBlogPostsFeedbackCollection string `mapstructure:"AI_BLOG_POSTS_FEEDBACK_COLLECTION"`
}

// Viper can be made injectable
func NewEnv(env_file_path string) (*Env, error) {
	v := viper.New()
	v.SetConfigFile(env_file_path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("failed to unmarshal env: %w", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env, nil
}
