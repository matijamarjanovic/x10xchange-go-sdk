package user

// Account represents the authenticated user's account details
type Account struct {
	Status                string `json:"status"`
	L2Key                 string `json:"l2Key"`
	L2Vault               string `json:"l2Vault"`
	AccountID             int    `json:"accountId"`
	Description           string `json:"description"`
	BridgeStarknetAddress string `json:"bridgeStarknetAddress"`
}
