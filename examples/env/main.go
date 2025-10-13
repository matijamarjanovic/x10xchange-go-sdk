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
	// Load configuration from environment variables
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load credentials from environment variables
	creds, err := models.LoadCredentialsFromEnv()
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("Vault ID: %d\n", creds.VaultID)
	fmt.Printf("API Key: %s...\n", creds.APIKey[:8])

	// Create HTTP client with API key
	httpClient := client.NewHTTPClientWithAPIKey(cfg, creds.APIKey)

	// Try to get account details (proper endpoint from API docs)
	var response models.AccountResponse
	err = httpClient.Get(context.Background(), "/user/account/info", &response)
	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		fmt.Printf("Account Status: %s\n", response.Status)
		fmt.Printf("Account ID: %d\n", response.Data.AccountID)
		fmt.Printf("Description: %s\n", response.Data.Description)
		fmt.Printf("L2 Vault: %s\n", response.Data.L2Vault)
		fmt.Printf("Status: %s\n", response.Data.Status)
		fmt.Printf("API Keys: %d keys\n", len(response.Data.APIKeys))
	}
}
