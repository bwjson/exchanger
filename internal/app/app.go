package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"exchanger/internal/clients/currency"
	"exchanger/internal/config"
	"exchanger/internal/ports/http"
	"exchanger/internal/ports/worker"
	"exchanger/internal/postgres"
	"exchanger/internal/repository/rates"
	"exchanger/internal/service"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *pgxpool.Pool
	server *http.Server
	worker *worker.Worker
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
	provider := currency.New(cfg.Currency)
	svc := service.New(repo, provider, cfg.Currency.Supported, logger)
	server := http.New(cfg.ServerAddress(), svc, logger)

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		server: server,
		worker: worker.New(svc, cfg.Worker, logger),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()

	go func() {
		if err := a.server.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	var wg sync.WaitGroup
	wg.Go(func() { a.worker.Run(ctx) })

	<-ctx.Done()
	err := a.server.Shutdown()
	wg.Wait()
	return err
}
