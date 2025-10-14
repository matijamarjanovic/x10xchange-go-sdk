package info

import "github.com/shopspring/decimal"

// OpenInterest represents a single open interest record
type OpenInterest struct {
	Interest  decimal.Decimal `json:"i"` // Open interest amount
	Interest2 decimal.Decimal `json:"I"` // Secondary open interest amount
	Timestamp int64           `json:"t"` // Timestamp in epoch milliseconds
}
