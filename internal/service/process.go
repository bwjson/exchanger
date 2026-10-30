package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"exchanger/internal/domain/entity"
	"exchanger/internal/domain/enum"
)

func (s *Service) ProcessUpdates(ctx context.Context, limit int) error {
	start := time.Now()
	var processed, failed int

	err := s.repo.ProcessPending(ctx, limit, func(updates []entity.RateUpdate) error {
		snapshot, err := s.provider.Latest(ctx, s.currencies)
		if err != nil {
			return err
		}
		for i := range updates {
			resolve(&updates[i], snapshot)
			if updates[i].Status == enum.UpdateStatusFailed {
				failed++
				s.logger.Warn("rate is not available", "update_id", updates[i].ID, "pair", updates[i].Pair.String())
			}
		}
		processed = len(updates)
		return nil
	})
	if err != nil {
		return err
	}

	if processed > 0 {
		s.logger.Info("rate updates processed", "count", processed, "failed", failed, "duration", time.Since(start).String())
	}
	return nil
}

// resolve applies the snapshot to u: it records the cross rate and marks u as done,
// or marks u as failed when the rate cannot be derived.
func resolve(u *entity.RateUpdate, snapshot entity.RateSnapshot) {
	price, ok := crossRate(snapshot, u.Pair)
	if !ok {
		u.Status = enum.UpdateStatusFailed
		u.Error = "rate is not available"
		return
	}
	u.Status = enum.UpdateStatusDone
	u.Price = &price
	u.RateTimestamp = &snapshot.Timestamp
}

// crossRate derives the price of pair as quote/base, where both rates are expressed
// against the snapshot base currency. The result is rounded to 10 decimal places to
// match the database precision. It reports false when either rate is unavailable,
// which also guards against division by zero.
func crossRate(snapshot entity.RateSnapshot, pair entity.Pair) (decimal.Decimal, bool) {
	base, quote := rateOf(snapshot, pair.Base), rateOf(snapshot, pair.Quote)
	if base.IsZero() || quote.IsZero() {
		return decimal.Decimal{}, false
	}
	return quote.DivRound(base, 10), true
}

// rateOf returns the amount of currency per one unit of the snapshot base currency.
// The base currency itself is always 1; an unknown currency yields zero.
func rateOf(snapshot entity.RateSnapshot, currency string) decimal.Decimal {
	if currency == snapshot.Base {
		return decimal.NewFromInt(1)
	}
	return snapshot.Rates[currency]
}
