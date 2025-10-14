package starknet

import (
	"fmt"
	"math/big"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/shopspring/decimal"
)

// extractAssetData extracts asset IDs and resolutions from market l2Config
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

// convertToStarkAmount converts a decimal string to Stark amount using resolution
func convertToStarkAmount(decimalStr string, resolution int) (*big.Int, error) {
	dec, err := decimal.NewFromString(decimalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal string: %s", decimalStr)
	}

	resolutionDecimal := decimal.NewFromInt(int64(resolution))
	result := dec.Mul(resolutionDecimal)

	resultInt := result.BigInt()

	return resultInt, nil
}
