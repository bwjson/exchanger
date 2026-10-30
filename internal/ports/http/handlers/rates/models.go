package rates

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"exchanger/internal/domain/entity"
	"exchanger/internal/domain/enum"
)

type createUpdateRequest struct {
	Pair string `json:"pair"` // "EUR/MXN"
}

type createUpdateResponse struct {
	UpdateID uuid.UUID `json:"update_id"`
}

type updateResponse struct {
	UpdateID      uuid.UUID         `json:"update_id"`
	Pair          string            `json:"pair"`
	Status        enum.UpdateStatus `json:"status"`
	Price         *decimal.Decimal  `json:"price"`
	RateTimestamp *time.Time        `json:"rate_timestamp"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Error         string            `json:"error,omitempty"`
}

type rateResponse struct {
	Pair          string          `json:"pair"`
	Price         decimal.Decimal `json:"price"`
	RateTimestamp time.Time       `json:"rate_timestamp"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func newUpdateResponse(u entity.RateUpdate) updateResponse {
	resp := updateResponse{
		UpdateID:  u.ID,
		Pair:      u.Pair.String(),
		Status:    u.Status,
		Price:     u.Price,
		UpdatedAt: u.UpdatedAt.UTC(),
		Error:     u.Error,
	}
	if u.RateTimestamp != nil {
		ts := u.RateTimestamp.UTC()
		resp.RateTimestamp = &ts
	}
	return resp
}

func newRateResponse(r entity.Rate) rateResponse {
	return rateResponse{
		Pair:          r.Pair.String(),
		Price:         r.Price,
		RateTimestamp: r.RateTimestamp.UTC(),
		UpdatedAt:     r.UpdatedAt.UTC(),
	}
}
