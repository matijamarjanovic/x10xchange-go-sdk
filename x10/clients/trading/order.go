package trading

import (
	"context"
	"fmt"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/perpetual"
	"github.com/shopspring/decimal"
)

// PlaceOrder creates and submits a LIMIT order, matching Python's place_order method.
// This is the main entrypoint for placing orders on the exchange.
func (c *TradingClient) PlaceOrder(ctx context.Context, market string, amountOfSynthetic decimal.Decimal, price decimal.Decimal, side string, opts *perpetual.PlaceOrderOptions) (*user.CreateOrderResponse, error) {
	if c.account == nil {
		return nil, fmt.Errorf("stark account is not set")
	}

	mkt, err := c.FetchMarketData(ctx, market)
	if err != nil {
		return nil, err
	}

	req, err := perpetual.CreateOrder(c.account, mkt, amountOfSynthetic, price, side, opts)
	if err != nil {
		return nil, err
	}

	return c.PlaceOrderPostRequest(ctx, *req)
}

// CreateOrder submits a fully-formed order request to the API.
// The request must include all required fields including settlement signature and nonce.
// Users should build and sign the CreateOrderRequest themselves before calling this method.
func (c *TradingClient) PlaceOrderPostRequest(ctx context.Context, req user.CreateOrderRequest) (*user.CreateOrderResponse, error) {
	endpoint := "/user/order"

	var response struct {
		Status string                   `json:"status"`
		Data   user.CreateOrderResponse `json:"data"`
	}

	if err := c.httpClient.Post(ctx, endpoint, req, &response); err != nil {
		return nil, fmt.Errorf("failed to create/edit order: %w", err)
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf("failed to create/edit order: status=%s", response.Status)
	}
	return &response.Data, nil
}
