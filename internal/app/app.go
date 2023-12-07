package app

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"task/internal/config"
	"task/internal/delivery/register_layers"
	v1 "task/internal/delivery/v1"
	"task/internal/server"
	"task/pkg/database/postgresql"
	"time"
)

type App struct {
	ctx context.Context
	log *zap.Logger
	s   *server.Server
	db  postgresql.Storage
}

func NewApp(ctx context.Context, log *zap.Logger) *App {
	return &App{
		ctx: ctx,
		log: log,
	}
}

func (app *App) Run() error {
	cfg := config.GetConfig()

	app.log.Debug("successfully read config")

	postgresClient, err := postgresql.NewPostgresqlClient(cfg, app.log)
	if err != nil {
		app.log.Error("failed to get postgresql client", zap.Error(err))
		return errors.Wrap(err, "failed to get postgresql client")
	}
	app.log.Debug("successfully get postgresql client")

	app.db = postgresClient

	s, err := server.New(cfg, app.log)
	if err != nil {
		app.log.Error("failed to get server", zap.Error(err))
		return errors.Wrap(err, "failed to get server")
	}
	app.s = s

	rGroup := v1.GetRoutes(s.Engine)

	gRepo := register_layers.NewGRepository(postgresClient, app.log)
	gUsecase := register_layers.NewGUsecase(gRepo)
	gHandler := register_layers.NewGDelivery(gUsecase)
	gHandler.RegisterRoutes(rGroup)

	if err := app.Start(); err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func (app *App) Start() error {
	go func() {
		if err := app.s.Start(); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
				app.log.Warn("server is closed")
			default:
				app.log.Fatal("failed to run server", zap.Error(err))
			}
		}
	}()

	return nil
}

func (app *App) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	app.db.Close()
	app.log.Debug("successfully close database")

	app.s.Stop(ctx)
	app.log.Debug("successfully server stopped")
}
