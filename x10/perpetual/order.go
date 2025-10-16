package perpetual

import (
	"fmt"
	"time"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/utils/starknet"
)

// OrderClient defines the minimal client surface needed to build and submit orders.
// OrderOptions mirrors trading.OrderOptions to avoid import cycles.
type OrderOptions struct {
	ReduceOnly  bool
	PostOnly    bool
	TimeInForce string
	ExpireIn    time.Duration
	BuilderFee  string
	BuilderID   int
}

// CreateOrder builds, signs, and submits a LIMIT order using the provided TradingClient,
// mirroring Python's order_object flow. Supports cancelID for replacements.
func CreateOrder(
	account *starknet.StarknetAccount,
	marketData *info.Market,
	market, side, qty, price string,
	opts *OrderOptions,
	cancelID string,
) (*user.CreateOrderRequest, error) {
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
			feeRate = "0.0002"
		} else {
			feeRate = "0.0005"
		}
	}
	fee := feeRate

	if account == nil {
		return nil, fmt.Errorf("no Starknet account configured on TradingClient")
	}

	if marketData == nil {
		return nil, fmt.Errorf("market data is required")
	}

	nonce, err := starknet.GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	orderHash, err := starknet.HashOrder(marketData, side, qty, price, fee, expireMs, nonce, account.Vault)
	if err != nil {
		return nil, fmt.Errorf("failed to create order hash: %w", err)
	}

	r, s, err := account.Sign(orderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign order: %w", err)
	}

	settlement := user.Settlement{
		Signature: user.SettlementSignature{
			R: fmt.Sprintf("0x%x", r),
			S: fmt.Sprintf("0x%x", s),
		},
		StarkKey:           account.GetPublicKeyHex(),
		CollateralPosition: account.GetVaultIDString(),
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
	return &req, nil
}
