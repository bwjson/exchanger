package rates

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"exchanger/internal/domain/entity"
)

type Service interface {
	CreateUpdate(ctx context.Context, rawPair string) (uuid.UUID, error)
	GetUpdate(ctx context.Context, id uuid.UUID) (entity.RateUpdate, error)
	GetLatest(ctx context.Context, rawPair string) (entity.Rate, error)
}

type Controller struct {
	service Service
	logger  *slog.Logger
}

func New(service Service, logger *slog.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}

func (c *Controller) Init(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/rates/updates", c.createUpdate)
	mux.HandleFunc("GET /api/v1/rates/updates/{id}", c.getUpdate)
	mux.HandleFunc("GET /api/v1/rates/latest", c.getLatest)
}
