# ThorStreamer Integration Guide

## Table of Contents
1. [Overview](#overview)
2. [Features](#features)
3. [Protocol Definition](#protocol-definition)
4. [Service Methods](#service-methods)
5. [Message Structures](#message-structures)
6. [Integration Guide](#integration-guide)
7. [Examples](#examples)
8. [Error Handling](#error-handling)
9. [Best Practices](#best-practices)

## Overview

ThorStreamer is a high-performance gRPC service providing real-time access to Solana blockchain data. It's specifically designed for trading bots and DeFi applications requiring low-latency access to on-chain data, offering filtered streams of transactions, account updates, and slot information. Our service filters transactions, focusing on selected accounts of interest, ensuring you only receive the most relevant data. This filtering enables faster response times and reduces overhead, making our streamer service one of the **fastest available for Solana event data stream**.

## Features

- Real-time transaction streaming with program filtering
- Account state change notifications
- Slot status monitoring
- Wallet-specific transaction tracking
- Automatic reconnection handling
- Protocol Buffer message encoding
- SSL/TLS encryption
- Token-based authentication
- Support for program-specific filters

## Protocol Definition
We provide `.proto` files that you can use to generate client libraries in the programming language of your choice. The `.proto` files define the structure of the data and the services that your client will interact with. You can use these files to generate client code for various languages using the respective gRPC tools.

To use the `.proto` files, follow the instructions provided in the gRPC documentation for your specific language to generate the necessary client libraries.

### Service Definition
```protobuf
service EventPublisher {
    // Stream all transactions matching configured program filters
    rpc SubscribeToTransactions(google.protobuf.Empty) returns (stream TransactionEvent) {}
    
    // Stream account data changes
    rpc SubscribeToAccountUpdates(google.protobuf.Empty) returns (stream UpdateAccountEvent) {}
    
    // Stream slot status updates
    rpc SubscribeToSlotStatus(google.protobuf.Empty) returns (stream SlotStatusEvent) {}
    
    // Stream transactions for specific wallets
    rpc SubscribeToWalletTransactions(SubscribeWalletRequest) returns (stream TransactionEvent) {}
}

message SubscribeWalletRequest {
    repeated string WalletAddress = 1;  // Array of Base58 encoded wallet addresses, max 10
}
```

## Service Methods

### 1. SubscribeToTransactions
- Purpose: Stream real-time transactions matching configured program filters
- Input: None (Empty message)
- Output: Stream of TransactionEvent messages
- Use cases:
  - Monitor specific program interactions
  - Track DeFi protocol transactions
  - Market making and trading strategies
    
| Program Name                | Program ID                                    |
|-----------------------------|-----------------------------------------------|
| Magic Eden V2               | `M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K` |
| Tensor Swap                 | `TSWAPaqyCSx2KABk68Shruf4rp7CxcNi8hAsbdwmHbN` |
| Raydium Liquidity Pool V4   | `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8`|
| Raydium Raydium CPMM        | `CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C`|
| OpenBook                    | `srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX`|
| Orca                        | `whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc` |
| Fluxbeam Program            | `FLUXubRmkEi2q6K3Y9kBPg9248ggaZVsoSFhtJHSrm1X`|
| Meteora Pools Program       | `Eo7WjKq67rjJQSZxS6z3YkapzY3eMj6Xy8X5EQVn5UaB`|
| Pump.fun                    | `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` |
| Pump.fun Fee Account        | `CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM` |
| Moonshot                    | `MoonCVVNZFSYkqNXP6bxHLPL6QQJiMagDL3qcqUQTrG` |
| game.com                    | `GameEs6zXFFGhE5zCdx2sqeRZkL7uYzPsZuSVn1fdxHF` |


*Note: Feel free to ask for any missing or requested program IDs.*


### 2. SubscribeToAccountUpdates (TO-DO)
- Purpose: Stream account data changes in real-time
*(In progress)*

### 3. SubscribeToSlotStatus
- Purpose: Monitor Solana slot progression
- Input: None (Empty message)
- Output: Stream of SlotStatusEvent messages
- Use cases:
  - Transaction confirmation tracking
  - Block finality monitoring
  - Network health checks

### 4. SubscribeToWalletTransactions
- Purpose: Stream transactions for specific wallet addresses
- Input: SubscribeWalletRequest with array of wallet addresses
- Output: Stream of TransactionEvent messages
- Limitations: Maximum 10 wallet addresses per subscription
- Use cases:
  - Portfolio tracking
  - User activity monitoring
  - Wallet analytics

## Message Structures

### Core Message Types

#### 1. TransactionEvent
```protobuf
message TransactionEvent {
    uint64 slot = 1;                           // Slot number when transaction was processed
    bytes signature = 2;                       // Transaction signature (64 bytes)
    uint64 index = 3;                         // Transaction index within the block
    bool is_vote = 4;                         // Indicates if this is a vote transaction
    SanitizedTransaction transaction = 5;      // Full transaction details
    TransactionStatusMeta transaction_status_meta = 6;  // Transaction status and metadata
}

message SanitizedTransaction {
    SanitizedMessage message = 1;              // Transaction message content
    bytes message_hash = 2;                    // Hash of the message (32 bytes)
    repeated bytes signatures = 3;             // Array of transaction signatures
    bool is_simple_vote_transaction = 4;       // Simple vote transaction indicator
}

message TransactionStatusMeta {
    bool is_status_err = 1;                    // Transaction success/failure flag
    uint64 fee = 2;                           // Transaction fee in lamports
    repeated uint64 pre_balances = 3;         // Account balances before transaction
    repeated uint64 post_balances = 4;        // Account balances after transaction
    repeated InnerInstructions inner_instructions = 5;  // CPI instruction details
    repeated string log_messages = 6;          // Program log messages
    repeated TransactionTokenBalance pre_token_balances = 7;   // Token balances before tx
    repeated TransactionTokenBalance post_token_balances = 8;  // Token balances after tx
    repeated Reward rewards = 9;               // Transaction rewards
    string error_info = 10;                    // Error details if transaction failed
}
```

#### 2. UpdateAccountEvent
```protobuf
under development
```

#### 3. SlotStatusEvent
```protobuf
message SlotStatusEvent {
    uint64 slot = 1;                          // Slot number
    uint64 parent = 2;                        // Parent slot number
    SlotStatus status = 3;                    // Current slot status
}

enum SlotStatus {
    PROCESSED = 0;    // Slot has been processed
    CONFIRMED = 1;    // Slot has been confirmed
    ROOTED = 2;       // Slot has been rooted (finalized)
}
```

### Supporting Types

#### Token Balance Information
```protobuf
message TransactionTokenBalance {
    uint32 account_index = 1;                 // Index of the token account
    string mint = 2;                          // Token mint address
    UiTokenAmount ui_token_account = 3;       // Token amount information
    string owner = 4;                         // Token account owner
}

message UiTokenAmount {
    string amount = 1;                        // Raw token amount
    uint32 decimals = 2;                      // Token decimals
    string ui_amount_string = 3;              // Human-readable amount
    google.protobuf.DoubleValue ui_amount = 4;  // Decimal amount
}
```

#### Instruction Details
```protobuf
message CompiledInstruction {
    uint32 program_id_index = 1;              // Index of program ID
    bytes data = 2;                           // Instruction data
    repeated uint32 accounts = 3;             // Account index array
}

message InnerInstructions {
    uint32 index = 1;                         // Instruction index
    repeated InnerInstruction instructions = 2;  // Inner instruction array
}

message InnerInstruction {
    CompiledInstruction instruction = 1;       // Instruction details
    optional uint32 stack_height = 2;          // CPI stack height
}
```

## Integration Guide

### Connection Details
- Server Address: `example-grpc.thornode.io:50051` (replace with the actual server address)
- Protocol: gRPC (HTTP/2)
- TLS: Enabled
  
### Authentication and Client Identification

### Token System
- Each authentication token serves as a unique client identifier
- Tokens must be included in all gRPC requests via metadata
- One token represents one client connection
```
metadata: {
    "authorization": "your_token"
}
```

## Client Limits

### Subscription Limits per Client (Token)
1. **Maximum Concurrent Subscriptions**: 4 total
2. **One Subscription per Event Type**:
   - 1x Transaction Stream (`SubscribeToTransactions`)
   - 1x Account Updates Stream (`SubscribeToAccountUpdates`)
   - 1x Slot Status Stream (`SubscribeToSlotStatus`)
   - 1x Wallet Transactions Stream (`SubscribeToWalletTransactions`)

### Subscription Rules
- Each client can maintain only one active subscription of each type
- New subscription attempts of the same type will be rejected
- Existing subscription must be closed before starting a new one
- Maximum 10 wallet addresses per wallet subscription

## Performance Requirements
- Clients MUST process messages quickly enough to keep up with the stream
- Slow clients will be automatically disconnected
- Warning messages will be sent before disconnection
- Disconnection triggers:
  - High message drop rate (>50% messages dropped)
  - Consistent processing latency spikes
  - Buffer overflow
  - Extended inactivity periods

## Implementation Examples

### Python Example
```python
import grpc
import asyncio
from google.protobuf.empty_pb2 import Empty
from thor_streamer.proto import publisher_pb2_grpc as pb2_grpc
from thor_streamer.proto import publisher_pb2 as pb2

class ThorClient:
    def __init__(self, server_addr: str, token: str):
        self.token = token
        self.server_addr = server_addr
        self.metadata = [('authorization', token)]

    async def process_transaction(self, tx):
        print(f"Processing tx: {tx.signature.hex()}")
        # Add your processing logic here

    async def subscribe_transactions(self):
        async with grpc.aio.insecure_channel(self.server_addr) as channel:
            stub = pb2_grpc.EventPublisherStub(channel)
            try:
                stream = stub.SubscribeToTransactions(
                    Empty(), metadata=self.metadata)
                async for tx in stream:
                    await self.process_transaction(tx)
            except grpc.RpcError as e:
                print(f"Stream error: {e.code()}")
                
client = ThorClient("example-grpc.thornode.io:50051", "your-token")
asyncio.run(client.subscribe_transactions())
```

### Rust Example
```rust
use tonic::{Request, Streaming};
use thor_proto::event_publisher_client::EventPublisherClient;
use thor_proto::{Empty, TransactionEvent};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = EventPublisherClient::connect("http://example-grpc.thornode.io:50051").await?;
    
    let mut request = Request::new(Empty {});
    request.metadata_mut().insert(
        "authorization",
        "your-token".parse()?
    );

    let stream: Streaming<TransactionEvent> = 
        client.subscribe_to_transactions(request).await?.into_inner();

    tokio::pin!(stream);
    while let Some(tx) = stream.try_next().await? {
        tokio::spawn(async move {
            process_transaction(tx).await;
        });
    }

    Ok(())
}

async fn process_transaction(tx: TransactionEvent) {
    println!("Processing tx: {:?}", hex::encode(&tx.signature));
    // Add your processing logic here
}
```

### TypeScript Example
```typescript
import * as grpc from '@grpc/grpc-js';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { EventPublisherClient } from './proto/publisher_grpc_pb';
import { TransactionEvent } from './proto/publisher_pb';

class ThorClient {
    private client: EventPublisherClient;
    private metadata: grpc.Metadata;

    constructor(serverAddr: string, token: string) {
        this.client = new EventPublisherClient(
            serverAddr,
            grpc.credentials.createInsecure()
        );
        this.metadata = new grpc.Metadata();
        this.metadata.set('authorization', token);
    }

    async subscribeTransactions() {
        const stream = this.client.subscribeToTransactions(
            new Empty(),
            this.metadata
        );

        return new Promise((resolve, reject) => {
            stream.on('data', (tx: TransactionEvent) => {
                this.processTransaction(tx);
            });

            stream.on('error', (err) => {
                console.error('Stream error:', err);
                reject(err);
            });

            stream.on('end', () => {
                resolve(undefined);
            });
        });
    }

    private async processTransaction(tx: TransactionEvent) {
        console.log('Processing tx:', tx.getSignature_asB64());
        // Add your processing logic here
    }
}

const client = new ThorClient('example-grpc.thornode.io:50051', 'your-token');
client.subscribeTransactions().catch(console.error);
```

### Go Example
```go
package main

import (
    "context"
    "log"
    "time"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "github.com/golang/protobuf/ptypes/empty"
    pb "thor_grpc/proto"
)

type ThorClient struct {
    client pb.EventPublisherClient
    token  string
}

func NewThorClient(serverAddr, token string) (*ThorClient, error) {
    conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    return &ThorClient{
        client: pb.NewEventPublisherClient(conn),
        token:  token,
    }, nil
}

func (c *ThorClient) SubscribeTransactions(ctx context.Context) error {
    md := metadata.New(map[string]string{
        "authorization": c.token,
    })
    ctx = metadata.NewOutgoingContext(ctx, md)

    stream, err := c.client.SubscribeToTransactions(ctx, &empty.Empty{})
    if err != nil {
        return err
    }

    for {
        tx, err := stream.Recv()
        if err != nil {
            return err
        }

        // Process asynchronously
        go func(tx *pb.TransactionEvent) {
            if err := c.processTransaction(tx); err != nil {
                log.Printf("Error processing tx: %v", err)
            }
        }(tx)
    }
}

func (c *ThorClient) processTransaction(tx *pb.TransactionEvent) error {
    log.Printf("Processing tx: %x", tx.Signature)
    // Add your processing logic here
    return nil
}

func main() {
    client, err := NewThorClient("example-grpc.thornode.io:50051", "your-token")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    if err := client.SubscribeTransactions(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Error Codes Overview

### Authentication Errors
- `UNAUTHENTICATED`: Token is missing or invalid
- `TOKEN_EXPIRED`: Authentication token has expired
- `INVALID_TOKEN`: Token format is incorrect

### Subscription Limits
- `SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum number of allowed subscriptions (4)
- `DUPLICATE_SUBSCRIPTION_TYPE`: Client already has an active subscription of this type
- `TOO_MANY_WALLET_ADDRESSES`: Wallet subscription request exceeds maximum allowed addresses (10)

### Connection Errors
- `CONNECTION_ERROR`: Failed to establish connection
- `CONNECTION_CLOSED`: Connection was closed by the server
- `CONNECTION_TIMEOUT`: Connection timed out
- `STREAM_CLOSED`: Stream was closed by the server

### Request Errors
- `INVALID_REQUEST`: Malformed request
- `INVALID_PROGRAM_ID`: Program ID format is incorrect
- `INVALID_WALLET_ADDRESS`: Wallet address format is incorrect
- `EMPTY_WALLET_LIST`: No wallet addresses provided in subscription request

### Server Errors
- `INTERNAL`: Internal server error
- `UNAVAILABLE`: Service temporarily unavailable
- `RESOURCE_EXHAUSTED`: Server resources are exhausted
- `UNKNOWN`: Unknown server error

### Data Errors
- `INVALID_MESSAGE_FORMAT`: Received message format is invalid
- `DESERIALIZATION_ERROR`: Unable to deserialize message
- `DATA_CORRUPTION`: Received corrupted data

## Error Handling

### Critical Errors (Require Immediate Action)
- `UNAUTHENTICATED`
- `TOKEN_EXPIRED`
- `CONNECTION_ERROR`
- `INTERNAL`
- `RESOURCE_EXHAUSTED`

### Recoverable Errors (Can Retry)
- `CONNECTION_TIMEOUT`
- `UNAVAILABLE`
- `STREAM_CLOSED`
- `CONNECTION_CLOSED`

### Validation Errors (Check Input)
- `INVALID_REQUEST`
- `INVALID_PROGRAM_ID`
- `INVALID_WALLET_ADDRESS`
- `TOO_MANY_WALLET_ADDRESSES`
- `EMPTY_WALLET_LIST`

### Data Handling Errors
- `INVALID_MESSAGE_FORMAT`
- `DESERIALIZATION_ERROR`
- `DATA_CORRUPTION`

### Error Handling Recommendations

#### Critical Errors
- Require manual intervention or service restart
- Should trigger alerts/notifications
- Need immediate attention from developers/operators

#### Recoverable Errors
- Implement automatic retry with backoff
- Maintain subscription state for recovery
- Log errors for monitoring

#### Validation Errors
- Check input before making requests
- Provide clear error messages to users
- Update client configuration as needed

#### Data Errors
- Log error details for debugging
- Monitor frequency of occurrences
- Report to service provider if persistent

### Monitoring Recommendations

- Track error frequencies by type
- Set up alerts for critical errors
- Monitor error patterns
- Track recovery success rates
- Log detailed error contexts

### Additional Considerations

1. **Error Persistence**
   - Some errors may be persistent until token refresh
   - Some may require service restart
   - Others may resolve automatically

2. **Rate Limiting**
   - Track error rates
   - Implement circuit breakers where appropriate
   - Use exponential backoff for retries

3. **Error Reporting**
   - Include error context in logs
   - Track error frequencies
   - Monitor error patterns

## Best Practices

1. Always use the latest version of the generated client code for your language, following updates to the `.proto` files.
2. Implement reconnection logic to handle disconnects gracefully in your application.
3. Use SSL/TLS for secure communication (gRPC supports this natively in most languages).
4. Handle backpressure by implementing flow control if your application processes events slower than it receives them.
5. Consider implementing a worker pool for parallel processing of events, especially if the event processing is CPU-intensive.

## Disclaimer

This service and associated `.proto` files are provided by **ThorLabs** for use in interacting with real-time Solana blockchain data streams. While we strive to ensure that the service and the generated code are reliable and accurate, **ThorLabs** provides no warranties or guarantees regarding the availability, accuracy, completeness, or suitability of the provided data for any particular purpose.

### Important Notes:
- **Use at Your Own Risk**: The use of this service, including the integration of the provided `.proto` files and the client-side implementation, is at your own risk. ThorLabs will not be held responsible for any direct, indirect, or consequential damages resulting from the use of the service.
- **Data Accuracy**: While the service is designed to deliver accurate, filtered event data, there may be cases where events are delayed, missing, or out-of-sync. Always consider additional mechanisms to verify critical transaction data.
- **Compliance**: By using this service, you agree to comply with all applicable laws and regulations in your jurisdiction. It is the responsibility of the user to ensure that the use of the provided data does not violate any applicable regulations or third-party rights.
- **Modification**: ThorLabs reserves the right to modify, suspend, or discontinue the service at any time without notice.

This documentation provides comprehensive technical information for integrating with ThorStreamer. For additional support or questions, please contact our support team via our [discord server](https://discord.gg/thorlabs).

