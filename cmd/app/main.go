package main

import (
	"context"
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/handler"
	"github.com/gavrl/app/internal/repository"
	"github.com/gavrl/app/internal/service"
	"github.com/gavrl/app/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	util.GetConfig()
	util.InitLogger()
}

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: viper.GetString("POSTGRES_USER"),
		DBName:   viper.GetString("POSTGRES_DB"),
		SSLMode:  viper.GetString("POSTGRES_SSL_MODE"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(internal.Server)

	go func() {
		if err := srv.Run(viper.GetString("HTTP_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Println("Customer-balance-service Started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Customer-balance-service Shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
