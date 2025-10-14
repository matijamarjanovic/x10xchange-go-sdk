// Wrappers for public client methods
package trading

import (
	"context"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/info"
)

func (c *TradingClient) GetAllMarkets(ctx context.Context) ([]info.Market, error) {
	return c.PublicClient.GetAllMarkets(ctx)
}

func (c *TradingClient) GetMarkets(ctx context.Context, markets ...string) ([]info.Market, error) {
	return c.PublicClient.GetMarkets(ctx, markets...)
}

func (c *TradingClient) GetMarketStats(ctx context.Context, market string) (*info.MarketStats, error) {
	return c.PublicClient.GetMarketStats(ctx, market)
}

func (c *TradingClient) GetOrderBook(ctx context.Context, market string) (*info.OrderBook, error) {
	return c.PublicClient.GetOrderBook(ctx, market)
}

func (c *TradingClient) GetMarketTrades(ctx context.Context, market string) ([]info.Trade, error) {
	return c.PublicClient.GetMarketTrades(ctx, market)
}

func (c *TradingClient) GetCandles(ctx context.Context, market, candleType, interval string, limit int, endTime *int64) ([]info.Candle, error) {
	return c.PublicClient.GetCandles(ctx, market, candleType, interval, limit, endTime)
}

func (c *TradingClient) GetTradesCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]info.Candle, error) {
	return c.PublicClient.GetTradesCandles(ctx, market, interval, limit, endTime)
}

func (c *TradingClient) GetMarkPriceCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]info.Candle, error) {
	return c.PublicClient.GetMarkPriceCandles(ctx, market, interval, limit, endTime)
}

func (c *TradingClient) GetIndexPriceCandles(ctx context.Context, market, interval string, limit int, endTime *int64) ([]info.Candle, error) {
	return c.PublicClient.GetIndexPriceCandles(ctx, market, interval, limit, endTime)
}

func (c *TradingClient) GetFundingRates(ctx context.Context, market string, startTime, endTime int64, cursor *int64, limit *int) (*info.FundingRatesResponse, error) {
	return c.PublicClient.GetFundingRates(ctx, market, startTime, endTime, cursor, limit)
}

func (c *TradingClient) GetOpenInterest(ctx context.Context, market, interval string, startTime, endTime int64, limit *int) ([]info.OpenInterest, error) {
	return c.PublicClient.GetOpenInterest(ctx, market, interval, startTime, endTime, limit)
}
