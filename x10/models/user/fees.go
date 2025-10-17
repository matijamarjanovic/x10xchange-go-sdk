package user

type Fee struct {
	Market         string `json:"market"`
	MakerFeeRate   string `json:"makerFeeRate"`
	TakerFeeRate   string `json:"takerFeeRate"`
}
