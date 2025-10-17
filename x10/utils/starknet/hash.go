package starknet

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models"
)

const (
	OPLimitOrderWithFees = 3
	HoursInDay           = 24
	SecondsInHour        = 60 * 60
)

//todo: add godocs 
func HashOrder(amounts models.StarkOrderAmounts, isBuyingSynthetic bool, expireTime *time.Time, nonce int64, vaultID int) (*felt.Felt, error) {
	syntheticStark := amounts.SyntheticAmountInternal.ToStarkAmount(amounts.RoundingMode)
	collateralStark := amounts.CollateralAmountInternal.ToStarkAmount(amounts.RoundingMode)
	feeStark := amounts.FeeAmountInternal.ToStarkAmount(models.RoundingModeFee)

	syntheticAsset := syntheticStark.Asset
	collateralAsset := collateralStark.Asset

	syntheticAssetID, ok := new(big.Int).SetString(syntheticAsset.SettlementExternalID, 0)
	if !ok {
		return nil, fmt.Errorf("invalid synthetic asset ID: %s", syntheticAsset.SettlementExternalID)
	}

	collateralAssetID, ok := new(big.Int).SetString(collateralAsset.SettlementExternalID, 0)
	if !ok {
		return nil, fmt.Errorf("invalid collateral asset ID: %s", collateralAsset.SettlementExternalID)
	}

	qtyStark := syntheticStark.Value
	collateralStarkBig := collateralStark.Value
	feeStarkBig := feeStark.Value

	expireTimeWithBuffer := expireTime.AddDate(0, 0, 14)
	expireTimeInHours := int64(math.Ceil(float64(expireTimeWithBuffer.Unix()) / float64(SecondsInHour)))

	positionID := int64(vaultID)

	hash, err := getLimitOrderMsg(
		syntheticAssetID,
		collateralAssetID,
		isBuyingSynthetic,
		collateralAssetID,
		qtyStark,
		collateralStarkBig,
		feeStarkBig,
		nonce,
		positionID,
		expireTimeInHours,
	)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

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
