package public

import "github.com/shopspring/decimal"

// Trade represents a single trade in the market
type Trade struct {
	ID        int64           `json:"i"`  // Trade ID
	Market    string          `json:"m"`  // Market symbol
	Side      string          `json:"S"`  // Side: BUY or SELL
	TradeType string          `json:"tT"` // Trade type: TRADE
	Time      int64           `json:"T"`  // Timestamp
	Price     decimal.Decimal `json:"p"`  // Price
	Quantity  decimal.Decimal `json:"q"`  // Quantity
}
