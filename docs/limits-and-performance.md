# Limits & Performance

## Subscription Limits

Each authentication token (client) is limited to:

| Limit | Value |
|-------|-------|
| **Total concurrent subscriptions** | 6 |
| Transaction streams | 2 max |
| Account update streams | 5 max |
| Slot status streams | 2 max |

### Per-Subscription Limits

| Subscription Type | Limit |
|-------------------|-------|
| Account addresses per request | 100 |

---

## Performance Requirements

### EventPublisher Service (Individual Streams)

Recommended for most use cases.

| Resource | Requirement |
|----------|-------------|
| CPU | 2-4 cores |
| RAM | 2-4 GB |
| Network | Standard broadband |

### ThorStreamer Service (Unified Stream)

High-throughput stream combining all events. **Requires robust infrastructure.**

| Resource | Requirement |
|----------|-------------|
| CPU | 8+ cores |
| RAM | 16+ GB |
| Network | Low-latency, high-bandwidth |

**Mandatory architecture**:
- Worker pools for parallel message processing
- Asynchronous event handling
- Message queuing for backpressure management

---

## Disconnection Policy

Clients are automatically disconnected when they cannot keep up with the stream:

- **High message drop rate** — >50% messages dropped consistently
- **Processing latency** — Sustained delays in message handling
- **Buffer overflow** — Client or server-side buffer saturation

### Warning Signs

Before disconnection, you may observe:
- Increasing latency in message delivery
- Gaps in slot numbers or transaction indices
- Server-side warning messages

### Mitigation

If experiencing disconnections:

1. **Upgrade infrastructure** — More CPU, RAM, better network
2. **Optimize code** — Implement async processing, reduce per-message latency
3. **Use individual streams** — EventPublisher streams have lower throughput requirements
4. **Monitor resources** — Track CPU, memory, and network utilization

---

## Choosing the Right Service

| Use Case | Recommended Service |
|----------|---------------------|
| Monitor specific programs | EventPublisher: `SubscribeToTransactions` |
| Monitor account states | EventPublisher: `SubscribeToAccountUpdates` |
| Block confirmations | EventPublisher: `SubscribeToSlotStatus` |
| All data, high-performance infra | ThorStreamer: `StreamUpdates` |

---

## Optimization Tips

### Message Processing

```go
// Use goroutines for parallel processing (Go)
for {
    msg, err := stream.Recv()
    if err != nil {
        break
    }
    go processMessage(msg)  // Process concurrently
}
```

```rust
// Use tokio::spawn for async processing (Rust)
while let Some(response) = stream.message().await? {
    let data = response.data.clone();
    tokio::spawn(async move {
        process_message(data).await;
    });
}
```

### Buffer Management

- Use bounded channels/queues to prevent memory exhaustion
- Implement circuit breakers for downstream systems
- Monitor queue depths as a health indicator

### Network

- Deploy clients in same region as ThorStreamer server
- Use persistent connections (gRPC handles this automatically)
- Avoid proxy layers that add latency
