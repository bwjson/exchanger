package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"exchanger/internal/app"
	_ "exchanger/internal/docs"
)

// @title Exchanger API
// @version 1.0
// @description Asynchronous currency rate service.
// @BasePath /
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
