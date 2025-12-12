# Quick Start

Get streaming Solana data in under 5 minutes.

## Prerequisites

- Authentication token (contact [ThorLabs Discord](https://discord.gg/thorlabs))
- Server address

## Choose Your Language

{% tabs %}
{% tab title="Go" %}
### Installation

```bash
go get github.com/thorlabsDev/ThorStreamer/sdks/go@v0.1.0
```

### Stream Transactions

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
            break
        }
        if tx := msg.GetTransaction(); tx != nil {
            log.Printf("Transaction: slot=%d", tx.Transaction.Slot)
        }
    }
}
```
{% endtab %}

{% tab title="Rust" %}
### Installation

Add to `Cargo.toml`:

```toml
[dependencies]
thor-grpc-client = "0.1.0"
tokio = { version = "1", features = ["full"] }
```

### Stream Transactions

```rust
use thor_grpc_client::{ClientConfig, ThorClient, parse_message};
use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = ClientConfig {
        server_addr: std::env::var("SERVER_ADDRESS")?,
        token: std::env::var("AUTH_TOKEN")?,
        ..Default::default()
    };

    let mut client = ThorClient::new(config).await?;
    let mut stream = client.subscribe_to_transactions().await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        if let Some(EventMessage::Transaction(tx)) = msg.event_message {
            println!("Transaction: slot={}", tx.slot);
        }
    }
    Ok(())
}
```
{% endtab %}

{% tab title="TypeScript" %}
### Installation

```bash
npm install @grpc/grpc-js google-protobuf
```

### Stream Transactions

```typescript
import * as grpc from '@grpc/grpc-js';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { EventPublisherClient } from './proto/publisher_grpc_pb';

const client = new EventPublisherClient(
    process.env.SERVER_ADDRESS!,
    grpc.credentials.createInsecure()
);

const metadata = new grpc.Metadata();
metadata.set('authorization', process.env.AUTH_TOKEN!);

const stream = client.subscribeToTransactions(new Empty(), metadata);

stream.on('data', (response) => {
    console.log('Transaction received:', response.getData().length, 'bytes');
});

stream.on('error', (err) => {
    console.error('Stream error:', err);
});
```
{% endtab %}

{% tab title="Python" %}
### Installation

```bash
pip install grpcio grpcio-tools
```

### Stream Transactions

```python
import grpc
import asyncio
from google.protobuf.empty_pb2 import Empty

# Generate proto files first:
# python -m grpc_tools.protoc -I./proto --python_out=. --grpc_python_out=. proto/*.proto

from proto import publisher_pb2_grpc

async def main():
    metadata = [('authorization', 'your-token')]

    async with grpc.aio.insecure_channel('your-server:50051') as channel:
        stub = publisher_pb2_grpc.EventPublisherStub(channel)
        stream = stub.SubscribeToTransactions(Empty(), metadata=metadata)

        async for response in stream:
            print(f"Transaction: {len(response.data)} bytes")

asyncio.run(main())
```
{% endtab %}
{% endtabs %}

## Environment Setup

Create a `.env` file:

```bash
SERVER_ADDRESS=your-server:50051
AUTH_TOKEN=your-auth-token
```

## Next Steps

- [API Reference](api-reference/services.md) — Explore all available methods
- [SDK Documentation](sdks/) — Detailed guides for each language
- [Error Handling](error-handling.md) — Handle errors gracefully
