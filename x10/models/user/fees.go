package user

import "github.com/shopspring/decimal"

type TradingFee struct {
	Market       string          `json:"market"`
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate decimal.Decimal `json:"takerFeeRate"`
}

var DefaultFees = TradingFee{
	Market:       "BTC-USD",
	MakerFeeRate: decimal.NewFromInt(2).Div(decimal.NewFromInt(10000)),
	TakerFeeRate: decimal.NewFromInt(5).Div(decimal.NewFromInt(10000)),
}
