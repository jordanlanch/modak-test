package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              int    `mapstructure:"REDIS_PORT"`
}

func NewEnv(envfile string) *Env {
	viper.SetConfigFile(envfile)
	if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Can't find the file .env : %s", err)
	}

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
			log.Fatalf("Environment can't be loaded: %s", err)
	}

	if env.AppEnv == "development" {
			log.Println("The App is running in development environment")
	}

	return &env
}
