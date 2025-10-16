package starknet

import (
	"fmt"
	"math/big"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
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

// bigIntToFelt converts a big.Int to a felt.Felt, returning an error if the value
// is negative or greater than or equal to the Stark field prime. No modulo is applied.
func bigIntToFelt(x *big.Int) (*felt.Felt, error) {
	if x == nil {
		return nil, fmt.Errorf("nil big.Int")
	}
	p, _ := new(big.Int).SetString("0x800000000000011000000000000000000000000000000000000000000000001", 0)
	if x.Sign() < 0 || x.Cmp(p) >= 0 {
		return nil, fmt.Errorf("value %s out of field range", x.String())
	}
	f := new(felt.Felt)
	f.SetBigInt(x)
	return f, nil
}
