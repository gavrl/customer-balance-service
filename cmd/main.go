package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

}

// initConfig reads config file
func initConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile("./.env")

	return viper.ReadInConfig()
}
