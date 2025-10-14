package starknet

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/shopspring/decimal"
)

// GenerateNonce generates a random nonce (0 to 2^31 - 1, like Python SDK)
// This matches the Python SDK's nonce generation logic
func GenerateNonce() (int64, error) {
	// Generate a random number between 0 and 2^31 - 1
	max := big.NewInt(1 << 31) // 2^31
	nonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce.Int64(), nil
}

// Constants from Python SDK
const (
	OPLimitOrderWithFees = 3
	HoursInDay           = 24
	SecondsInHour        = 60 * 60
)

// CreateOrderHash creates the order hash using market data
// This implements the same logic as Python SDK's hash_order function
func CreateOrderHash(market *info.Market, orderType, side, qty, price, fee string, expireAtMs int64, nonce int64, vaultID int) (*big.Int, error) {
	// 1. Extract asset IDs and resolutions from market l2Config
	syntheticAssetID, collateralAssetID, syntheticResolution, collateralResolution, err := extractAssetData(market)
	if err != nil {
		return nil, fmt.Errorf("failed to extract asset data: %w", err)
	}

	// 2. Convert decimal amounts to Stark amounts using settlement resolution
	qtyStark, err := convertToStarkAmount(qty, syntheticResolution)
	if err != nil {
		return nil, fmt.Errorf("failed to convert qty to Stark amount: %w", err)
	}

	feeStark, err := convertToStarkAmount(fee, collateralResolution)
	if err != nil {
		return nil, fmt.Errorf("failed to convert fee to Stark amount: %w", err)
	}

	// 3. Determine if buying synthetic (1) or selling (0)
	isBuyingSynthetic := side == "BUY"

	// 4. Calculate collateral amount (qty * price) in human-readable form first
	qtyDecimal, _ := decimal.NewFromString(qty)
	priceDecimal, _ := decimal.NewFromString(price)
	collateralAmountHuman := qtyDecimal.Mul(priceDecimal) // 0.001 * 50000 = 50

	// Convert to Stark amount using collateral resolution
	collateralAmountStark := collateralAmountHuman.Mul(decimal.NewFromInt(int64(collateralResolution)))
	collateralStark := collateralAmountStark.BigInt()

	// 5. Calculate expiration timestamp in hours (like Python SDK)
	expireTime := time.UnixMilli(expireAtMs)
	expireTimeWithBuffer := expireTime.AddDate(0, 0, 14) // Add 14 days buffer

	// 6. Use the exact hashing logic from get_limit_order_msg_without_bounds
	// Position ID should be vault ID, not collateral asset ID
	positionID := int64(vaultID)

	// Convert expiration to seconds instead of hours
	expireTimeInSeconds := expireTimeWithBuffer.Unix()

	hash := getLimitOrderMsg(
		syntheticAssetID,
		collateralAssetID,
		isBuyingSynthetic,
		collateralAssetID, // Fee asset is same as collateral
		qtyStark,
		collateralStark,
		feeStark,
		nonce,
		positionID,
		expireTimeInSeconds, // Use seconds instead of hours
	)

	return hash, nil
}

// getLimitOrderMsg implements the exact logic from Python SDK's get_limit_order_msg_without_bounds
func getLimitOrderMsg(
	assetIDSynthetic, assetIDCollateral *big.Int,
	isBuyingSynthetic bool,
	assetIDFee *big.Int,
	amountSynthetic, amountCollateral, maxAmountFee *big.Int,
	nonce, positionID, expirationTimestamp int64,
) *big.Int {
	var assetIDSell, assetIDBuy, amountSell, amountBuy *big.Int

	if isBuyingSynthetic {
		assetIDSell, assetIDBuy = assetIDCollateral, assetIDSynthetic
		amountSell, amountBuy = amountCollateral, amountSynthetic
	} else {
		assetIDSell, assetIDBuy = assetIDSynthetic, assetIDCollateral
		amountSell, amountBuy = amountSynthetic, amountCollateral
	}

	// First hash: asset_id_sell, asset_id_buy
	// Use HashPedersenElements for two elements to match Python SDK's pedersen_hash(a, b)
	msg := curve.HashPedersenElements([]*big.Int{assetIDSell, assetIDBuy})

	// Second hash: msg, asset_id_fee
	msg = curve.HashPedersenElements([]*big.Int{msg, assetIDFee})

	packedMessage0 := new(big.Int).Set(amountSell)
	packedMessage0.Lsh(packedMessage0, 64)                // packed_message0 * 2^64
	packedMessage0.Add(packedMessage0, amountBuy)         // + amount_buy
	packedMessage0.Lsh(packedMessage0, 64)                // * 2^64
	packedMessage0.Add(packedMessage0, maxAmountFee)      // + max_amount_fee
	packedMessage0.Lsh(packedMessage0, 32)                // * 2^32
	packedMessage0.Add(packedMessage0, big.NewInt(nonce)) // + nonce

	// Third hash: msg, packed_message0
	msg = curve.HashPedersenElements([]*big.Int{msg, packedMessage0})

	packedMessage1 := big.NewInt(OPLimitOrderWithFees)
	packedMessage1.Lsh(packedMessage1, 64)                              // packed_message1 * 2^64
	packedMessage1.Add(packedMessage1, big.NewInt(positionID))          // + position_id
	packedMessage1.Lsh(packedMessage1, 64)                              // * 2^64
	packedMessage1.Add(packedMessage1, big.NewInt(positionID))          // + position_id
	packedMessage1.Lsh(packedMessage1, 64)                              // * 2^64
	packedMessage1.Add(packedMessage1, big.NewInt(positionID))          // + position_id
	packedMessage1.Lsh(packedMessage1, 32)                              // * 2^32
	packedMessage1.Add(packedMessage1, big.NewInt(expirationTimestamp)) // + expiration_timestamp
	packedMessage1.Lsh(packedMessage1, 17)                              // * 2^17 (Padding)

	// Final hash: msg, packed_message1
	return curve.HashPedersenElements([]*big.Int{msg, packedMessage1})
}
