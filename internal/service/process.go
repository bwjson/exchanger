package service

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
	"exchanger/internal/domain/enum"
)

const (
	maxAttempts = 3
	retryDelay  = 5 * time.Second
)

func (s *Service) ProcessUpdates(ctx context.Context, limit int) error {
	start := time.Now()
	var processed, failed int

	err := s.repo.ProcessPending(ctx, limit, func(updates []entity.RateUpdate) {
		snapshot, err := s.provider.Latest(ctx, s.currencies)
		if err != nil {
			s.logger.Warn("get latest rates", "error", err)
		}
		for i := range updates {
			if err != nil {
				retry(&updates[i], err)
			} else {
				resolve(&updates[i], snapshot)
			}
			if updates[i].Status == enum.UpdateStatusFailed {
				failed++
				s.logger.Warn("rate update failed", "update_id", updates[i].ID, "pair", updates[i].Pair.String(), "error", updates[i].Error)
			}
		}
		processed = len(updates)
	})
	if err != nil {
		return err
	}

	if processed > 0 {
		s.logger.Info("rate updates processed", "count", processed, "failed", failed, "duration", time.Since(start).String())
	}
	return nil
}

// retry schedules u for another attempt with an exponentially growing delay.
// It marks u as failed when the provider rejected the request or the attempts are exhausted.
func retry(u *entity.RateUpdate, cause error) {
	u.Attempts++
	if errors.Is(cause, domain.ErrProviderRejected) || u.Attempts >= maxAttempts {
		u.Status = enum.UpdateStatusFailed
		u.Error = "rate provider is not available"
		return
	}
	u.Status = enum.UpdateStatusPending
	u.NextAttemptAt = time.Now().Add(retryDelay << (u.Attempts - 1))
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
