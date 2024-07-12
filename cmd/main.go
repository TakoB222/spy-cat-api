package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"spy-cat-api/env"
	"spy-cat-api/handlers"
	"spy-cat-api/pkg/logger"
	"spy-cat-api/services"
	"syscall"
	"time"
)

func main() {
	config, err := env.LoadConfig("")
	if err != nil {
		return
	}

	environment, err := env.NewBuilder().
		SetConfig(config).
		ConnectToPostgresDB().
		Build()
	if err != nil {
		panic(err)
	}

	handler := handlers.NewHandler(environment)
	server := services.NewServer(config.Http, handler.Init())
	go func() {
		if err := server.Run(); err != nil {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
