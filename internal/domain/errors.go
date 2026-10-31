package domain

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrInvalidPair         = errors.New("invalid currency pair")
	ErrIdempotencyConflict = errors.New("idempotency key is already used with another pair")
	ErrProviderRejected    = errors.New("rate provider rejected the request")
)
