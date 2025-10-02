# ThorStreamer Go SDK

Official Go client library for ThorStreamer gRPC services.

## Installation

```bash
go get github.com/thorlabsDev/ThorStreamer/sdks/go@v0.1.0
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    thorclient "github.com/thorlabsDev/ThorStreamer/sdks/go/client"
)

func main() {
    // Load environment variables
    godotenv.Load()

    client, err := thorclient.NewClient(thorclient.Config{
        ServerAddr:     os.Getenv("SERVER_ADDRESS"),
        Token:          os.Getenv("AUTH_TOKEN"),
        DefaultTimeout: 30 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Subscribe to transactions
    stream, err := client.SubscribeToTransactions(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for {
        msg, err := stream.Recv()
        if err != nil {
            if thorclient.IsStreamDone(err) {
                break
            }
            log.Printf("Error: %v", err)
            break
        }

        if tx := msg.GetTransaction(); tx != nil {
            log.Printf("Transaction: slot=%d", tx.Transaction.Slot)
        }
    }
}
```

## Environment Setup

Create a `.env` file:

```bash
SERVER_ADDRESS=your-server:50051
AUTH_TOKEN=your-auth-token
```

## Features

- **Transaction Streaming**: Real-time Solana transaction updates
- **Slot Status**: Subscribe to slot confirmations and updates
- **Wallet Tracking**: Monitor specific wallet addresses (up to 10 per request)
- **Account Updates**: Track account state changes with owner filtering
- **Automatic Reconnection**: Built-in error handling and stream management
- **Type Safety**: Full protobuf type definitions

## API Reference

### Creating a Client

```go
client, err := thorclient.NewClient(thorclient.Config{
    ServerAddr:     "server-address:50051",
    Token:          "your-auth-token",
    DefaultTimeout: 30 * time.Second,
})
```

### Subscribe to Transactions

```go
stream, err := client.SubscribeToTransactions(ctx)
if err != nil {
    log.Fatal(err)
}

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    
    if tx := msg.GetTransaction(); tx != nil {
        log.Printf("Transaction: %x", tx.Transaction.Signature)
    }
}
```

### Subscribe to Slot Status

```go
stream, err := client.SubscribeToSlotStatus(ctx)
if err != nil {
    log.Fatal(err)
}

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    
    if slot := msg.GetSlot(); slot != nil {
        log.Printf("Slot: %d, Status: %d", slot.Slot, slot.Status)
    }
}
```

### Subscribe to Wallet Transactions

Monitor specific wallet addresses (max 10 per request):

```go
wallets := []string{
    "wallet1base58address...",
    "wallet2base58address...",
}
stream, err := client.SubscribeToWalletTransactions(ctx, wallets)
if err != nil {
    log.Fatal(err)
}

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    
    if tx := msg.GetTransaction(); tx != nil {
        log.Printf("Wallet transaction: slot=%d", tx.Transaction.Slot)
    }
}
```

### Subscribe to Account Updates

Track account state changes with optional owner filtering:

```go
accounts := []string{"account1...", "account2..."}
owners := []string{"owner1..."}

stream, err := client.SubscribeToAccountUpdates(ctx, accounts, owners)
if err != nil {
    log.Fatal(err)
}

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    
    if acc := msg.GetAccountUpdate(); acc != nil {
        log.Printf("Account: lamports=%d", acc.Lamports)
    }
}
```


## Error Handling

```go
msg, err := stream.Recv()
if err != nil {
    if thorclient.IsStreamDone(err) {
        // Stream closed normally (EOF or context cancelled)
        return
    }
    // Handle error
    log.Printf("Stream error: %v", err)
    return
}
```

## Examples

See the [examples directory](../../examples/golang-advanced) for complete working examples.

## Building from Source

```bash
# Clone repository
git clone https://github.com/thorlabsDev/ThorStreamer.git
cd ThorStreamer/sdks/go

# Install dependencies
go mod download

# Generate protobuf code (if needed)
cd ../../proto
protoc --go_out=../sdks/go/proto --go_opt=paths=source_relative \
    --go-grpc_out=../sdks/go/proto --go-grpc_opt=paths=source_relative \
    events.proto publisher.proto

# Build the client library
cd ../sdks/go
go build ./client

# Run examples
cd ../../examples/golang-advanced
go run main.go
```

## Module Information

```go
module github.com/thorlabsDev/ThorStreamer/sdks/go

go 1.23.2

require (
    google.golang.org/grpc v1.75.1
    google.golang.org/protobuf v1.36.10
)
```

## Version History

- `v0.1.0` - Initial release
    - Transaction streaming
    - Slot status updates
    - Wallet tracking
    - Account updates

## Contributing

Contributions are welcome!

## Support

- GitHub Issues: https://github.com/thorlabsDev/ThorStreamer/issues
- Documentation: https://pkg.go.dev/github.com/thorlabsDev/ThorStreamer/sdks/go
- Main Repository: https://github.com/thorlabsDev/ThorStreamer