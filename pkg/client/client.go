package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/x10xchange/go-sdk/pkg/config"
)

// HTTPClient represents the HTTP client
type HTTPClient struct {
	config     *config.Config
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get makes a GET request to the specified endpoint
func (c *HTTPClient) Get(ctx context.Context, endpoint string, result interface{}) error {
	url := c.config.APIBaseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "X10GoSDK/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
