package info

import "github.com/shopspring/decimal"

type L2Config struct {
	Type                 string `json:"type"`
	CollateralID         string `json:"collateralId"`
	CollateralResolution int    `json:"collateralResolution"`
	SyntheticID          string `json:"syntheticId"`
	SyntheticResolution  int    `json:"syntheticResolution"`
}

type RiskFactorConfig struct {
	UpperBound decimal.Decimal `json:"upperBound"`
	RiskFactor decimal.Decimal `json:"riskFactor"`
}

func (r *RiskFactorConfig) MaxLeverage() decimal.Decimal {
	if r.RiskFactor.IsZero() {
		return decimal.Zero
	}
	return decimal.NewFromInt(1).Div(r.RiskFactor).Round(2)
}
type TradingConfig struct {
	MinOrderSize        decimal.Decimal    `json:"minOrderSize"`
	MinOrderSizeChange  decimal.Decimal    `json:"minOrderSizeChange"`
	MinPriceChange      decimal.Decimal    `json:"minPriceChange"`
	MaxMarketOrderValue decimal.Decimal    `json:"maxMarketOrderValue"`
	MaxLimitOrderValue  decimal.Decimal    `json:"maxLimitOrderValue"`
	MaxPositionValue    decimal.Decimal    `json:"maxPositionValue"`
	MaxLeverage         decimal.Decimal    `json:"maxLeverage"`
	MaxNumOrders        int                `json:"maxNumOrders"`
	LimitPriceCap       decimal.Decimal    `json:"limitPriceCap"`
	LimitPriceFloor     decimal.Decimal    `json:"limitPriceFloor"`
	RiskFactorConfig    []RiskFactorConfig `json:"riskFactorConfig"`
}
type Market struct {
	Name                     string        `json:"name"`
	AssetName                string        `json:"assetName"`
	AssetPrecision           int           `json:"assetPrecision"`
	Category                 string        `json:"category"`
	CollateralAssetName      string        `json:"collateralAssetName"`
	CollateralAssetPrecision int           `json:"collateralAssetPrecision"`
	Active                   bool          `json:"active"`
	Status                   string        `json:"status"`
	UIName                   string        `json:"uiName"`
	VisibleOnUI              bool          `json:"visibleOnUi"`
	CreatedAt                int64         `json:"createdAt"`
	MarketStats              *MarketStats  `json:"marketStats,omitempty"`
	TradingConfig            TradingConfig `json:"tradingConfig"`
	L2Config                 L2Config      `json:"l2Config"`
}

func (m *Market) SyntheticAsset() Asset {
	return Asset{
		ID:                   1,
		Name:                 m.AssetName,
		Precision:            m.AssetPrecision,
		Active:               m.Active,
		IsCollateral:         false,
		SettlementExternalID: m.L2Config.SyntheticID,
		SettlementResolution: m.L2Config.SyntheticResolution,
		L1ExternalID:         "",
		L1Resolution:         0,
	}
}

func (m *Market) CollateralAsset() Asset {
	return Asset{
		ID:                   2,
		Name:                 m.CollateralAssetName,
		Precision:            m.CollateralAssetPrecision,
		Active:               m.Active,
		IsCollateral:         true,
		SettlementExternalID: m.L2Config.CollateralID,
		SettlementResolution: m.L2Config.CollateralResolution,
		L1ExternalID:         "",
		L1Resolution:         0,
	}
}
