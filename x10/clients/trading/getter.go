// Opposed to the getter.go file in the public package, this file contains only the private get methods for the trading client.
package trading

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/matijamarjanovic/x10xchange-go-sdk/x10/models/user"
)

// GetAccountInfo retrieves the current authenticated account details.
func (c *TradingClient) GetAccountInfo(ctx context.Context) (*user.Account, error) {
	endpoint := "/user/account/info"

	var response struct {
		Status string       `json:"status"`
		Data   user.Account `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	return &response.Data, nil
}

// GetBalance retrieves the key balance details for the authenticated sub-account.
// Note: API may return 404 if user's balance is 0.
func (c *TradingClient) GetBalance(ctx context.Context) (*user.Balance, error) {
	endpoint := "/user/balance"

	var response struct {
		Status string       `json:"status"`
		Data   user.Balance `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &response.Data, nil
}

// GetPositions returns open positions for the authenticated sub-account filtered by markets and/or side.
// markets is optional variadic list; side can be "LONG" or "SHORT".
func (c *TradingClient) GetPositions(ctx context.Context, side *string, markets ...string) ([]user.Position, error) {
	base := "/user/positions"
	q := url.Values{}
	if len(markets) > 0 {
		for _, m := range markets {
			if m != "" {
				q.Add("market", m)
			}
		}
	}
	if side != nil && *side != "" {
		q.Set("side", *side)
	}
	endpoint := base
	if encoded := q.Encode(); encoded != "" {
		endpoint = base + "?" + encoded
	}

	var response struct {
		Status string          `json:"status"`
		Data   []user.Position `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}
	return response.Data, nil
}

// GetAssetOperations returns history of deposits, withdrawals, and transfers with optional filters.
// typeFilter and statusFilter are optional; cursor/limit enable pagination.
func (c *TradingClient) GetAssetOperations(ctx context.Context, typeFilter, statusFilter *string, cursor *int64, limit *int) ([]user.AssetOperation, *user.Pagination, error) {
	base := "/user/assetOperations"
	q := url.Values{}
	if typeFilter != nil && *typeFilter != "" {
		q.Set("type", *typeFilter)
	}
	if statusFilter != nil && *statusFilter != "" {
		q.Set("status", *statusFilter)
	}
	if cursor != nil {
		q.Set("cursor", strconv.FormatInt(*cursor, 10))
	}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	endpoint := base
	if encoded := q.Encode(); encoded != "" {
		endpoint = base + "?" + encoded
	}

	var response struct {
		Status     string                `json:"status"`
		Data       []user.AssetOperation `json:"data"`
		Pagination user.Pagination       `json:"pagination"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to list asset operations: %w", err)
	}

	return response.Data, &response.Pagination, nil
}

// GetPositionsHistory returns historical positions with optional filters and pagination.
func (c *TradingClient) GetPositionsHistory(ctx context.Context, side *string, cursor *int64, limit *int, markets ...string) ([]user.PositionHistory, *user.Pagination, error) {
	base := "/user/positions/history"
	q := url.Values{}
	if len(markets) > 0 {
		for _, m := range markets {
			if m != "" {
				q.Add("market", m)
			}
		}
	}
	if side != nil && *side != "" {
		q.Set("side", *side)
	}
	if cursor != nil {
		q.Set("cursor", strconv.FormatInt(*cursor, 10))
	}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	endpoint := base
	if encoded := q.Encode(); encoded != "" {
		endpoint = base + "?" + encoded
	}

	var response struct {
		Status     string                 `json:"status"`
		Data       []user.PositionHistory `json:"data"`
		Pagination user.Pagination        `json:"pagination"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to get positions history: %w", err)
	}
	return response.Data, &response.Pagination, nil
}

// GetOrderByID retrieves a single order by its ID for the authenticated sub-account.
func (c *TradingClient) GetOrderByID(ctx context.Context, id int64) (*user.Order, error) {
	endpoint := fmt.Sprintf("/user/orders/%d", id)

	var response struct {
		Status string     `json:"status"`
		Data   user.Order `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}
	return &response.Data, nil
}

// GetOrdersByExternalID retrieves orders by user-provided external ID.
func (c *TradingClient) GetOrdersByExternalID(ctx context.Context, externalID string) ([]user.Order, error) {
	endpoint := fmt.Sprintf("/user/orders/external/%s", url.PathEscape(externalID))

	var response struct {
		Status string       `json:"status"`
		Data   []user.Order `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get orders by external id: %w", err)
	}
	return response.Data, nil
}

// GetOpenOrders returns open orders filtered by markets, type, and/or side.
// markets is optional variadic; typeFilter can be LIMIT | CONDITIONAL | TPSL | TWAP; sideFilter BUY | SELL.
func (c *TradingClient) GetOpenOrders(ctx context.Context, typeFilter, sideFilter *string, markets ...string) ([]user.Order, error) {
	base := "/user/orders"
	q := url.Values{}
	if len(markets) > 0 {
		for _, m := range markets {
			if m != "" {
				q.Add("market", m)
			}
		}
	}
	if typeFilter != nil && *typeFilter != "" {
		q.Set("type", *typeFilter)
	}
	if sideFilter != nil && *sideFilter != "" {
		q.Set("side", *sideFilter)
	}
	endpoint := base
	if encoded := q.Encode(); encoded != "" {
		endpoint = base + "?" + encoded
	}

	var response struct {
		Status string       `json:"status"`
		Data   []user.Order `json:"data"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}
	return response.Data, nil
}

// GetOrdersHistory returns orders history with optional filters and pagination.
// Filters: markets..., typeFilter, sideFilter, ids, externalIDs; pagination via cursor, limit.
func (c *TradingClient) GetOrdersHistory(
	ctx context.Context,
	typeFilter, sideFilter *string,
	ids []int64,
	externalIDs []string,
	cursor *int64,
	limit *int,
	markets ...string,
) ([]user.Order, *user.Pagination, error) {
	base := "/user/orders/history"
	q := url.Values{}
	for _, m := range markets {
		if m != "" {
			q.Add("market", m)
		}
	}
	if typeFilter != nil && *typeFilter != "" {
		q.Set("type", *typeFilter)
	}
	if sideFilter != nil && *sideFilter != "" {
		q.Set("side", *sideFilter)
	}
	for _, id := range ids {
		q.Add("id", strconv.FormatInt(id, 10))
	}
	for _, eid := range externalIDs {
		if eid != "" {
			q.Add("externalId", eid)
		}
	}
	if cursor != nil {
		q.Set("cursor", strconv.FormatInt(*cursor, 10))
	}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	endpoint := base
	if encoded := q.Encode(); encoded != "" {
		endpoint = base + "?" + encoded
	}

	var response struct {
		Status     string          `json:"status"`
		Data       []user.Order    `json:"data"`
		Pagination user.Pagination `json:"pagination"`
	}

	if err := c.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to get orders history: %w", err)
	}
	return response.Data, &response.Pagination, nil
}
