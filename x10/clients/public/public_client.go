package client

import (
	"context"
	"fmt"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/clients"
	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/public"
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

// GetAllMarkets fetches information for all available markets.
// Returns a list of all markets with their basic configuration and status.
func (c *PublicClient) GetAllMarkets(ctx context.Context) ([]public.Market, error) {
	endpoint := "/info/markets"

	var response struct {
		Status string          `json:"status"`
		Data   []public.Market `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get markets: %w", err)
	}

	return response.Data, nil
}

// GetMarkets fetches information for one or more specific markets by name.
// Accepts variadic market names (e.g., "BTC-USD", "ETH-USD").
func (c *PublicClient) GetMarkets(ctx context.Context, markets ...string) ([]public.Market, error) {
	if len(markets) == 0 {
		return nil, fmt.Errorf("at least one market name must be provided")
	}

	endpoint := "/info/markets"

	params := ""
	for i, market := range markets {
		if i == 0 {
			params = fmt.Sprintf("?market=%s", market)
		} else {
			params += fmt.Sprintf("&market=%s", market)
		}
	}
	endpoint += params

	var response struct {
		Status string          `json:"status"`
		Data   []public.Market `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get markets: %w", err)
	}

	return response.Data, nil
}

// GetMarketStats fetches real-time market statistics including prices, volume, and funding rates.
// Returns daily volume, price changes, bid/ask prices, mark price, and deleverage levels.
func (c *PublicClient) GetMarketStats(ctx context.Context, market string) (*public.MarketStats, error) {
	endpoint := fmt.Sprintf("/info/markets/%s/stats", market)

	var response struct {
		Status string             `json:"status"`
		Data   public.MarketStats `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get market stats: %w", err)
	}

	return &response.Data, nil
}

// GetOrderBook fetches the current order book showing bid and ask orders with quantities.
// Returns arrays of bid and ask orders sorted by price.
func (c *PublicClient) GetOrderBook(ctx context.Context, market string) (*public.OrderBook, error) {
	endpoint := fmt.Sprintf("/info/markets/%s/orderbook", market)

	var response struct {
		Status string           `json:"status"`
		Data   public.OrderBook `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	return &response.Data, nil
}

// GetMarketTrades fetches the latest 50 trades for a market.
// Returns trade data including price, quantity, side, and timestamp.
func (c *PublicClient) GetMarketTrades(ctx context.Context, market string) ([]public.Trade, error) {
	endpoint := fmt.Sprintf("/info/markets/%s/trades", market)

	var response struct {
		Status string         `json:"status"`
		Data   []public.Trade `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get market trades: %w", err)
	}

	return response.Data, nil
}

// GetCandles fetches OHLCV candlestick data for a market.
// candleType can be "trades", "mark-prices", or "index-prices".
func (c *PublicClient) GetCandles(ctx context.Context, market, candleType, interval string, limit int, endTime *int64) ([]public.Candle, error) {
	endpoint := fmt.Sprintf("/info/candles/%s/%s", market, candleType)

	params := fmt.Sprintf("?interval=%s&limit=%d", interval, limit)
	if endTime != nil {
		params += fmt.Sprintf("&endTime=%d", *endTime)
	}
	endpoint += params

	var response struct {
		Status string          `json:"status"`
		Data   []public.Candle `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get candles: %w", err)
	}

	return response.Data, nil
}

// GetTradesCandles fetches candlestick data based on actual trade prices.
func (c *PublicClient) GetTradesCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]public.Candle, error) {
	return c.GetCandles(ctx, market, "trades", interval, limit, endTime)
}

// GetMarkPriceCandles fetches candlestick data based on mark prices (fair value prices).
func (c *PublicClient) GetMarkPriceCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]public.Candle, error) {
	return c.GetCandles(ctx, market, "mark-prices", interval, limit, endTime)
}

// GetIndexPriceCandles fetches candlestick data based on index prices (spot market prices).
func (c *PublicClient) GetIndexPriceCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]public.Candle, error) {
	return c.GetCandles(ctx, market, "index-prices", interval, limit, endTime)
}

// GetFundingRates fetches historical funding rates with pagination support.
// Funding rates are applied hourly and returned sorted by timestamp (descending).
func (c *PublicClient) GetFundingRates(ctx context.Context, market string, startTime, endTime int64, cursor *int64, limit *int) (*public.FundingRatesResponse, error) {
	endpoint := fmt.Sprintf("/info/%s/funding?startTime=%d&endTime=%d", market, startTime, endTime)

	if cursor != nil {
		endpoint += fmt.Sprintf("&cursor=%d", *cursor)
	}
	if limit != nil {
		endpoint += fmt.Sprintf("&limit=%d", *limit)
	}

	var response public.FundingRatesResponse
	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get funding rates: %w", err)
	}

	return &response, nil
}

// GetOpenInterest fetches historical open interest data with configurable intervals.
// interval can be "P1H" (hourly) or "P1D" (daily).
func (c *PublicClient) GetOpenInterest(ctx context.Context, market, interval string, startTime, endTime int64, limit *int) ([]public.OpenInterest, error) {
	endpoint := fmt.Sprintf("/info/%s/open-interests?interval=%s&startTime=%d&endTime=%d", market, interval, startTime, endTime)

	if limit != nil {
		endpoint += fmt.Sprintf("&limit=%d", *limit)
	}

	var response struct {
		Status string                `json:"status"`
		Data   []public.OpenInterest `json:"data"`
	}

	err := c.httpClient.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get open interest: %w", err)
	}

	return response.Data, nil
}
