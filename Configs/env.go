package config

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

	// blog defaults
	Page           int    `mapstructure:"PAGE"`
	PageSize       int    `mapstructure:"PAGE_SIZE"`
	Recency        string `mapstructure:"RECENCY"`
	BlogCollection string `mapstructure:"BLOG_COLLECTION"`

	// blog comment defaults
	BlogCommentCollection string `mapstructure:"BLOG_COMMENT_COLLECTION"`
	// blog user reaction defaults
	BlogUserReactionCollection string `mapstructure:"BLOG_USER_REACTION_COLLECTION"`

	// user collection
	UserCollection string `mapstructure:"USER_COLLECTION"`

	// user refresh token collection
	RefreshTokenCollection string `mapstructure:"REFRESH_TOKEN_COLLECTION"`
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
