package user

// Trade represents a single account trade record
type Trade struct {
	ID          int64  `json:"id"`
	AccountID   int    `json:"accountId"`
	Market      string `json:"market"`
	OrderID     int64  `json:"orderId"`
	ExternalID  string `json:"externalId"`
	Side        string `json:"side"` // BUY | SELL
	Price       string `json:"price"`
	Qty         string `json:"qty"`
	Value       string `json:"value"`
	Fee         string `json:"fee"`
	TradeType   string `json:"tradeType"` // TRADE | LIQUIDATION | DELEVERAGE
	CreatedTime int64  `json:"createdTime"`
	IsTaker     bool   `json:"isTaker"`
}
