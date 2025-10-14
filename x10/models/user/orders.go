package user

// TriggerConfig describes trigger-related fields for conditional/TPSL orders
type TriggerConfig struct {
	TriggerPrice          string `json:"triggerPrice"`
	TriggerPriceType      string `json:"triggerPriceType"`      // LAST | INDEX | MARK
	TriggerPriceDirection string `json:"triggerPriceDirection"` // UP | DOWN
	ExecutionPriceType    string `json:"executionPriceType"`    // MARKET | LIMIT (where applicable)
}

// TpConfig describes take profit parameters
type TpConfig struct {
	TriggerPrice     string `json:"triggerPrice"`
	TriggerPriceType string `json:"triggerPriceType"`
	Price            string `json:"price"`
	PriceType        string `json:"priceType"`
}

// SlConfig describes stop loss parameters
type SlConfig struct {
	TriggerPrice     string `json:"triggerPrice"`
	TriggerPriceType string `json:"triggerPriceType"`
	Price            string `json:"price"`
	PriceType        string `json:"priceType"`
}

// Order represents an open order
type Order struct {
	ID           int64          `json:"id"`
	AccountID    int            `json:"accountId"`
	ExternalID   string         `json:"externalId"`
	Market       string         `json:"market"`
	Type         string         `json:"type"` // LIMIT | CONDITIONAL | TPSL | TWAP
	Side         string         `json:"side"` // BUY | SELL
	Status       string         `json:"status"`
	Price        string         `json:"price"`
	AveragePrice string         `json:"averagePrice"`
	Qty          string         `json:"qty"`
	FilledQty    string         `json:"filledQty"`
	PayedFee     string         `json:"payedFee"`
	Trigger      *TriggerConfig `json:"trigger,omitempty"`
	TakeProfit   *TpConfig      `json:"takeProfit,omitempty"`
	StopLoss     *SlConfig      `json:"stopLoss,omitempty"`
	ReduceOnly   bool           `json:"reduceOnly"`
	PostOnly     bool           `json:"postOnly"`
	CreatedTime  int64          `json:"createdTime"`
	UpdatedTime  int64          `json:"updatedTime"`
	TimeInForce  string         `json:"timeInForce"`
	ExpireTime   int64          `json:"expireTime"`
}
