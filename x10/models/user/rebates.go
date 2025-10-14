package user

// RebatesStats contains rebate-related metrics for the authenticated sub-account
type RebatesStats struct {
	TotalPaid          string `json:"totalPaid"`
	RebatesRate        string `json:"rebatesRate"`
	MarketShare        string `json:"marketShare"`
	NextTierMakerShare string `json:"nextTierMakerShare"`
	NextTierRebateRate string `json:"nextTierRebateRate"`
}
