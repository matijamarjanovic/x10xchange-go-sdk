package public

import "github.com/shopspring/decimal"

// Candle represents a single candle/OHLCV data point
type Candle struct {
	Open      decimal.Decimal `json:"o"` // Open price
	Low       decimal.Decimal `json:"l"` // Low price
	High      decimal.Decimal `json:"h"` // High price
	Close     decimal.Decimal `json:"c"` // Close price
	Volume    decimal.Decimal `json:"v"` // Volume
	Timestamp int64           `json:"T"` // Timestamp in epoch milliseconds
}
