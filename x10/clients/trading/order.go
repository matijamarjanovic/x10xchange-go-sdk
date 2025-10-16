package trading

import (
	"context"
	"fmt"
	"time"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/perpetual"
	"github.com/shopspring/decimal"
)

// PlaceOrder mirrors the Python SDK shape: market, amount_of_synthetic, price, side.
// This is the simplest entrypoint for placing a LIMIT order:
// - Defaults to LIMIT + GTT with a 24h expiry
// - Picks a sensible default fee (maker if post-only is not set here)
// - Handles signing internally via the configured Signer
// If you need to customize TIF/expiry/flags, use CreateLimitOrder instead.
func (c *TradingClient) PlaceOrder(ctx context.Context, market string, amountOfSynthetic decimal.Decimal, price decimal.Decimal, side string) (*user.CreateOrderResponse, error) {
	return c.CreateLimitOrder(ctx, market, side, amountOfSynthetic.String(), price.String(), &perpetual.OrderOptions{TimeInForce: "GTT", ExpireIn: 24 * time.Hour})
}

// ReplaceOrder mirrors PlaceOrder but performs replacement via cancelId.
func (c *TradingClient) ReplaceOrder(ctx context.Context, cancelID, market string, amountOfSynthetic decimal.Decimal, price decimal.Decimal, side string) (*user.CreateOrderResponse, error) {
	return c.ReplaceLimitOrder(ctx, cancelID, market, side, amountOfSynthetic.String(), price.String(), &perpetual.OrderOptions{TimeInForce: "GTT", ExpireIn: 24 * time.Hour})
}

// CreateLimitOrder is a higher-control helper. It still abstracts settlement/signing,
// but lets you customize time-in-force (GTT/FOK/IOC), expiry window, reduce-only,
// post-only and optional builderFee/builderId. It computes a default fee if not provided
// (maker for post-only, otherwise taker) and signs the request via the Signer.
func (c *TradingClient) CreateLimitOrder(ctx context.Context, market, side, qty, price string, opts *perpetual.OrderOptions) (*user.CreateOrderResponse, error) {
	mkt, err := c.FetchMarketData(ctx, market)
	if err != nil {
		return nil, err
	}
	req, err := perpetual.CreateOrder(c.account, mkt, market, side, qty, price, opts, "")
	if err != nil {
		return nil, err
	}
	return c.PlaceOrderPostRequest(ctx, *req) //TODO should this be done from here?
}

// ReplaceLimitOrder is a wrapper that sets cancelId and delegates to the shared builder.
func (c *TradingClient) ReplaceLimitOrder(ctx context.Context, cancelID, market, side, qty, price string, opts *perpetual.OrderOptions) (*user.CreateOrderResponse, error) {
	mkt, err := c.FetchMarketData(ctx, market)
	if err != nil {
		return nil, err
	}
	req, err := perpetual.CreateOrder(c.account, mkt, market, side, qty, price, opts, cancelID)
	if err != nil {
		return nil, err
	}
	return c.PlaceOrderPostRequest(ctx, *req) //TODO should this be done from here?
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
