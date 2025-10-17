package starknet

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	felt "github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
)

// StarknetAccount represents a Starknet account with signing capabilities
// This is the Go equivalent of Python's StarkPerpetualAccount
type StarknetPerpetualAccount struct {
	Vault       int
	PrivateKey  *big.Int
	PublicKey   *big.Int
	APIKey      string
	TradingFees map[string]user.Fee
}

// NewStarknetAccountFromEnv creates a StarknetAccount by loading credentials from environment variables
func NewStarknetAccount() (*StarknetPerpetualAccount, error) {
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

	return &StarknetPerpetualAccount{
		Vault:       vaultID,
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		APIKey:      apiKey,
		TradingFees: make(map[string]user.Fee),
	}, nil
}

// Sign signs a message hash using the account's private key
// This is the Go equivalent of Python's account.sign() method
// Returns r, s signature components as *big.Int for easy hex formatting
func (a *StarknetPerpetualAccount) Sign(msgHash *felt.Felt) (*big.Int, *big.Int, error) {
	privateKeyFelt := new(felt.Felt).SetBigInt(a.PrivateKey)
	r, s, err := curve.SignFelts(msgHash, privateKeyFelt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign message hash: %w", err)
	}
	rBigInt := r.BigInt(new(big.Int))
	sBigInt := s.BigInt(new(big.Int))
	return rBigInt, sBigInt, nil
}
