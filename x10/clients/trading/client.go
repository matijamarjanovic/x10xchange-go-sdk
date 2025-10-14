package trading

import (
	"context"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients"
	pub "github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients/public"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
)

// TradingClient is the main authenticated client. It wraps PublicClient for public endpoints
// and adds access to private (authenticated) endpoints using the provided API key.
type TradingClient struct {
	*pub.PublicClient
	httpClient *clients.HTTPClient
	streaming  bool
	signer     Signer
}

// NewTradingClient creates a new TradingClient. Credentials are required for private endpoints.
// The streaming flag can be used later to opt-in to WebSocket features when available.
// The signer is optional - if provided, it enables order creation methods.
func NewTradingClient(cfg *x10.Config, apiKey string, enableStreaming bool, signer Signer) *TradingClient {
	return &TradingClient{
		PublicClient: pub.NewPublicClient(cfg, enableStreaming),
		httpClient:   clients.NewHTTPClientWithAPIKey(cfg, apiKey),
		streaming:    enableStreaming,
		signer:       signer,
	}
}

// StreamingEnabled returns whether streaming features are enabled on this client.
func (c *TradingClient) StreamingEnabled() bool {
	return c.streaming
}

// SignInputs contains normalized fields provided to the Signer to produce
// Starknet settlement and a nonce. Callers never set these directly; they're
// derived internally from the high-level order parameters.
type SignInputs struct {
	Market     string
	Type       string
	Side       string
	Qty        string
	Price      string
	Fee        string
	ExpireAtMs int64
	PostOnly   bool
}

// Signer provides order-settlement signatures (Starknet) for order creation.
type Signer interface {
	// SignOrder should produce settlement and nonce for the given order payload inputs.
	SignOrder(ctx context.Context, inputs SignInputs) (settlement user.Settlement, nonce string, err error)
}

// SetSigner configures a signer used to sign orders before submission.
func (c *TradingClient) SetSigner(s Signer) {
	c.signer = s
}
