package starknet

import (
	"fmt"
	"math/big"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/shopspring/decimal"
)

// extractAssetData extracts synthetic and collateral asset IDs along with their settlement resolutions
// from the market's L2 configuration. Used for order hash calculation and amount conversion.
func extractAssetData(market *info.Market) (syntheticAssetID, collateralAssetID *big.Int, syntheticResolution, collateralResolution int, err error) {
	syntheticIDStr, ok := market.L2Config["syntheticId"].(string)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("syntheticId not found in l2Config")
	}

	syntheticAssetID, ok = new(big.Int).SetString(syntheticIDStr, 0)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("invalid synthetic asset ID: %s", syntheticIDStr)
	}

	collateralIDStr, ok := market.L2Config["collateralId"].(string)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("collateralId not found in l2Config")
	}

	collateralAssetID, ok = new(big.Int).SetString(collateralIDStr, 0)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("invalid collateral asset ID: %s", collateralIDStr)
	}

	syntheticResolutionFloat, ok := market.L2Config["syntheticResolution"].(float64)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("syntheticResolution not found in l2Config")
	}
	syntheticResolution = int(syntheticResolutionFloat)

	collateralResolutionFloat, ok := market.L2Config["collateralResolution"].(float64)
	if !ok {
		return nil, nil, 0, 0, fmt.Errorf("collateralResolution not found in l2Config")
	}
	collateralResolution = int(collateralResolutionFloat)

	return syntheticAssetID, collateralAssetID, syntheticResolution, collateralResolution, nil
}

// convertToStarkAmount converts a human-readable decimal amount to Stark's internal integer representation
// by multiplying with the settlement resolution. Returns a felt.Felt for cryptographic operations.
func convertToStarkAmount(decimalStr string, resolution int) (*felt.Felt, error) {
	dec, err := decimal.NewFromString(decimalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal string: %s", decimalStr)
	}

	resolutionDecimal := decimal.NewFromInt(int64(resolution))
	result := dec.Mul(resolutionDecimal)

	resultInt := result.BigInt()

	return bigIntToFelt(resultInt), nil
}

// bigIntToFelt converts a big.Int to a felt.Felt by reducing it modulo the Stark field prime.
// This ensures values stay within the Stark field bounds (uint256-like) and prevents overflow
// in cryptographic operations. Handles negative numbers by adding the field prime.
func bigIntToFelt(x *big.Int) *felt.Felt {
	p, _ := new(big.Int).SetString("0x800000000000011000000000000000000000000000000000000000000000001", 0)
	mod := new(big.Int).Mod(new(big.Int).Set(x), p)
	if mod.Sign() < 0 {
		mod.Add(mod, p)
	}
	f := new(felt.Felt)
	f.SetBigInt(mod)
	return f
}
