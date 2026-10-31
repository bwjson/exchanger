package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"exchanger/internal/config"
	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
)

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New(cfg config.CurrencyConfig) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) Latest(ctx context.Context, symbols []string) (entity.RateSnapshot, error) {
	endpoint := c.baseURL + "/latest?symbols=" + strings.Join(symbols, ",")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return entity.RateSnapshot{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return entity.RateSnapshot{}, fmt.Errorf("get latest rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if retryable(resp.StatusCode) {
			return entity.RateSnapshot{}, fmt.Errorf("get latest rates: status %d", resp.StatusCode)
		}
		return entity.RateSnapshot{}, fmt.Errorf("%w: status %d", domain.ErrProviderRejected, resp.StatusCode)
	}

	var out latestResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return entity.RateSnapshot{}, fmt.Errorf("decode latest rates: %w", err)
	}

	return entity.RateSnapshot{
		Base:      out.Response.Base,
		Rates:     out.Response.Rates,
		Timestamp: out.Response.Date,
	}, nil
}

func retryable(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}
