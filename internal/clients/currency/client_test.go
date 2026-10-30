package currency_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"exchanger/internal/clients/currency"
	"exchanger/internal/config"
)

func TestLatest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"response":{"date":"2024-03-15T14:30:00Z","base":"USD","rates":{"MXN":18.5437788018}}}`))
	}))
	defer srv.Close()

	client := currency.New(config.CurrencyConfig{BaseURL: srv.URL})

	snapshot, err := client.Latest(context.Background(), []string{"MXN"})

	require.NoError(t, err)
	require.Equal(t, "USD", snapshot.Base)
	require.Equal(t, "18.5437788018", snapshot.Rates["MXN"].String())
}

func TestLatestBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := currency.New(config.CurrencyConfig{BaseURL: srv.URL})

	_, err := client.Latest(context.Background(), []string{"MXN"})

	require.Error(t, err)
}
