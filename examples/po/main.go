package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients/trading"
	"github.com/shopspring/decimal"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Load config
	cfg, err := x10.LoadFromEnv()
	if err != nil || cfg == nil {
		cfg = x10.Testnet()
	}

	// Create trading client with embedded account (loads from .env)
	tc, err := trading.NewTradingClient(cfg, false)
	if err != nil {
		log.Fatalf("failed to create trading client: %v", err)
	}

	// Test that the account was loaded correctly
	account := tc.GetAccount()
	fmt.Printf("Testing order creation with embedded Starknet account:\n")
	fmt.Printf("  Public Key: %s\n", account.GetPublicKeyHex())
	fmt.Printf("  Vault ID: %s\n", account.GetVaultIDString())

	// Test order creation parameters
	market := "BTC-USD"
	qty := decimal.NewFromFloat(0.0001) // Small amount for testing
	price := decimal.NewFromFloat(50000)

	fmt.Printf("\nAttempting to place order:\n")
	fmt.Printf("  Market: %s\n", market)
	fmt.Printf("  Quantity: %s BTC\n", qty.String())
	fmt.Printf("  Price: $%s\n", price.String())
	fmt.Printf("  Side: BUY\n")

	// Try to place a buy order
	resp, err := tc.PlaceOrder(context.Background(), market, qty, price, "BUY")
	if err != nil {
		log.Printf("Order creation failed: %v", err)
		return
	}

	fmt.Printf("\n✅ Order created successfully!\n")
	fmt.Printf("  Order ID: %d\n", resp.ID)
	fmt.Printf("  External ID: %s\n", resp.ExternalID)
}
