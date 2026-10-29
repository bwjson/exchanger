package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"exchanger/internal/ports/http/handlers/rates"
	"exchanger/internal/ports/http/middlewares"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(addr string, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	rates.New().Init(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           middlewares.Recovery(logger, middlewares.Logging(logger, mux)),
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Info("http server listening", "address", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down http server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
