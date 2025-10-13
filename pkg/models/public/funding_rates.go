package public

import "github.com/shopspring/decimal"

// FundingRate represents a single funding rate record
type FundingRate struct {
	Market    string          `json:"m"` // Market symbol
	Timestamp int64           `json:"T"` // Timestamp in epoch seconds
	Rate      decimal.Decimal `json:"f"` // Funding rate
}

// FundingRatesResponse represents the response from the funding rates endpoint
type FundingRatesResponse struct {
	Data       []FundingRate `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Cursor int64 `json:"cursor"`
	Count  int   `json:"count"`
}
