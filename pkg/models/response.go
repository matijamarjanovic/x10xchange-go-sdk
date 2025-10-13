package models

// SimpleResponse represents a basic API response for testing
type SimpleResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// MarketsResponse represents the response from the markets endpoint
type MarketsResponse struct {
	Status string   `json:"status"`
	Data   []Market `json:"data,omitempty"`
}

// AccountResponse represents the response from the account endpoint
type AccountResponse struct {
	Status string  `json:"status"`
	Data   Account `json:"data,omitempty"`
}
