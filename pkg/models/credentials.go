package models

import (
	"os"
	"strconv"
)

// Credentials represents X10 API credentials
type Credentials struct {
	APIKey     string
	PublicKey  string
	PrivateKey string
	VaultID    int
}

// LoadCredentialsFromEnv loads credentials from environment variables
func LoadCredentialsFromEnv() (*Credentials, error) {
	apiKey := os.Getenv("X10_API_KEY")
	publicKey := os.Getenv("X10_PUBLIC_KEY")
	privateKey := os.Getenv("X10_PRIVATE_KEY")
	vaultIDStr := os.Getenv("X10_VAULT_ID")

	if apiKey == "" || publicKey == "" || privateKey == "" || vaultIDStr == "" {
		return nil, &X10Error{
			Code:    400,
			Message: "Missing required environment variables",
			Details: "Please set X10_API_KEY, X10_PUBLIC_KEY, X10_PRIVATE_KEY, and X10_VAULT_ID",
		}
	}

	vaultID, err := strconv.Atoi(vaultIDStr)
	if err != nil {
		return nil, &X10Error{
			Code:    400,
			Message: "Invalid vault ID",
			Details: "X10_VAULT_ID must be a valid integer",
		}
	}

	return &Credentials{
		APIKey:     apiKey,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		VaultID:    vaultID,
	}, nil
}
