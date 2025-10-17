# X10 Extended Exchange Go SDK

**DISCLAIMER**: Existing GO packages don't implement RFC6979 signing as API server expects so development is paused until it becomes available.
Codebase temporarily contains Python SDK for reference during development.

A pure Go SDK for the X10 Extended Exchange API, providing comprehensive trading functionality for perpetual futures. This sdk is written to be as similar to [Extended Python SDK](https://github.com/x10xchange/python_sdk) as possible in terms of functionality and code architecure as much as the language difference allows it. 

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
