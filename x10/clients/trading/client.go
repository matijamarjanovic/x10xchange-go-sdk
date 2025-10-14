package trading

import (
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients"
	pub "github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients/public"
)

// TradingClient is the main authenticated client. It wraps PublicClient for public endpoints
// and adds access to private (authenticated) endpoints using the provided API key.
type TradingClient struct {
	*pub.PublicClient
	httpClient *clients.HTTPClient
	streaming  bool
}

// NewTradingClient creates a new TradingClient. Credentials are required for private endpoints.
// The streaming flag can be used later to opt-in to WebSocket features when available.
func NewTradingClient(cfg *x10.Config, apiKey string, enableStreaming bool) *TradingClient {
	return &TradingClient{
		PublicClient: pub.NewPublicClient(cfg, enableStreaming),
		httpClient:   clients.NewHTTPClientWithAPIKey(cfg, apiKey),
		streaming:    enableStreaming,
	}
}

// StreamingEnabled returns whether streaming features are enabled on this client.
func (c *TradingClient) StreamingEnabled() bool {
	return c.streaming
}
