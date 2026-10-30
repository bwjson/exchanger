package worker

import (
	"context"
	"log/slog"
	"time"

	"exchanger/internal/config"
)

type Service interface {
	ProcessUpdates(ctx context.Context, limit int) error
}

type Worker struct {
	service   Service
	interval  time.Duration
	batchSize int
	logger    *slog.Logger
}

func New(service Service, cfg config.WorkerConfig, logger *slog.Logger) *Worker {
	return &Worker{service: service, interval: cfg.Interval, batchSize: cfg.BatchSize, logger: logger}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("rate worker started", "interval", w.interval.String(), "batch_size", w.batchSize)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("rate worker stopped")
			return
		case <-ticker.C:
			if err := w.service.ProcessUpdates(ctx, w.batchSize); err != nil {
				w.logger.Error("process rate updates", "error", err)
			}
		}
	}
}
