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
- Unified streaming interface with message wrapping

## Protocol Definition

We provide `.proto` files that you can use to generate client libraries in the programming language of your choice. The `.proto` files define the structure of the data and the services that your client will interact with. You can use these files to generate client code for various languages using the respective gRPC tools.

To use the `.proto` files, follow the instructions provided in the gRPC documentation for your specific language to generate the necessary client libraries.

### Service Definition

The ThorStreamer service provides two main interfaces:

#### ThorStreamer Service (events.proto)
```protobuf
service ThorStreamer {
  rpc StreamUpdates(Empty) returns (stream MessageWrapper);
}
```

#### EventPublisher Service (publisher.proto)
```protobuf
service EventPublisher {
  rpc SubscribeToTransactions(google.protobuf.Empty) returns (stream StreamResponse) {}
  rpc SubscribeToSlotStatus(google.protobuf.Empty) returns (stream StreamResponse) {}
  rpc SubscribeToWalletTransactions(SubscribeWalletRequest) returns (stream StreamResponse) {}
  rpc SubscribeToAccountUpdates(SubscribeAccountsRequest) returns (stream StreamResponse) {}
}

message SubscribeWalletRequest {
  repeated string wallet_address = 1; // Array of Base58 encoded wallet addresses, max 10
}

message SubscribeAccountsRequest {
  repeated string account_address = 1;
}
```

## Service Methods

### ThorStreamer Service

#### StreamUpdates
- **Purpose**: Unified stream for all event types (transactions, account updates, slot status)
- **Input**: Empty message
- **Output**: Stream of MessageWrapper messages containing different event types
- **Performance Requirements**: High-throughput stream requiring powerful client infrastructure
- **Use cases**: Single connection for all data types with message type identification
- **⚠️ Warning**: This is an extremely high-volume stream that requires robust server infrastructure to process effectively

### EventPublisher Service

#### 1. SubscribeToTransactions
- **Purpose**: Stream real-time transactions matching configured program filters
- **Input**: None (Empty message)
- **Output**: Stream of StreamResponse messages containing transaction data
- **Use cases**:
    - Monitor specific program interactions
    - Track DeFi protocol transactions
    - Market making and trading strategies


| Program Name | Program ID |
|---------|---|
| Fluxbeam | `FLUXubRmkEi2q6K3Y9kBPg9248ggaZVsoSFhtJHSrm1X` |
| game.com | `GameEs6zXFFGhE5zCdx2sqeRZkL7uYzPsZuSVn1fdxHF` |
| Lifinity V2 | `2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c` |
| Magic Eden V2 | `M2mx93ekt1fmXSVkTrUL9xVFHkmME8HTUi5Cyc5aF7K` |
| Meteora DAMM v2 | `cpamdpZCGKUy5JxQXB4dcpGPiikHawvSWAd6mEn1sGG` |
| Meteora DLMM | `LBUZKhRxPF3XUpBCjp4YzTKgLccjZhTSDM9YuVaPwxo` |
| Meteora Pools | `Eo7WjKq67rjJQSZxS6z3YkapzY3eMj6Xy8X5EQVn5UaB` |
| Moonshot | `MoonCVVNZFSYkqNXP6bxHLPL6QQJiMagDL3qcqUQTrG` |
| OpenBook | `srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX` |
| Orca | `whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc` |
| Pump.fun | `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` |
| Pump.fun AMM | `pAMMBay6oceH9fJKBRHGP5D4bD4sWpmSwMn52FMfXEA` |
| Pump.fun Fee Account | `CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM` |
| Pump.fun: Raydium Migration | `39azUYFWPz3VHgKCf3VChUwbpURdCHRxjWVowf5jUJjg` |
| Raydium CLMM | `CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK` |
| Raydium CPMM | `CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C` |
| Raydium Launchpad | `LanMV9sAd7wArD4vJFi2qDdfnVhFxYSUg6eADduJ3uj` |
| Raydium Launchpad Authority | `WLHv2UAZm6z4KyaaELi5pjdbJh6RESMva1Rnn8pJVVh` |
| Raydium Liquidity Pool V4 | `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8` |
| SolFi | `SoLFiHG9TfgtdUXUjWAxi3LtvYuFyDLVhBWxdMZxyCe` |
| Tensor cNFT | `TCMPhJdwDryooaGtiocG1u3xcYbRpiJzb283XfCZsDp` |
| Tensor Swap | `TSWAPaqyCSx2KABk68Shruf4rp7CxcNi8hAsbdwmHbN` |

*Note: Feel free to ask for any missing or requested program IDs.*

#### 2. SubscribeToAccountUpdates
- **Purpose**: Stream account data changes in real-time
- **Input**: SubscribeAccountsRequest with array of account addresses
- **Output**: Stream of StreamResponse messages containing account update data
- **Use cases**:
    - Monitor token account balances
    - Track program account state changes
    - Real-time portfolio tracking

#### 3. SubscribeToSlotStatus
- **Purpose**: Monitor Solana slot progression
- **Input**: None (Empty message)
- **Output**: Stream of StreamResponse messages containing slot status data
- **Use cases**:
    - Transaction confirmation tracking
    - Block finality monitoring
    - Network health checks

#### 4. SubscribeToWalletTransactions
- **Purpose**: Stream transactions for specific wallet addresses
- **Input**: SubscribeWalletRequest with array of wallet addresses
- **Output**: Stream of StreamResponse messages containing transaction data
- **Limitations**: Maximum 10 wallet addresses per subscription
- **Use cases**:
    - Portfolio tracking
    - User activity monitoring
    - Wallet analytics

## Message Structures

### Core Message Types

#### 1. MessageWrapper (Unified Stream)
```protobuf
message MessageWrapper {
  oneof event_message {
    SubscribeUpdateAccountInfo account_update = 1;
    SlotStatusEvent slot = 2;
    TransactionEventWrapper transaction = 3;
  }
}

message TransactionEventWrapper {
  StreamType stream_type = 1;
  TransactionEvent transaction = 2;
}

enum StreamType {
  STREAM_TYPE_UNSPECIFIED = 0;
  STREAM_TYPE_FILTERED = 1;
  STREAM_TYPE_WALLET = 2;
  STREAM_TYPE_ACCOUNT = 3;
}
```

#### 2. TransactionEvent
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
  Message message = 1;                       // Transaction message content
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

#### 3. SubscribeUpdateAccountInfo
```protobuf
message SubscribeUpdateAccountInfo {
  bytes pubkey = 1;                         // Account public key
  uint64 lamports = 2;                      // Account balance in lamports
  bytes owner = 3;                          // Account owner program
  bool executable = 4;                      // Whether account is executable
  uint64 rent_epoch = 5;                    // Rent epoch
  bytes data = 6;                           // Account data
  uint64 write_version = 7;                 // Write version
  optional bytes txn_signature = 8;         // Transaction signature that caused update
  optional SlotStatus slot = 9;             // Slot information
}
```

#### 4. SlotStatusEvent
```protobuf
message SlotStatusEvent {
  uint64 slot = 1;                          // Slot number
  uint64 parent = 2;                        // Parent slot number
  int32 status = 3;                         // Slot status (processed/confirmed/rooted)
  bytes block_hash = 4;                     // Block hash (32 bytes)
  uint64 block_height = 5;                  // Block height
}
```

### Supporting Types

#### Message Structure (Optimized for Legacy and v0 transactions)
```protobuf
message Message {
  uint32 version = 1;                       // 0 for legacy, 1 for v0
  MessageHeader header = 2;
  bytes recent_block_hash = 3;              // 32 bytes
  repeated bytes account_keys = 4;          // Array of 32 byte keys
  repeated CompiledInstruction instructions = 5;
  repeated MessageAddressTableLookup address_table_lookups = 6;  // Only used for v0
  LoadedAddresses loaded_addresses = 7;                          // Only used for v0
  repeated bool is_writable = 8;                                // Account write permissions
}

message MessageHeader {
  uint32 num_required_signatures = 1;
  uint32 num_readonly_signed_accounts = 2;
  uint32 num_readonly_unsigned_accounts = 3;
}
```

#### Token Balance Information
```protobuf
message TransactionTokenBalance {
  uint32 account_index = 1;                 // Index of the token account
  string mint = 2;                          // Token mint address
  UiTokenAmount ui_token_amount = 3;        // Token amount information
  string owner = 4;                         // Token account owner
}

message UiTokenAmount {
  double ui_amount = 1;                     // Decimal amount
  uint32 decimals = 2;                      // Token decimals
  string amount = 3;                        // Raw token amount string
  string ui_amount_string = 4;              // Human-readable amount string
}
```

#### Instruction Details
```protobuf
message CompiledInstruction {
  uint32 program_id_index = 1;              // Index of program ID
  bytes data = 2;                           // Instruction data
  repeated uint32 accounts = 3;             // Account index array (packed)
}

message InnerInstructions {
  uint32 index = 1;                         // Instruction index
  repeated InnerInstruction instructions = 2;  // Inner instruction array
}

message InnerInstruction {
  CompiledInstruction instruction = 1;       // Instruction details
  optional uint32 stack_height = 2;         // CPI stack height
}
```

## Integration Guide

### Connection Details
- Server Address: `example-grpc.thornode.io:50051` (replace with the actual server address)
- Protocol: gRPC (HTTP/2)
- TLS: Enabled

### Authentication and Client Identification

#### Token System
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
1. **Maximum Concurrent Subscriptions**: 6 total
2. **Multiple Subscriptions per Event Type**:
    - Up to 2x Transaction Streams (`SubscribeToTransactions`)
    - Up to 5x Account Updates Streams (`SubscribeToAccountUpdates`)
    - Up to 2x Slot Status Streams (`SubscribeToSlotStatus`)
    - Up to 10x Wallet Transactions Streams (`SubscribeToWalletTransactions`)

### Subscription Rules
- Each client can maintain multiple active subscriptions of each type within limits
- Maximum 10 wallet addresses per wallet subscription
- Maximum 100 account addresses per account subscription
- Total subscriptions across all types cannot exceed 6 per client

## Performance Requirements

### Client Infrastructure Requirements

#### For Unified Stream (StreamUpdates)
The unified stream is an **extremely high-throughput data feed** that combines all event types. Using this stream requires:

- **High-Performance Servers**: Multi-core processors (8+ cores recommended)
- **Sufficient RAM**: Minimum 8GB, 16GB+ recommended for production
- **Fast Network Connection**: Low-latency, high-bandwidth connection
- **Optimized Code**: Efficient message processing and deserialization
- **Proper Architecture**: Asynchronous processing, worker pools, message queuing

#### For Individual Streams (EventPublisher)
Individual streams have more manageable throughput and can run on modest hardware:
- **Standard Servers**: 2-4 core processors sufficient for most use cases
- **Moderate RAM**: 2-4GB typically adequate
- **Regular Network**: Standard broadband connections work well

### Disconnection Policy

#### Automatic Disconnection Triggers
Clients will be **automatically disconnected** if they cannot keep up with the stream:

- **High Message Drop Rate**: >50% of messages dropped consistently
- **Processing Latency Spikes**: Consistent delays in message acknowledgment
- **Buffer Overflow**: Client-side or server-side buffer saturation
- **Extended Inactivity**: No message processing activity detected

#### Performance Monitoring
- **Warning Messages**: Sent before disconnection to allow optimization
- **Performance Metrics**: Server monitors client processing capabilities
- **Automatic Cleanup**: Slow clients are removed to maintain service quality

#### Client Responsibility
**If you experience frequent disconnections:**
1. **Upgrade Your Infrastructure**: More powerful servers, better network
2. **Optimize Your Code**: Implement async processing, reduce per-message latency
3. **Consider Individual Streams**: Use EventPublisher for lower-volume needs
4. **Monitor Resource Usage**: CPU, memory, and network utilization

**⚠️ Important**: Frequent disconnections indicate your client infrastructure cannot handle the unified stream's throughput. This is not a service issue - you need more powerful hardware and optimized code.

## Implementation Examples

### Python Example (Unified Stream)
```python
import grpc
import asyncio
from google.protobuf.empty_pb2 import Empty
from thor_streamer.proto import events_pb2_grpc as events_grpc
from thor_streamer.proto import events_pb2 as events_pb2

class ThorClient:
    def __init__(self, server_addr: str, token: str):
        self.token = token
        self.server_addr = server_addr
        self.metadata = [('authorization', token)]

    async def process_message(self, wrapper):
        if wrapper.HasField('transaction'):
            tx_wrapper = wrapper.transaction
            tx = tx_wrapper.transaction
            print(f"Processing tx: {tx.signature.hex()}, type: {tx_wrapper.stream_type}")
        elif wrapper.HasField('account_update'):
            account = wrapper.account_update
            print(f"Account update: {account.pubkey.hex()}")
        elif wrapper.HasField('slot'):
            slot = wrapper.slot
            print(f"Slot update: {slot.slot}, status: {slot.status}")

    async def stream_updates(self):
        async with grpc.aio.insecure_channel(self.server_addr) as channel:
            stub = events_grpc.ThorStreamerStub(channel)
            try:
                stream = stub.StreamUpdates(Empty(), metadata=self.metadata)
                async for message in stream:
                    await self.process_message(message)
            except grpc.RpcError as e:
                print(f"Stream error: {e.code()}")
                
client = ThorClient("example-grpc.thornode.io:50051", "your-token")
asyncio.run(client.stream_updates())
```

### Python Example (EventPublisher - Account Updates)
```python
import grpc
import asyncio
from thor_streamer.proto import publisher_pb2_grpc as pb2_grpc
from thor_streamer.proto import publisher_pb2 as pb2

class ThorAccountClient:
    def __init__(self, server_addr: str, token: str):
        self.token = token
        self.server_addr = server_addr
        self.metadata = [('authorization', token)]

    async def process_account_update(self, data):
        # Deserialize the StreamResponse data field
        # Implementation depends on how the data is encoded
        print(f"Account update received: {len(data)} bytes")

    async def subscribe_account_updates(self, account_addresses):
        async with grpc.aio.insecure_channel(self.server_addr) as channel:
            stub = pb2_grpc.EventPublisherStub(channel)
            request = pb2.SubscribeAccountsRequest(account_address=account_addresses)
            
            try:
                stream = stub.SubscribeToAccountUpdates(request, metadata=self.metadata)
                async for response in stream:
                    await self.process_account_update(response.data)
            except grpc.RpcError as e:
                print(f"Stream error: {e.code()}")

# Subscribe to specific accounts
accounts = ["11111111111111111111111111111112", "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"]
client = ThorAccountClient("example-grpc.thornode.io:50051", "your-token")
asyncio.run(client.subscribe_account_updates(accounts))
```

### Rust Example (Unified Stream)
```rust
use tonic::{Request, Streaming};
use thor_proto::thor_streamer_client::ThorStreamerClient;
use thor_proto::{Empty, MessageWrapper};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = ThorStreamerClient::connect("http://example-grpc.thornode.io:50051").await?;
    
    let mut request = Request::new(Empty {});
    request.metadata_mut().insert(
        "authorization",
        "your-token".parse()?
    );

    let stream: Streaming<MessageWrapper> = 
        client.stream_updates(request).await?.into_inner();

    tokio::pin!(stream);
    while let Some(wrapper) = stream.try_next().await? {
        tokio::spawn(async move {
            process_message(wrapper).await;
        });
    }

    Ok(())
}

async fn process_message(wrapper: MessageWrapper) {
    match wrapper.event_message {
        Some(thor_proto::message_wrapper::EventMessage::Transaction(tx_wrapper)) => {
            let tx = tx_wrapper.transaction.unwrap();
            println!("Processing tx: {:?}, type: {:?}", 
                hex::encode(&tx.signature), tx_wrapper.stream_type);
        },
        Some(thor_proto::message_wrapper::EventMessage::AccountUpdate(account)) => {
            println!("Account update: {:?}", hex::encode(&account.pubkey));
        },
        Some(thor_proto::message_wrapper::EventMessage::Slot(slot)) => {
            println!("Slot update: {}, status: {}", slot.slot, slot.status);
        },
        None => {}
    }
}
```

### TypeScript Example (EventPublisher)
```typescript
import * as grpc from '@grpc/grpc-js';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { EventPublisherClient } from './proto/publisher_grpc_pb';
import { SubscribeAccountsRequest, StreamResponse } from './proto/publisher_pb';

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

    async subscribeAccountUpdates(accountAddresses: string[]) {
        const request = new SubscribeAccountsRequest();
        request.setAccountAddressList(accountAddresses);

        const stream = this.client.subscribeToAccountUpdates(request, this.metadata);

        return new Promise((resolve, reject) => {
            stream.on('data', (response: StreamResponse) => {
                this.processAccountUpdate(response.getData_asU8());
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

    private async processAccountUpdate(data: Uint8Array) {
        console.log('Account update received:', data.length, 'bytes');
        // Process the account update data
    }
}

const client = new ThorClient('example-grpc.thornode.io:50051', 'your-token');
const accounts = ['11111111111111111111111111111112', 'TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA'];
client.subscribeAccountUpdates(accounts).catch(console.error);
```

### Go Example (EventPublisher - Wallet Transactions)
```go
package main

import (
    "context"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
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

func (c *ThorClient) SubscribeWalletTransactions(ctx context.Context, wallets []string) error {
    md := metadata.New(map[string]string{
        "authorization": c.token,
    })
    ctx = metadata.NewOutgoingContext(ctx, md)

    request := &pb.SubscribeWalletRequest{
        WalletAddress: wallets,
    }

    stream, err := c.client.SubscribeToWalletTransactions(ctx, request)
    if err != nil {
        return err
    }

    for {
        response, err := stream.Recv()
        if err != nil {
            return err
        }

        // Process asynchronously
        go func(data []byte) {
            if err := c.processWalletTransaction(data); err != nil {
                log.Printf("Error processing wallet tx: %v", err)
            }
        }(response.Data)
    }
}

func (c *ThorClient) processWalletTransaction(data []byte) error {
    log.Printf("Processing wallet transaction: %d bytes", len(data))
    // Deserialize and process the transaction data
    return nil
}

func main() {
    client, err := NewThorClient("example-grpc.thornode.io:50051", "your-token")
    if err != nil {
        log.Fatal(err)
    }

    wallets := []string{
        "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM",
        "2ojv9BAiHUrvsm9gxDe7fJSzbNZSJcxZvf8dqmWGHG8S",
    }

    ctx := context.Background()
    if err := client.SubscribeWalletTransactions(ctx, wallets); err != nil {
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
- `SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum number of allowed subscriptions (6 total)
- `TRANSACTION_SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum transaction subscriptions (2)
- `ACCOUNT_SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum account subscriptions (5)
- `SLOT_SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum slot subscriptions (2)
- `WALLET_SUBSCRIPTION_LIMIT_REACHED`: Client has reached maximum wallet subscriptions (10)
- `TOO_MANY_WALLET_ADDRESSES`: Wallet subscription request exceeds maximum allowed addresses (10)
- `TOO_MANY_ACCOUNT_ADDRESSES`: Account subscription request exceeds maximum allowed addresses (100)

### Connection Errors
- `CONNECTION_ERROR`: Failed to establish connection
- `CONNECTION_CLOSED`: Connection was closed by the server
- `CONNECTION_TIMEOUT`: Connection timed out
- `STREAM_CLOSED`: Stream was closed by the server

### Request Errors
- `INVALID_REQUEST`: Malformed request
- `INVALID_PROGRAM_ID`: Program ID format is incorrect
- `INVALID_WALLET_ADDRESS`: Wallet address format is incorrect
- `INVALID_ACCOUNT_ADDRESS`: Account address format is incorrect
- `EMPTY_WALLET_LIST`: No wallet addresses provided in subscription request
- `EMPTY_ACCOUNT_LIST`: No account addresses provided in subscription request

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
- `INVALID_ACCOUNT_ADDRESS`
- `TRANSACTION_SUBSCRIPTION_LIMIT_REACHED`
- `ACCOUNT_SUBSCRIPTION_LIMIT_REACHED`
- `SLOT_SUBSCRIPTION_LIMIT_REACHED`
- `WALLET_SUBSCRIPTION_LIMIT_REACHED`
- `TOO_MANY_WALLET_ADDRESSES`
- `TOO_MANY_ACCOUNT_ADDRESSES`
- `EMPTY_WALLET_LIST`
- `EMPTY_ACCOUNT_LIST`

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

1. **Protocol Buffer Usage**
    - Always use the latest version of the generated client code for your language, following updates to the `.proto` files
    - Handle both message formats (unified MessageWrapper and individual StreamResponse)
    - Use proper deserialization for StreamResponse data fields

2. **Connection Management**
    - Implement reconnection logic to handle disconnects gracefully in your application
    - Use SSL/TLS for secure communication (gRPC supports this natively in most languages)
    - **Choose the right stream type**: Use unified stream (`StreamUpdates`) only if you have powerful infrastructure; otherwise use individual EventPublisher streams
    - **Monitor disconnections**: Frequent disconnections indicate insufficient client resources

3. **Performance Optimization**
    - Handle backpressure by implementing flow control if your application processes events slower than it receives them
    - **Mandatory for unified stream**: Implement worker pools for parallel processing - single-threaded processing will not work
    - Use asynchronous processing patterns to handle high-volume streams
    - **Hardware requirements**: Ensure adequate CPU, RAM, and network resources before using unified stream
    - Use the packed encoding benefits for repeated numeric fields

4. **Message Processing**
    - Handle oneof fields properly in MessageWrapper
    - Process different stream types (FILTERED, WALLET, ACCOUNT) appropriately
    - Implement proper error handling for message deserialization

5. **Resource Management**
    - Respect subscription limits (6 total: up to 2 transaction, 5 account, 2 slot, 10 wallet)
    - Monitor your subscription usage across different event types
    - Close unused subscriptions to free resources for new ones
    - Monitor memory usage when processing high-volume streams
    - Consider subscription strategy based on your application needs

## Disclaimer

This service and associated `.proto` files are provided by **ThorLabs** for use in interacting with real-time Solana blockchain data streams. While we strive to ensure that the service and the generated code are reliable and accurate, **ThorLabs** provides no warranties or guarantees regarding the availability, accuracy, completeness, or suitability of the provided data for any particular purpose.

### Important Notes:
- **Use at Your Own Risk**: The use of this service, including the integration of the provided `.proto` files and the client-side implementation, is at your own risk. ThorLabs will not be held responsible for any direct, indirect, or consequential damages resulting from the use of the service.
- **Data Accuracy**: While the service is designed to deliver accurate, filtered event data, there may be cases where events are delayed, missing, or out-of-sync. Always consider additional mechanisms to verify critical transaction data.
- **Compliance**: By using this service, you agree to comply with all applicable laws and regulations in your jurisdiction. It is the responsibility of the user to ensure that the use of the provided data does not violate any applicable regulations or third-party rights.
- **Modification**: ThorLabs reserves the right to modify, suspend, or discontinue the service at any time without notice.

This documentation provides comprehensive technical information for integrating with ThorStreamer. For additional support or questions, please contact our support team via our [discord server](https://discord.gg/thorlabs).