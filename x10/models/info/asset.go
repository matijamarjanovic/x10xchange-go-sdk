package info

import (
	"github.com/shopspring/decimal"
)

// Asset represents an asset with its settlement and L1 configuration
type Asset struct {
	ID                   int
	Name                 string
	Precision            int
	Active               bool
	IsCollateral         bool
	SettlementExternalID string
	SettlementResolution int
	L1ExternalID         string
	L1Resolution         int
}

// ConvertHumanReadableToStarkQuantity converts a human-readable amount to Stark internal units
func (a *Asset) ConvertHumanReadableToStarkQuantity(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(decimal.NewFromInt(int64(a.SettlementResolution)))
}

// ConvertStarkToInternalQuantity converts Stark internal units to human-readable amount
func (a *Asset) ConvertStarkToInternalQuantity(stark int64) decimal.Decimal {
	return decimal.NewFromInt(stark).Div(decimal.NewFromInt(int64(a.SettlementResolution)))
}
