package x10

import (
	"os"

	"github.com/joho/godotenv"
)

// Config represents the SDK configuration
type Config struct {
	APIBaseURL  string
	StreamURL   string
	Environment string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	godotenv.Load()

	env := getEnvOrDefault("X10_ENVIRONMENT", "testnet")

	if env == "mainnet" {
		return Mainnet(), nil
	}
	return Testnet(), nil
}

// Testnet returns testnet configuration
func Testnet() *Config {
	return &Config{
		APIBaseURL:  "https://api.starknet.sepolia.extended.exchange/api/v1",
		StreamURL:   "wss://starknet.sepolia.extended.exchange/stream.extended.exchange/v1",
		Environment: "testnet",
	}
}

// Mainnet returns mainnet configuration
func Mainnet() *Config {
	return &Config{
		APIBaseURL:  "https://api.starknet.extended.exchange/api/v1",
		StreamURL:   "wss://api.starknet.extended.exchange/stream.extended.exchange/v1",
		Environment: "mainnet",
	}
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
