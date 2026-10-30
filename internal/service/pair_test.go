package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
)

func TestParsePair(t *testing.T) {
	s := New(nil, []string{"USD", "EUR", "MXN"})

	pair, err := s.parsePair("EUR/MXN")

	require.NoError(t, err)
	require.Equal(t, entity.Pair{Base: "EUR", Quote: "MXN"}, pair)
}

func TestParsePairInvalid(t *testing.T) {
	s := New(nil, []string{"USD", "EUR", "MXN"})

	for _, raw := range []string{"EUR/EUR", "EUR/GBP", "eur/mxn", "EURMXN", "EUR/MXN/USD", ""} {
		_, err := s.parsePair(raw)

		require.ErrorIs(t, err, domain.ErrInvalidPair, raw)
	}
}
