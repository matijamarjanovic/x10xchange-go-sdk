package public

import "github.com/shopspring/decimal"

// MarketStats represents market statistics
type MarketStats struct {
	DailyVolume                decimal.Decimal  `json:"dailyVolume"`
	DailyVolumeBase            decimal.Decimal  `json:"dailyVolumeBase"`
	DailyPriceChange           decimal.Decimal  `json:"dailyPriceChange"`
	DailyPriceChangePercentage decimal.Decimal  `json:"dailyPriceChangePercentage"`
	DailyLow                   decimal.Decimal  `json:"dailyLow"`
	DailyHigh                  decimal.Decimal  `json:"dailyHigh"`
	LastPrice                  decimal.Decimal  `json:"lastPrice"`
	AskPrice                   decimal.Decimal  `json:"askPrice"`
	BidPrice                   decimal.Decimal  `json:"bidPrice"`
	MarkPrice                  decimal.Decimal  `json:"markPrice"`
	IndexPrice                 decimal.Decimal  `json:"indexPrice"`
	FundingRate                decimal.Decimal  `json:"fundingRate"`
	NextFundingRate            int64            `json:"nextFundingRate"`
	OpenInterest               decimal.Decimal  `json:"openInterest"`
	OpenInterestBase           decimal.Decimal  `json:"openInterestBase"`
	DeleverageLevels           DeleverageLevels `json:"deleverageLevels"`
}

// DeleverageLevels represents deleverage levels for long and short positions
type DeleverageLevels struct {
	ShortPositions []DeleverageLevel `json:"shortPositions"`
	LongPositions  []DeleverageLevel `json:"longPositions"`
}

// DeleverageLevel represents a single deleverage level
type DeleverageLevel struct {
	Level             int             `json:"level"`
	RankingLowerBound decimal.Decimal `json:"rankingLowerBound"`
}
