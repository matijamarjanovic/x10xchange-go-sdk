package main

import (
	"context"
	"fmt"
	"log"

	"github.com/x10xchange/go-sdk/pkg/client"
	"github.com/x10xchange/go-sdk/pkg/config"
	"github.com/x10xchange/go-sdk/pkg/models"
)

func main() {
	// Create configuration
	cfg := config.Testnet()

	// Create HTTP client
	httpClient := client.NewHTTPClient(cfg)

	// Make a simple request to get markets
	var response models.MarketsResponse
	err := httpClient.Get(context.Background(), "/info/markets?market=BTC-USD", &response)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Found %d markets\n", len(response.Data))

	if len(response.Data) > 0 {
		market := response.Data[0]
		fmt.Printf("Market: %s\n", market.Name)
		fmt.Printf("Asset: %s\n", market.AssetName)
		fmt.Printf("Active: %t\n", market.Active)
	}
}
