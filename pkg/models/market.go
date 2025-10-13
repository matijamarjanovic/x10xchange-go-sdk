package models

// Market represents a trading market
type Market struct {
	Name                     string                 `json:"name"`
	AssetName                string                 `json:"assetName"`
	AssetPrecision           int                    `json:"assetPrecision"`
	Category                 string                 `json:"category"`
	CollateralAssetName      string                 `json:"collateralAssetName"`
	CollateralAssetPrecision int                    `json:"collateralAssetPrecision"`
	Active                   bool                   `json:"active"`
	Status                   string                 `json:"status"`
	UIName                   string                 `json:"uiName"`
	VisibleOnUI              bool                   `json:"visibleOnUi"`
	CreatedAt                int64                  `json:"createdAt"`
	MarketStats              map[string]interface{} `json:"marketStats"`
	TradingConfig            map[string]interface{} `json:"tradingConfig"`
	L2Config                 map[string]interface{} `json:"l2Config"`
}
