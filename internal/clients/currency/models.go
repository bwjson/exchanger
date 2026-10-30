package currency

import (
	"time"

	"github.com/shopspring/decimal"
)

type latestResponse struct {
	Response latestRates `json:"response"`
}

type latestRates struct {
	Date  time.Time                  `json:"date"`
	Base  string                     `json:"base"`
	Rates map[string]decimal.Decimal `json:"rates"`
}
