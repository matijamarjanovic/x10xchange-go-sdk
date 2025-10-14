package public

import "github.com/shopspring/decimal"

// OrderBookEntry represents a single entry in the order book
type OrderBookEntry struct {
	Qty   decimal.Decimal `json:"qty"`
	Price decimal.Decimal `json:"price"`
}

// OrderBook represents the order book for a market
type OrderBook struct {
	Market string           `json:"market"`
	Bid    []OrderBookEntry `json:"bid"`
	Ask    []OrderBookEntry `json:"ask"`
}
