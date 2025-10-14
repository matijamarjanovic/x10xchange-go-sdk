package public

import (
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients"
)

// PublicClient provides access to public market data endpoints that don't require authentication.
// Use this client to fetch market information, order books, trades, candles, and other public data.
// PublicClient can only be used to fetch public data, it is uncapable of making POST requests.
type PublicClient struct {
	httpClient *clients.HTTPClient
	streaming  bool
}

func NewPublicClient(cfg *x10.Config, enableStreaming bool) *PublicClient {
	return &PublicClient{
		httpClient: clients.NewHTTPClient(cfg),
		streaming:  enableStreaming,
	}
}

// StreamingEnabled returns whether streaming features are enabled on this client.
func (c *PublicClient) StreamingEnabled() bool {
	return c.streaming
}
