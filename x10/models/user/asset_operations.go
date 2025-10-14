package user

// AssetOperation represents a deposit, withdrawal, or transfer record
type AssetOperation struct {
	ID                  string `json:"id"`
	Type                string `json:"type"`   // DEPOSIT | WITHDRAWAL | TRANSFER | CLAIM
	Status              string `json:"status"` // COMPLETED | IN_PROGRESS | REJECTED | COMPLETED
	Amount              string `json:"amount"`
	Fee                 string `json:"fee"`
	Asset               int    `json:"asset"`
	Time                int64  `json:"time"` // epoch milliseconds
	AccountID           int    `json:"accountId"`
	CounterpartyAccount int    `json:"counterpartyAccountId,omitempty"`
	TransactionHash     string `json:"transactionHash,omitempty"`
	Chain               string `json:"chain,omitempty"` // e.g., ETH
}

// Pagination represents pagination metadata
type Pagination struct {
	Cursor int64 `json:"cursor"`
	Count  int   `json:"count"`
}
