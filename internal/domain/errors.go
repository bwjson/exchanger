package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidPair = errors.New("invalid currency pair")
)
