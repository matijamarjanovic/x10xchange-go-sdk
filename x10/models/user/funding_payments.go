package user

// FundingPayment represents a single funding payment record
type FundingPayment struct {
	ID          int64  `json:"id"`
	AccountID   int    `json:"accountId"`
	Market      string `json:"market"`
	PositionID  int64  `json:"positionId"`
	Side        string `json:"side"` // LONG | SHORT
	Size        string `json:"size"`
	Value       string `json:"value"`
	MarkPrice   string `json:"markPrice"`
	FundingFee  string `json:"fundingFee"`
	FundingRate string `json:"fundingRate"`
	PaidTime    int64  `json:"paidTime"`
}
