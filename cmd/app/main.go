package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"task/internal/app"
	"task/internal/logger"
)

var log *zap.Logger

func init() {
	log = logger.CreateLogger()
}

// main is the entry point of the application
// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// @in header
// @name Authorization
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.NewApp(ctx, log)
	if err := application.Run(); err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("Started application")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.Stop(ctx)

	log.Debug("shutting down server...")

	log.Info("Stopped application")
}
