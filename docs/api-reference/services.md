# Services

ThorStreamer provides two gRPC services for streaming Solana blockchain data.

## EventPublisher Service

Recommended for most use cases. Provides individual streams with lower throughput requirements.

### SubscribeToTransactions

Stream real-time transactions matching configured program filters.

```protobuf
rpc SubscribeToTransactions(Empty) returns (stream StreamResponse)
```

**Input**: None
**Output**: Stream of transactions from [supported programs](program-filters.md)

**Use cases**:
- Monitor DeFi protocol interactions
- Track trading activity
- Market making strategies

---

### SubscribeToAccountUpdates

Stream account state changes with flexible filtering.

```protobuf
rpc SubscribeToAccountUpdates(SubscribeAccountsRequest) returns (stream StreamResponse)
```

**Input**:
```protobuf
message SubscribeAccountsRequest {
  repeated string account_address = 1;  // Specific accounts (max 100)
  repeated string owner_address = 2;    // Filter by owner program
}
```

**Filtering options**:
- Monitor specific accounts by address
- Monitor all accounts owned by specific programs
- Combine both filters

**Use cases**:
- Token balance tracking
- Program state monitoring
- Portfolio analytics

---

### SubscribeToSlotStatus

Monitor Solana slot progression and confirmations.

```protobuf
rpc SubscribeToSlotStatus(Empty) returns (stream StreamResponse)
```

**Input**: None
**Output**: Stream of slot status events

**Use cases**:
- Transaction confirmation tracking
- Block finality monitoring
- Network health checks

---

## ThorStreamer Service

High-throughput unified stream combining all event types. **Requires robust infrastructure.**

### StreamUpdates

Single stream delivering all events through `MessageWrapper`.

```protobuf
rpc StreamUpdates(Empty) returns (stream MessageWrapper)
```

**Output**: Stream containing one of:
- `TransactionEvent` — Transaction data
- `SubscribeUpdateAccountInfo` — Account state changes
- `SlotStatusEvent` — Slot progression

**Requirements**:
- 8+ CPU cores
- 16GB+ RAM
- Low-latency network
- Worker pools for parallel processing

**Warning**: Clients that cannot keep up will be automatically disconnected.
