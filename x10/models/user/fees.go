package user

// Fee represents current fee rates for a market
type Fee struct {
	Market         string `json:"market"`
	MakerFeeRate   string `json:"makerFeeRate"`
	TakerFeeRate   string `json:"takerFeeRate"`
	BuilderFeeRate string `json:"builderFeeRate,omitempty"`
}
