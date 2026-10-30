package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"exchanger/internal/config"
	"exchanger/internal/ports/http"
	"exchanger/internal/postgres"
	"exchanger/internal/repository/rates"
	"exchanger/internal/service"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *pgxpool.Pool
	server *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	db, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	repo := rates.New(db)
	svc := service.New(repo, cfg.Currency.Supported)
	server := http.New(cfg.ServerAddress(), svc, logger)

	return &App{cfg: cfg, logger: logger, db: db, server: server}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()

	go func() {
		if err := a.server.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	return a.server.Shutdown()
}
