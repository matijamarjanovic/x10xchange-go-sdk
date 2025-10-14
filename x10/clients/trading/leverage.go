package trading

import (
	"context"
	"fmt"
)

type updateLeverageRequest struct {
	Market   string `json:"market"`
	Leverage string `json:"leverage"`
}

// UpdateLeverage updates leverage for an individual market.
func (c *TradingClient) UpdateLeverage(ctx context.Context, market string, leverage string) error {
	endpoint := "/user/leverage"
	req := updateLeverageRequest{Market: market, Leverage: leverage}

	var response struct {
		Status string `json:"status"`
	}

	if err := c.httpClient.Patch(ctx, endpoint, req, &response); err != nil {
		return fmt.Errorf("failed to update leverage: %w", err)
	}

	if response.Status != "OK" {
		return fmt.Errorf("failed to update leverage: status=%s", response.Status)
	}
	return nil
}
