# Go SDK

Official Go client library for ThorStreamer.

## Installation

```bash
go get github.com/thorlabsDev/ThorStreamer/sdks/go@v0.1.0
```

## Requirements

- Go 1.23.2+

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

## API Reference

### Creating a Client

```go
client, err := thorclient.NewClient(thorclient.Config{
    ServerAddr:     "server:50051",
    Token:          "your-token",
    DefaultTimeout: 30 * time.Second,
})
defer client.Close()
```

### SubscribeToTransactions

```go
stream, err := client.SubscribeToTransactions(ctx)

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    if tx := msg.GetTransaction(); tx != nil {
        log.Printf("Signature: %x", tx.Transaction.Signature)
    }
}
```

### SubscribeToSlotStatus

```go
stream, err := client.SubscribeToSlotStatus(ctx)

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

### SubscribeToWalletTransactions

Monitor up to 10 wallet addresses:

```go
wallets := []string{
    "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM",
    "2ojv9BAiHUrvsm9gxDe7fJSzbNZSJcxZvf8dqmWGHG8S",
}
stream, err := client.SubscribeToWalletTransactions(ctx, wallets)

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    if tx := msg.GetTransaction(); tx != nil {
        log.Printf("Wallet tx: slot=%d", tx.Transaction.Slot)
    }
}
```

### SubscribeToAccountUpdates

Monitor accounts with optional owner filtering:

```go
accounts := []string{"account1...", "account2..."}
owners := []string{"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"}

stream, err := client.SubscribeToAccountUpdates(ctx, accounts, owners)

for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    if acc := msg.GetAccountUpdate(); acc != nil {
        log.Printf("Account: %x, lamports=%d", acc.Pubkey, acc.Lamports)
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
    // Handle other errors
    log.Printf("Stream error: %v", err)
    return
}
```

## Full Example

See [examples/golang-advanced](https://github.com/thorlabsDev/ThorStreamer/tree/master/examples/golang-advanced) for a complete implementation with filtering, logging, and event handling.

## Resources

- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/thorlabsDev/ThorStreamer/sdks/go)
- [GitHub Repository](https://github.com/thorlabsDev/ThorStreamer)
