package user

// Position represents an open position for the authenticated sub-account
type Position struct {
	ID               int     `json:"id"`
	AccountID        int     `json:"accountId"`
	Market           string  `json:"market"`
	Side             string  `json:"side"` // LONG | SHORT
	Leverage         string  `json:"leverage"`
	Size             string  `json:"size"`
	Value            string  `json:"value"`
	OpenPrice        string  `json:"openPrice"`
	MarkPrice        string  `json:"markPrice"`
	LiquidationPrice string  `json:"liquidationPrice"`
	Margin           string  `json:"margin"`
	UnrealisedPnl    string  `json:"unrealisedPnl"`
	RealisedPnl      string  `json:"realisedPnl"`
	TPTriggerPrice   string  `json:"tpTriggerPrice"`
	TPLimitPrice     string  `json:"tpLimitPrice"`
	SLTriggerPrice   string  `json:"slTriggerPrice"`
	SLLimitPrice     string  `json:"slLimitPrice"`
	ADL              float64 `json:"adl"`
	MaxPositionSize  string  `json:"maxPositionSize"`
	CreatedTime      int64   `json:"createdTime"`
	UpdatedTime      int64   `json:"updatedTime"`
}
