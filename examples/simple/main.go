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
	var response models.Response
	err := httpClient.Get(context.Background(), "/info/markets?market=BTC-USD", &response)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Printf("Status: %s\n", response.Status)

	// Cast the data to markets
	if markets, ok := response.Data.([]interface{}); ok {
		fmt.Printf("Found %d markets\n", len(markets))

		if len(markets) > 0 {
			if marketData, ok := markets[0].(map[string]interface{}); ok {
				fmt.Printf("Market: %v\n", marketData["name"])
				fmt.Printf("Asset: %v\n", marketData["assetName"])
				fmt.Printf("Active: %v\n", marketData["active"])
			}
		}
	}
}
