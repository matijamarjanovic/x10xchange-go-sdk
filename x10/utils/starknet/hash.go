package starknet

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/shopspring/decimal"
)

const (
	OPLimitOrderWithFees = 3
	HoursInDay           = 24
	SecondsInHour        = 60 * 60
)

// HashOrder constructs the canonical order hash from the provided market and order parameters.
// It derives synthetic and collateral amounts in human units, computes the fee amount from feeRate,
// converts these values to internal integers using the market resolutions, then packs fields and hashes them.
// The feeRate is a decimal string (e.g., "0.0005" for 0.05%).
func HashOrder(market *info.Market, side, qty, price, feeRate string, expireAtMs int64, nonce int64, vaultID int) (*felt.Felt, error) {
	syntheticAssetID, collateralAssetID, syntheticResolution, collateralResolution, err := extractAssetData(market)
	if err != nil {
		return nil, fmt.Errorf("failed to extract asset data: %w", err)
	}

	qtyDecimal, err := decimal.NewFromString(qty)
	if err != nil {
		return nil, fmt.Errorf("invalid qty: %w", err)
	}
	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}
	feeRateDec, err := decimal.NewFromString(feeRate)
	if err != nil {
		return nil, fmt.Errorf("invalid fee rate: %w", err)
	}

	qtyStark := qtyDecimal.Mul(decimal.NewFromInt(int64(syntheticResolution))).BigInt()

	collateralAmountHuman := qtyDecimal.Mul(priceDecimal)
	collateralStark := collateralAmountHuman.Mul(decimal.NewFromInt(int64(collateralResolution))).BigInt()

	feeAmountHuman := feeRateDec.Mul(collateralAmountHuman)
	feeStark := feeAmountHuman.Mul(decimal.NewFromInt(int64(collateralResolution))).BigInt()

	if side != "BUY" && side != "SELL" {
		return nil, fmt.Errorf("invalid side: %s, must be BUY or SELL", side)
	}
	isBuyingSynthetic := side == "BUY"

	expireTime := time.UnixMilli(expireAtMs)
	expireTimeWithBuffer := expireTime.AddDate(0, 0, 14)
	expireTimeInHours := int64(math.Ceil(float64(expireTimeWithBuffer.Unix()) / float64(SecondsInHour)))

	positionID := int64(vaultID)

	hash, err := getLimitOrderMsg(
		syntheticAssetID,
		collateralAssetID,
		isBuyingSynthetic,
		collateralAssetID,
		qtyStark,
		collateralStark,
		feeStark,
		nonce,
		positionID,
		expireTimeInHours,
	)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// getLimitOrderMsg packs the order fields and performs the pairwise Pedersen hash chaining.
func getLimitOrderMsg(
	assetIDSynthetic, assetIDCollateral *big.Int,
	isBuyingSynthetic bool,
	assetIDFee *big.Int,
	amountSynthetic, amountCollateral, maxAmountFee *big.Int,
	nonce, positionID, expirationTimestamp int64,
) (*felt.Felt, error) {
	var assetIDSell, assetIDBuy, amountSell, amountBuy *big.Int

	if isBuyingSynthetic {
		assetIDSell, assetIDBuy = assetIDCollateral, assetIDSynthetic
		amountSell, amountBuy = amountCollateral, amountSynthetic
	} else {
		assetIDSell, assetIDBuy = assetIDSynthetic, assetIDCollateral
		amountSell, amountBuy = amountSynthetic, amountCollateral
	}

	// First hash: asset_id_sell, asset_id_buy (using curve.Pedersen with felt.Felt)
	fsell, err := bigIntToFelt(assetIDSell)
	if err != nil {
		return nil, err
	}
	fbuy, err := bigIntToFelt(assetIDBuy)
	if err != nil {
		return nil, err
	}
	msg := curve.Pedersen(fsell, fbuy)

	// Second hash: msg, asset_id_fee
	ffee, err := bigIntToFelt(assetIDFee)
	if err != nil {
		return nil, err
	}
	msg = curve.Pedersen(msg, ffee)

	packedMessage0 := new(big.Int).Set(amountSell)
	packedMessage0.Lsh(packedMessage0, 64)                // packed_message0 * 2^64
	packedMessage0.Add(packedMessage0, amountBuy)         // + amount_buy
	packedMessage0.Lsh(packedMessage0, 64)                // * 2^64
	packedMessage0.Add(packedMessage0, maxAmountFee)      // + max_amount_fee
	packedMessage0.Lsh(packedMessage0, 32)                // * 2^32
	packedMessage0.Add(packedMessage0, big.NewInt(nonce)) // + nonce

	// Third hash: msg, packed_message0
	fpm0, err := bigIntToFelt(packedMessage0)
	if err != nil {
		return nil, err
	}
	msg = curve.Pedersen(msg, fpm0)

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
	fpm1, err := bigIntToFelt(packedMessage1)
	if err != nil {
		return nil, err
	}
	return curve.Pedersen(msg, fpm1), nil
}

// GenerateNonce returns a uniformly random nonce in the range [0, 2^31).
// This is suitable for order uniqueness and replay protection.
func GenerateNonce() (int64, error) {
	max := big.NewInt(1 << 31)
	nonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce.Int64(), nil
}
