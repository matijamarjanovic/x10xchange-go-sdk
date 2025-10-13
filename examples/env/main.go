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
	var response models.Response
	err = httpClient.Get(context.Background(), "/user/account/info", &response)
	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		fmt.Printf("Account Status: %s\n", response.Status)

		// Cast the data to account info
		if accountData, ok := response.Data.(map[string]interface{}); ok {
			fmt.Printf("Account ID: %v\n", accountData["accountId"])
			fmt.Printf("Description: %v\n", accountData["description"])
			fmt.Printf("L2 Vault: %v\n", accountData["l2Vault"])
			fmt.Printf("Status: %v\n", accountData["status"])
			if apiKeys, ok := accountData["apiKeys"].([]interface{}); ok {
				fmt.Printf("API Keys: %d keys\n", len(apiKeys))
			}
		}
	}
}
