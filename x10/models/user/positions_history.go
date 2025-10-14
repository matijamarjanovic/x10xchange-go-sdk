package user

// PositionHistory represents a historical position (open or closed)
type PositionHistory struct {
	ID              int64  `json:"id"`
	AccountID       int    `json:"accountId"`
	Market          string `json:"market"`
	Side            string `json:"side"`     // LONG | SHORT
	ExitType        string `json:"exitType"` // e.g., TRADE, LIQUIDATION
	Leverage        string `json:"leverage"`
	Size            string `json:"size"`
	MaxPositionSize string `json:"maxPositionSize"`
	OpenPrice       string `json:"openPrice"`
	ExitPrice       string `json:"exitPrice"`
	RealisedPnl     string `json:"realisedPnl"`
	CreatedTime     int64  `json:"createdTime"`
	ClosedTime      int64  `json:"closedTime"`
}
