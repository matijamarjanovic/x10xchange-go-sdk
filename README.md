# X10 Extended Exchange Go SDK

A Go SDK for the X10 Extended Exchange API, providing comprehensive trading functionality for perpetual futures.

## Features

- **Trading Operations**: Place, cancel, and manage orders
- **Account Management**: Balance, positions, and trades
- **Real-time Data**: WebSocket streams for orderbooks, trades, and account updates
- **StarkEx Integration**: Full support for StarkEx cryptographic operations
- **User Onboarding**: Account creation and management
- **Asset Operations**: Deposits, withdrawals, and transfers

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/matijamarjanovic/x10xchange-go-sdk/pkg/client"
    "github.com/matijamarjanovic/x10xchange-go-sdk/pkg/config"
)

func main() {
    // Create configuration
    cfg := config.Mainnet()
    
    // Create trading client
    tradingClient := client.NewTradingClient(cfg)
    
    // Get markets
    markets, err := tradingClient.Markets().GetMarkets(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Found %d markets", len(markets))
}
```

## Installation

```bash
go get github.com/matijamarjanovic/x10xchange-go-sdk
```

## Documentation

See the [examples](./examples/) directory for usage examples.

## License

MIT
