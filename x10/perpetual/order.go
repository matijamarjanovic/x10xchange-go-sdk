package perpetual

import (
	"fmt"
	"math/big"
	"time"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/utils/starknet"
	"github.com/shopspring/decimal"
)

// PlaceOrderOptions contains optional parameters for PlaceOrder
type PlaceOrderOptions struct {
	PostOnly                 *bool
	PreviousOrderID          *string
	ExpireTime               *time.Time
	OrderExternalID          *string
	TimeInForce              *string
	SelfTradeProtectionLevel *string
}

// CreateOrder creates an order object to be placed on the exchange.
// This matches Python's place_order function signature and behavior.
func CreateOrder(
	account *starknet.StarknetPerpetualAccount,
	market *info.Market,
	amountOfSynthetic decimal.Decimal,
	price decimal.Decimal,
	side string,
	opts *PlaceOrderOptions,
) (*user.CreateOrderRequest, error) {
	if market == nil {
		return nil, fmt.Errorf("market is required")
	}

	if opts == nil {
		opts = &PlaceOrderOptions{}
	}

	fees := account.TradingFees[market.Name]
	if fees == (user.Fee{}) {
		fees = user.Fee{
			Market:       market.Name,
			TakerFeeRate: "0.00025",
			MakerFeeRate: "0.0000",
		}
	}

	return createOrder(
		market,
		amountOfSynthetic,
		price,
		side,
		account.Vault,
		fees,
		account.Sign,
		account.PublicKey,
		false,
		opts.ExpireTime,
		opts.PostOnly != nil && *opts.PostOnly,
		opts.PreviousOrderID,
		opts.OrderExternalID,
		opts.TimeInForce,
		opts.SelfTradeProtectionLevel,
	)
}

//todo: add godocs + continue matching python sdk
func createOrder(
	market *info.Market,
	syntheticAmount decimal.Decimal,
	price decimal.Decimal,
	side string,
	collateralPositionID int,
	fees user.Fee,
	signer func(*felt.Felt) (*big.Int, *big.Int, error),
	publicKey *big.Int,
	exactOnly bool,
	expireTime *time.Time,
	postOnly bool,
	previousOrderExternalID *string,
	orderExternalID *string,
	timeInForce *string,
	selfTradeProtectionLevel *string,
) (*user.CreateOrderRequest, error) {
	if exactOnly {
		return nil, fmt.Errorf("exact_only option is not supported yet")
	}

	if expireTime == nil {
		defaultExpire := time.Now().Add(8 * time.Hour)
		expireTime = &defaultExpire
	}

	if timeInForce == nil {
		defaultTIF := "GTT"
		timeInForce = &defaultTIF
	}

	if selfTradeProtectionLevel == nil {
		defaultSTP := "ACCOUNT"
		selfTradeProtectionLevel = &defaultSTP
	}

	nonce, err := starknet.GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	expireMs := expireTime.UnixMilli()

	feeRate, err := decimal.NewFromString(fees.TakerFeeRate)
	if err != nil {
		return nil, fmt.Errorf("invalid fee rate: %w", err)
	}

	isBuyingSynthetic := side == "BUY"

	amounts := models.NewStarkOrderAmounts(market, syntheticAmount, price, feeRate, isBuyingSynthetic)

	orderHash, err := starknet.HashOrder(amounts, isBuyingSynthetic, expireMs, nonce, collateralPositionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order hash: %w", err)
	}

	r, s, err := signer(orderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign order: %w", err)
	}

	settlement := user.Settlement{
		Signature: user.SettlementSignature{
			R: fmt.Sprintf("0x%x", r),
			S: fmt.Sprintf("0x%x", s),
		},
		StarkKey:           fmt.Sprintf("0x%x", publicKey),
		CollateralPosition: fmt.Sprintf("%d", collateralPositionID),
	}

	var orderID string
	if orderExternalID != nil {
		orderID = *orderExternalID
	} else {
		orderID = orderHash.String()
	}

	collateralStark, syntheticStark, feeStark := models.ToStarkAmounts(amounts)

	debuggingAmounts := &user.DebuggingAmounts{
		CollateralAmount: collateralStark.Value.String(),
		FeeAmount:        feeStark.Value.String(),
		SyntheticAmount:  syntheticStark.Value.String(),
	}

	req := user.CreateOrderRequest{
		ID:                       orderID,
		Market:                   market.Name,
		Type:                     "LIMIT",
		Side:                     side,
		Qty:                      syntheticAmount.String(),
		Price:                    price.String(),
		PostOnly:                 postOnly,
		TimeInForce:              *timeInForce,
		ExpiryEpochMillis:        expireMs,
		Fee:                      fees.TakerFeeRate,
		SelfTradeProtectionLevel: *selfTradeProtectionLevel,
		Nonce:                    fmt.Sprintf("%d", nonce),
		CancelID:                 getStringValue(previousOrderExternalID),
		Settlement:               settlement,
		DebuggingAmounts:         debuggingAmounts,
	}

	return &req, nil
}

// getStringValue safely extracts string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
