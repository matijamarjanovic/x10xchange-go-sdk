package starknet

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
)

// GenerateNonce generates a random nonce (0 to 2^31 - 1, like Python SDK)
// This matches the Python SDK's nonce generation logic
func GenerateNonce() (int64, error) {
	// Generate a random number between 0 and 2^31 - 1
	max := big.NewInt(1 << 31) // 2^31
	nonce, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce.Int64(), nil
}

// CreateOrderHash creates the order hash using market data
// This will implement the same logic as Python SDK's hash_order function
func CreateOrderHash(market *info.Market, orderType, side, qty, price, fee string, expireAtMs int64, nonce int64) (*big.Int, error) {
	// TODO: Implement proper order hash creation using market data
	// This will use the same logic as Python SDK's hash_order function
	// For now, return a placeholder hash
	return big.NewInt(12345), nil
}
