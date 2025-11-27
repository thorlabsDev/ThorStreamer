# Overview

## Architecture

ThorStreamer uses gRPC (HTTP/2) with Protocol Buffers for efficient, low-latency streaming. The service provides two interfaces:

### ThorStreamer Service (Unified Stream)
A single high-throughput stream delivering all event types through one connection.

```protobuf
service ThorStreamer {
  rpc StreamUpdates(Empty) returns (stream MessageWrapper);
}
```

**Use when**: You need all events and have high-performance infrastructure (8+ cores, 16GB+ RAM).

### EventPublisher Service (Individual Streams)
Separate subscription methods for each event type with filtering options.

```protobuf
service EventPublisher {
  rpc SubscribeToTransactions(Empty) returns (stream StreamResponse);
  rpc SubscribeToSlotStatus(Empty) returns (stream StreamResponse);
  rpc SubscribeToWalletTransactions(SubscribeWalletRequest) returns (stream StreamResponse);
  rpc SubscribeToAccountUpdates(SubscribeAccountsRequest) returns (stream StreamResponse);
}
```

**Use when**: You need specific event types or have modest hardware (2-4 cores, 2-4GB RAM).

## Authentication

All requests require a token via gRPC metadata:

```
metadata: {
  "authorization": "your_token"
}
```

Each token identifies a unique client and enforces subscription limits.

## Connection Details

| Property | Value |
|----------|-------|
| Protocol | gRPC (HTTP/2) |
| Encryption | TLS enabled |
| Port | 50051 |

## Data Flow

```
┌─────────────┐     gRPC Stream      ┌──────────────┐
│   Client    │ ◄──────────────────► │ ThorStreamer │
│   (SDK)     │   MessageWrapper     │   Service    │
└─────────────┘   or StreamResponse  └──────────────┘
```

Events are delivered as they occur on-chain. The unified stream wraps all events in `MessageWrapper` with type identification. Individual streams return `StreamResponse` containing serialized event data.
