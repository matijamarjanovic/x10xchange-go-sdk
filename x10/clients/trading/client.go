package trading

import (
	"context"
	"fmt"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients"
	pub "github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients/public"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/utils/starknet"
)

// TradingClient is the main authenticated client. It wraps PublicClient for public endpoints
// and adds access to private (authenticated) endpoints using the embedded StarknetAccount.
// This matches the Python SDK's architecture where the account is embedded in the client.
type TradingClient struct {
	*pub.PublicClient
	httpClient *clients.HTTPClient
	streaming  bool
	account    *starknet.StarknetAccount // Embedded Starknet account for signing
	markets    map[string]*info.Market   // Cached market data
}

// NewTradingClient creates a new TradingClient by loading credentials from environment variables.
// This is the main constructor that matches the Python SDK's approach.
func NewTradingClient(cfg *x10.Config, enableStreaming bool) (*TradingClient, error) {
	account, err := starknet.NewStarknetAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to load Starknet account from environment: %w", err)
	}

	return &TradingClient{
		PublicClient: pub.NewPublicClient(cfg, enableStreaming),
		httpClient:   clients.NewHTTPClientWithAPIKey(cfg, account.APIKey),
		streaming:    enableStreaming,
		account:      account,
		markets:      make(map[string]*info.Market), // Initialize market cache
	}, nil
}

// StreamingEnabled returns whether streaming features are enabled on this client.
func (c *TradingClient) StreamingEnabled() bool {
	return c.streaming
}

// GetAccount returns the embedded Starknet account
func (c *TradingClient) GetAccount() *starknet.StarknetAccount {
	return c.account
}

// fetchMarketData fetches market data for a specific market (with caching)
func (c *TradingClient) fetchMarketData(ctx context.Context, marketName string) (*info.Market, error) {
	if market, exists := c.markets[marketName]; exists {
		return market, nil
	}

	markets, err := c.GetMarkets(ctx, marketName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market data for %s: %w", marketName, err)
	}

	if len(markets) == 0 {
		return nil, fmt.Errorf("market %s not found", marketName)
	}

	c.markets[marketName] = &markets[0]
	return &markets[0], nil
}
