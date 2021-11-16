package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	PgMigrationsPath string `mapstructure:"PG_MIGRATIONS_PATH"`
	PgUser           string `mapstructure:"POSTGRES_USER"`
	PgPassword       string `mapstructure:"POSTGRES_PASSWORD"`
	PgDB             string `mapstructure:"POSTGRES_DB"`
}

var (
	config Config
	once   sync.Once
)

// GetConfig reads config from environment. Once.
func GetConfig() *Config {
	once.Do(func() {
		viper.AutomaticEnv()
		viper.SetConfigFile("./.env")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error initializing configs: %s", err.Error())
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("error unmarshalling configs: %s", err.Error())
		}
	})
	return &config
}
