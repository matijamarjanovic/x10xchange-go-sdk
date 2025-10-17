package models

import (
	"math/big"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/shopspring/decimal"
)

const (
	RoundingModeBuy  = "ROUND_UP"
	RoundingModeSell = "ROUND_DOWN"
	RoundingModeFee  = "ROUND_UP"
)

type HumanReadableAmount struct {
	Value decimal.Decimal
	Asset info.Asset
}

func (h HumanReadableAmount) ToStarkAmount(roundingMode string) StarkAmount {
	starkValue := h.Asset.ConvertHumanReadableToStarkQuantity(h.Value)

	var rounded decimal.Decimal
	switch roundingMode {
	case RoundingModeBuy:
		rounded = starkValue.RoundUp(0)
	case RoundingModeSell:
		rounded = starkValue.RoundDown(0)
	default:
		rounded = starkValue.Round(0)
	}

	return StarkAmount{
		Value: rounded.BigInt(),
		Asset: h.Asset,
	}
}

type StarkAmount struct {
	Value *big.Int
	Asset info.Asset
}

func (s StarkAmount) ToInternalAmount() HumanReadableAmount {
	converted := s.Asset.ConvertStarkToInternalQuantity(s.Value.Int64())
	return HumanReadableAmount{
		Value: converted,
		Asset: s.Asset,
	}
}

type StarkOrderAmounts struct {
	CollateralAmountInternal HumanReadableAmount
	SyntheticAmountInternal  HumanReadableAmount
	FeeAmountInternal        HumanReadableAmount
	FeeRate                  decimal.Decimal
	RoundingMode             string
}

func NewStarkOrderAmounts(
	market *info.Market,
	syntheticAmount decimal.Decimal,
	price decimal.Decimal,
	feeRate decimal.Decimal,
	isBuyingSynthetic bool,
) StarkOrderAmounts {
	var roundingMode string
	if isBuyingSynthetic {
		roundingMode = RoundingModeBuy
	} else {
		roundingMode = RoundingModeSell
	}

	syntheticAsset := market.SyntheticAsset()
	collateralAsset := market.CollateralAsset()

	collateralAmountHuman := syntheticAmount.Mul(price)

	collateralAmount := HumanReadableAmount{
		Value: collateralAmountHuman,
		Asset: collateralAsset,
	}

	syntheticAmountH := HumanReadableAmount{
		Value: syntheticAmount,
		Asset: syntheticAsset,
	}

	feeAmountHuman := feeRate.Mul(collateralAmountHuman)
	feeAmount := HumanReadableAmount{
		Value: feeAmountHuman,
		Asset: collateralAsset,
	}

	return StarkOrderAmounts{
		CollateralAmountInternal: collateralAmount,
		SyntheticAmountInternal:  syntheticAmountH,
		FeeAmountInternal:        feeAmount,
		FeeRate:                  feeRate,
		RoundingMode:             roundingMode,
	}
}
