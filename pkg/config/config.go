package config

// Config represents the SDK configuration
type Config struct {
	APIBaseURL string
	StreamURL  string
}

// Testnet returns testnet configuration
func Testnet() *Config {
	return &Config{
		APIBaseURL: "https://starknet.sepolia.extended.exchange/api/v1",
		StreamURL:  "wss://starknet.sepolia.extended.exchange/stream.extended.exchange/v1",
	}
}

// Mainnet returns mainnet configuration
func Mainnet() *Config {
	return &Config{
		APIBaseURL: "https://api.extended.exchange/api/v1",
		StreamURL:  "wss://api.extended.exchange/stream.extended.exchange/v1",
	}
}
