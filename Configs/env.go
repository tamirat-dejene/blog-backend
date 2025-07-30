package Configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	MongoUri       string `mapstructure:"MONGO_URI"`
	DBName         string `mapstructure:"DB_NAME"`
	RefTS          string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccTS          string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefTE          int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccTE          int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	ContextTimeout int    `mapstructure:"CTX_TIMEOUT"`
}

// Viper can be made injectable
func NewEnv(fileName string) (*Env, error) {
	v := viper.New()
	v.SetConfigFile(fileName)

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
