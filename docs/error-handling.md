# Error Handling

## Error Categories

### Authentication Errors

| Error | Description | Action |
|-------|-------------|--------|
| `UNAUTHENTICATED` | Token missing or invalid | Check token is set in metadata |
| `TOKEN_EXPIRED` | Token has expired | Obtain new token |
| `INVALID_TOKEN` | Token format incorrect | Verify token format |

### Subscription Limit Errors

| Error | Description | Limit |
|-------|-------------|-------|
| `SUBSCRIPTION_LIMIT_REACHED` | Total subscriptions exceeded | 6 total |
| `TRANSACTION_SUBSCRIPTION_LIMIT_REACHED` | Transaction streams exceeded | 2 max |
| `ACCOUNT_SUBSCRIPTION_LIMIT_REACHED` | Account streams exceeded | 5 max |
| `SLOT_SUBSCRIPTION_LIMIT_REACHED` | Slot streams exceeded | 2 max |
| `TOO_MANY_ACCOUNT_ADDRESSES` | Accounts in request exceeded | 100 per request |

### Connection Errors

| Error | Description | Action |
|-------|-------------|--------|
| `CONNECTION_ERROR` | Failed to establish connection | Check network/server address |
| `CONNECTION_CLOSED` | Server closed connection | Reconnect with backoff |
| `CONNECTION_TIMEOUT` | Connection timed out | Retry with backoff |
| `STREAM_CLOSED` | Stream was terminated | Reconnect |

### Request Errors

| Error | Description | Action |
|-------|-------------|--------|
| `INVALID_REQUEST` | Malformed request | Check request format |
| `INVALID_ACCOUNT_ADDRESS` | Bad account address format | Use valid Base58 address |
| `EMPTY_ACCOUNT_LIST` | No accounts provided | Add at least one account |

### Server Errors

| Error | Description | Action |
|-------|-------------|--------|
| `INTERNAL` | Internal server error | Contact support |
| `UNAVAILABLE` | Service unavailable | Retry with backoff |
| `RESOURCE_EXHAUSTED` | Server resources exhausted | Retry later |

---

## Recovery Strategies

### Reconnection with Exponential Backoff

```go
// Go example
func reconnectWithBackoff(ctx context.Context, maxRetries int) error {
    backoff := time.Second
    for i := 0; i < maxRetries; i++ {
        stream, err := client.SubscribeToTransactions(ctx)
        if err == nil {
            return processStream(stream)
        }

        log.Printf("Connection failed, retrying in %v", backoff)
        time.Sleep(backoff)
        backoff *= 2
        if backoff > 30*time.Second {
            backoff = 30 * time.Second
        }
    }
    return errors.New("max retries exceeded")
}
```

```rust
// Rust example
async fn reconnect_with_backoff(max_retries: u32) -> Result<(), Box<dyn Error>> {
    let mut backoff = Duration::from_secs(1);

    for _ in 0..max_retries {
        match client.subscribe_to_transactions().await {
            Ok(stream) => return process_stream(stream).await,
            Err(e) => {
                println!("Connection failed, retrying in {:?}", backoff);
                tokio::time::sleep(backoff).await;
                backoff = std::cmp::min(backoff * 2, Duration::from_secs(30));
            }
        }
    }
    Err("max retries exceeded".into())
}
```

```typescript
// TypeScript example
async function reconnectWithBackoff(maxRetries: number): Promise<void> {
    let backoff = 1000;

    for (let i = 0; i < maxRetries; i++) {
        try {
            await subscribeToTransactions();
            return;
        } catch (err) {
            console.log(`Connection failed, retrying in ${backoff}ms`);
            await new Promise(r => setTimeout(r, backoff));
            backoff = Math.min(backoff * 2, 30000);
        }
    }
    throw new Error('max retries exceeded');
}
```

### Error Type Classification

```go
// Go
func handleError(err error) {
    switch {
    case thorclient.IsStreamDone(err):
        // Normal closure, reconnect
        reconnect()
    case isAuthError(err):
        // Authentication issue, refresh token
        refreshToken()
    case isRateLimitError(err):
        // Too many subscriptions, close unused
        closeUnusedSubscriptions()
    default:
        // Unknown error, log and retry
        log.Printf("Unknown error: %v", err)
        reconnectWithBackoff()
    }
}
```

---

## Monitoring Recommendations

1. **Track error frequencies** — Alert on unusual patterns
2. **Monitor reconnection rates** — High rates indicate infrastructure issues
3. **Log error contexts** — Include timestamps and request details
4. **Set up health checks** — Verify stream connectivity periodically
