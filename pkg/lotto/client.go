package lotto

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	baseURL = "https://developers.lotto.pl/api/open/v1"
)

type Client interface {
	GetLastResults(ctx context.Context, gameType string) ([]Draw, error)
	GetGameInfo(ctx context.Context, gameType string) (*GameInfo, error)
}

type client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey: apiKey,
	}
}

func (c *client) prepareRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("secret", c.apiKey)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *client) GetLastResults(ctx context.Context, gameType string) ([]Draw, error) {
	url := fmt.Sprintf("%s/lotteries/draw-results/last-results-per-game?gameType=%s", baseURL, gameType)

	req, err := c.prepareRequest(ctx, "GET", url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "status", resp.StatusCode, "body", resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var results []Draw
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return results, nil
}

func (c *client) GetGameInfo(ctx context.Context, gameType string) (*GameInfo, error) {
	url := fmt.Sprintf("%s/lotteries/info?gameType=%s", baseURL, gameType)

	req, err := c.prepareRequest(ctx, "GET", url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "status", resp.StatusCode, "body", resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info GameInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}
