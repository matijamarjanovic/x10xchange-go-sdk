package user

// SettlementSignature holds Starknet signature values
type SettlementSignature struct {
	R string `json:"r"`
	S string `json:"s"`
}

// Settlement contains signature, stark key and collateral position
type Settlement struct {
	Signature          SettlementSignature `json:"signature"`
	StarkKey           string              `json:"starkKey"`
	CollateralPosition string              `json:"collateralPosition"`
}

// Trigger defines conditional trigger params
type Trigger struct {
	TriggerPrice       string `json:"triggerPrice,omitempty"`
	TriggerPriceType   string `json:"triggerPriceType,omitempty"`
	Direction          string `json:"direction,omitempty"`
	ExecutionPriceType string `json:"executionPriceType,omitempty"`
}

// TpslConfig defines TP/SL legs
type TpslConfig struct {
	TriggerPrice     string      `json:"triggerPrice,omitempty"`
	TriggerPriceType string      `json:"triggerPriceType,omitempty"`
	Price            string      `json:"price,omitempty"`
	PriceType        string      `json:"priceType,omitempty"`
	Settlement       *Settlement `json:"settlement,omitempty"`
}

// CreateOrderRequest payload for POST /user/order
type CreateOrderRequest struct {
	ID                       string      `json:"id"`
	Market                   string      `json:"market"`
	Type                     string      `json:"type"` // LIMIT | MARKET | CONDITIONAL | TPSL
	Side                     string      `json:"side"` // BUY | SELL
	Qty                      string      `json:"qty"`
	Price                    string      `json:"price"`
	TimeInForce              string      `json:"timeInForce"` // GTT | FOK | IOC
	ExpiryEpochMillis        int64       `json:"expiryEpochMillis"`
	Fee                      string      `json:"fee"`
	CancelID                 string      `json:"cancelId,omitempty"`
	Settlement               Settlement  `json:"settlement"`
	Nonce                    string      `json:"nonce"`
	SelfTradeProtectionLevel string      `json:"selfTradeProtectionLevel"` // DISABLED | ACCOUNT | CLIENT
	ReduceOnly               bool        `json:"reduceOnly,omitempty"`
	PostOnly                 bool        `json:"postOnly,omitempty"`
	Trigger                  *Trigger    `json:"trigger,omitempty"`
	TpSlType                 string      `json:"tpSlType,omitempty"` // ORDER | POSITION
	TakeProfit               *TpslConfig `json:"takeProfit,omitempty"`
	StopLoss                 *TpslConfig `json:"stopLoss,omitempty"`
	BuilderFee               string      `json:"builderFee,omitempty"`
	BuilderID                int         `json:"builderId,omitempty"`
}

// CreateOrderResponse returned on success
type CreateOrderResponse struct {
	ID         int64  `json:"id"`
	ExternalID string `json:"externalId"`
}
