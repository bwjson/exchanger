package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"exchanger/internal/domain/enum"
)

type Pair struct {
	Base  string
	Quote string
}

type RateUpdate struct {
	ID            uuid.UUID
	Pair          Pair
	Status        enum.UpdateStatus
	Price         *decimal.Decimal
	RateTimestamp *time.Time
	Error         string
	UpdatedAt     time.Time
}

type Rate struct {
	Pair          Pair
	Price         decimal.Decimal
	RateTimestamp time.Time
	UpdatedAt     time.Time
}
