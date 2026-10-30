package service

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"exchanger/internal/domain/entity"
)

func TestCrossRate(t *testing.T) {
	snapshot := entity.RateSnapshot{
		Base: "USD",
		Rates: map[string]decimal.Decimal{
			"EUR": decimal.RequireFromString("0.5"),
			"MXN": decimal.RequireFromString("20"),
		},
	}

	cases := map[entity.Pair]string{
		{Base: "EUR", Quote: "MXN"}: "40",
		{Base: "USD", Quote: "EUR"}: "0.5",
		{Base: "MXN", Quote: "USD"}: "0.05",
	}

	for pair, want := range cases {
		got, ok := crossRate(snapshot, pair)
		require.True(t, ok, pair.String())
		assert.True(t, got.Equal(decimal.RequireFromString(want)), pair.String())
	}
}

func TestCrossRateUnknownCurrency(t *testing.T) {
	snapshot := entity.RateSnapshot{Base: "USD"}

	_, ok := crossRate(snapshot, entity.Pair{Base: "USD", Quote: "EUR"})

	assert.False(t, ok)
}
