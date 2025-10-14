package private

// Account represents an X10 account
type Account struct {
	AccountID             int      `json:"accountId"`
	AccountIndex          int      `json:"accountIndex"`
	APIKeys               []string `json:"apiKeys"`
	BridgeStarknetAddress string   `json:"bridgeStarknetAddress"`
	Description           string   `json:"description"`
	L2Key                 string   `json:"l2Key"`
	L2Vault               string   `json:"l2Vault"`
	Status                string   `json:"status"`
}
