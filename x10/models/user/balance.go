package user

// Balance represents key balance details for the authenticated sub-account
type Balance struct {
	CollateralName         string `json:"collateralName"`         // collateral asset symbol (e.g., USDC)
	Balance                string `json:"balance"`                // deposits - withdrawals + realised PnL
	Equity                 string `json:"equity"`                 // balance + unrealised PnL
	AvailableForTrade      string `json:"availableForTrade"`      // equity - initial margin requirement
	AvailableForWithdrawal string `json:"availableForWithdrawal"` // max(0, wallet + min(0, unrealisedPnL) - initialMargin)
	UnrealisedPnl          string `json:"unrealisedPnl"`          // mark-price-based unrealised PnL across positions
	InitialMargin          string `json:"initialMargin"`          // total initial margin requirement
	MarginRatio            string `json:"marginRatio"`            // maintenance margin / equity
	Exposure               string `json:"exposure"`               // sum of all positions' values
	Leverage               string `json:"leverage"`               // exposure / equity
	UpdatedTime            int64  `json:"updatedTime"`            // last update time (epoch seconds)
}
