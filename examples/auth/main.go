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

	// Create HTTP client with API key (replace with a real API key for testing)
	apiKey := "your-api-key-here"
	httpClient := client.NewHTTPClientWithAPIKey(cfg, apiKey)

	// Try to get account balance (requires authentication)
	var response models.Response
	err := httpClient.Get(context.Background(), "/account/balance", &response)
	if err != nil {
		log.Printf("Request failed (expected with dummy API key): %v", err)
	} else {
		fmt.Printf("Response: %+v\n", response)
	}
}
