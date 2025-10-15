package trading

import (
	"context"
	"fmt"
	"time"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/utils/starknet"
	"github.com/shopspring/decimal"
)

// OrderOptions controls optional order settings beyond the simplest defaults.
// Use these to tweak reduce/post-only flags, time-in-force, expiry and builder fees.
type OrderOptions struct {
	ReduceOnly  bool
	PostOnly    bool
	TimeInForce string        // default GTT
	ExpireIn    time.Duration // default 24h
	BuilderFee  string
	BuilderID   int
}

// PlaceOrder mirrors the Python SDK shape: market, amount_of_synthetic, price, side.
// This is the simplest entrypoint for placing a LIMIT order:
// - Defaults to LIMIT + GTT with a 24h expiry
// - Picks a sensible default fee (maker if post-only is not set here)
// - Handles signing internally via the configured Signer
// If you need to customize TIF/expiry/flags, use CreateLimitOrder instead.
func (c *TradingClient) PlaceOrder(ctx context.Context, market string, amountOfSynthetic decimal.Decimal, price decimal.Decimal, side string) (*user.CreateOrderResponse, error) {
	return c.CreateLimitOrder(ctx, market, side, amountOfSynthetic.String(), price.String(), &OrderOptions{TimeInForce: "GTT", ExpireIn: 24 * time.Hour})
}

// ReplaceOrder mirrors PlaceOrder but performs replacement via cancelId.
func (c *TradingClient) ReplaceOrder(ctx context.Context, cancelID, market string, amountOfSynthetic decimal.Decimal, price decimal.Decimal, side string) (*user.CreateOrderResponse, error) {
	return c.ReplaceLimitOrder(ctx, cancelID, market, side, amountOfSynthetic.String(), price.String(), &OrderOptions{TimeInForce: "GTT", ExpireIn: 24 * time.Hour})
}

// CreateLimitOrder is a higher-control helper. It still abstracts settlement/signing,
// but lets you customize time-in-force (GTT/FOK/IOC), expiry window, reduce-only,
// post-only and optional builderFee/builderId. It computes a default fee if not provided
// (maker for post-only, otherwise taker) and signs the request via the Signer.
func (c *TradingClient) CreateLimitOrder(ctx context.Context, market, side, qty, price string, opts *OrderOptions) (*user.CreateOrderResponse, error) {
	return c.buildAndSubmitLimitOrder(ctx, market, side, qty, price, opts, "")
}

// ReplaceLimitOrder is a wrapper that sets cancelId and delegates to the shared builder.
func (c *TradingClient) ReplaceLimitOrder(ctx context.Context, cancelID, market, side, qty, price string, opts *OrderOptions) (*user.CreateOrderResponse, error) {
	return c.buildAndSubmitLimitOrder(ctx, market, side, qty, price, opts, cancelID)
}

// CreateOrder submits a fully-formed order request to the API.
// The request must include all required fields including settlement signature and nonce.
// Users should build and sign the CreateOrderRequest themselves before calling this method.
func (c *TradingClient) CreateOrder(ctx context.Context, req user.CreateOrderRequest) (*user.CreateOrderResponse, error) {
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

// EditOrder replaces an existing order by setting CancelID and delegating to CreateOrder.
// The request must be fully-formed and signed by the user.
func (c *TradingClient) EditOrder(ctx context.Context, cancelID string, req user.CreateOrderRequest) (*user.CreateOrderResponse, error) {
	req.CancelID = cancelID
	return c.CreateOrder(ctx, req)
}

// shared helper for create/replace variants
func (c *TradingClient) buildAndSubmitLimitOrder(ctx context.Context, market, side, qty, price string, opts *OrderOptions, cancelID string) (*user.CreateOrderResponse, error) {
	if opts == nil {
		opts = &OrderOptions{TimeInForce: "GTT", ExpireIn: 24 * time.Hour}
	}
	if opts.TimeInForce == "" {
		opts.TimeInForce = "GTT"
	}
	expireMs := time.Now().Add(opts.ExpireIn).UnixMilli()

	feeRate := opts.BuilderFee
	if feeRate == "" {
		if opts.PostOnly {
			feeRate = "0.0002" // 0.02% maker fee (2/10000)
		} else {
			feeRate = "0.0005" // 0.05% taker fee (5/10000)
		}
	}

	fee := feeRate

	if c.account == nil {
		return nil, fmt.Errorf("no Starknet account configured on TradingClient")
	}

	// 1. Fetch market data for the specific market (with caching)
	marketData, err := c.fetchMarketData(ctx, market)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market data for %s: %w", market, err)
	}

	// 2. Generate nonce (0 to 2^31 - 1, like Python SDK)
	nonce, err := starknet.GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 3. Create order hash using market data
	orderHash, err := starknet.CreateOrderHash(marketData, "LIMIT", side, qty, price, fee, expireMs, nonce, c.account.Vault)
	if err != nil {
		return nil, fmt.Errorf("failed to create order hash: %w", err)
	}

	// 4. Sign the hash using the embedded account (just like curve.Sign)
	r, s, err := c.account.Sign(orderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign order: %w", err)
	}

	// 5. Create settlement from signature components
	settlement := user.Settlement{
		Signature: user.SettlementSignature{
			R: fmt.Sprintf("0x%x", r),
			S: fmt.Sprintf("0x%x", s),
		},
		StarkKey:           c.account.GetPublicKeyHex(),
		CollateralPosition: c.account.GetVaultIDString(), //todo
	}

	nonceStr := fmt.Sprintf("%d", nonce)

	req := user.CreateOrderRequest{
		ID:                       fmt.Sprintf("%d", time.Now().UnixNano()),
		Market:                   market,
		Type:                     "LIMIT",
		Side:                     side,
		Qty:                      qty,
		Price:                    price,
		TimeInForce:              opts.TimeInForce,
		ExpiryEpochMillis:        expireMs,
		Fee:                      fee,
		Settlement:               settlement,
		Nonce:                    nonceStr,
		SelfTradeProtectionLevel: "ACCOUNT",
		ReduceOnly:               opts.ReduceOnly,
		PostOnly:                 opts.PostOnly,
		BuilderFee:               opts.BuilderFee,
		BuilderID:                opts.BuilderID,
		CancelID:                 cancelID,
	}
	return c.CreateOrder(ctx, req)
}
