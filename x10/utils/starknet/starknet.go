package starknet

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models"
)

// StarknetAccount represents a Starknet account with signing capabilities
// This is the Go equivalent of Python's StarkPerpetualAccount
type StarknetAccount struct {
	Vault      int
	PrivateKey *big.Int
	PublicKey  *big.Int
	APIKey     string
}

// NewStarknetAccountFromEnv creates a StarknetAccount by loading credentials from environment variables
func NewStarknetAccount() (*StarknetAccount, error) {
	apiKey := os.Getenv("X10_API_KEY")
	publicKeyHex := os.Getenv("X10_PUBLIC_KEY")
	privateKeyHex := os.Getenv("X10_PRIVATE_KEY")
	vaultIDStr := os.Getenv("X10_VAULT_ID")

	if apiKey == "" || publicKeyHex == "" || privateKeyHex == "" || vaultIDStr == "" {
		return nil, &models.X10Error{
			Code:    400,
			Message: "Missing required environment variables",
			Details: "Please set X10_API_KEY, X10_PUBLIC_KEY, X10_PRIVATE_KEY, and X10_VAULT_ID",
		}
	}

	vaultID, err := strconv.Atoi(vaultIDStr)
	if err != nil {
		return nil, &models.X10Error{
			Code:    400,
			Message: "Invalid vault ID",
			Details: "X10_VAULT_ID must be a valid integer",
		}
	}

	privateKey, ok := new(big.Int).SetString(privateKeyHex, 0)
	if !ok {
		return nil, &models.X10Error{
			Code:    400,
			Message: "Invalid private key format",
			Details: fmt.Sprintf("Private key must be a valid hex string: %s", privateKeyHex),
		}
	}

	publicKey, ok := new(big.Int).SetString(publicKeyHex, 0)
	if !ok {
		return nil, &models.X10Error{
			Code:    400,
			Message: "Invalid public key format",
			Details: fmt.Sprintf("Public key must be a valid hex string: %s", publicKeyHex),
		}
	}

	return &StarknetAccount{
		Vault:      vaultID,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		APIKey:     apiKey,
	}, nil
}

// Sign signs a message hash using the account's private key
// This is the Go equivalent of Python's account.sign() method
// Returns r, s signature components just like curve.Sign
func (a *StarknetAccount) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	r, s, err := curve.Sign(msgHash, a.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign message hash: %w", err)
	}
	return r, s, nil
}

// GetPublicKeyHex returns the public key as a hex string
func (a *StarknetAccount) GetPublicKeyHex() string {
	return fmt.Sprintf("0x%x", a.PublicKey)
}

// GetVaultIDString returns the vault ID as a string
func (a *StarknetAccount) GetVaultIDString() string {
	return strconv.Itoa(a.Vault)
}
