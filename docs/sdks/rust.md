# Rust SDK

Official Rust client library for ThorStreamer with async/await support.

## Installation

Add to `Cargo.toml`:

```toml
[dependencies]
thorstreamer-grpc-client = "0.1"
tokio = { version = "1", features = ["full"] }
```

## Quick Start

```rust
use thorstreamer_grpc_client::{ClientConfig, ThorClient, parse_message};
use thorstreamer_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = ClientConfig {
        server_addr: "http://your-server:50051".to_string(),
        token: "your-token".to_string(),
        ..Default::default()
    };

    let mut client = ThorClient::new(config).await?;
    let mut stream = client.subscribe_to_transactions().await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
            if let Some(tx) = tx_wrapper.transaction {
                println!("Transaction: slot={}", tx.slot);
            }
        }
    }
    Ok(())
}
```

## API Reference

### Creating a Client

```rust
use thorstreamer_grpc_client::{ClientConfig, ThorClient};
use std::time::Duration;

let config = ClientConfig {
    server_addr: "http://server:50051".to_string(),
    token: "your-token".to_string(),
    timeout: Duration::from_secs(30),
};

let client = ThorClient::new(config).await?;
```

### subscribe_to_transactions

```rust
use thorstreamer_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

let mut stream = client.subscribe_to_transactions().await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
        if let Some(tx) = tx_wrapper.transaction {
            let sig_hex: String = tx.signature.iter()
                .take(8)
                .map(|b| format!("{:02x}", b))
                .collect();
            println!("Transaction: slot={}, sig={}", tx.slot, sig_hex);
        }
    }
}
```

### subscribe_to_slot_status

```rust
let mut stream = client.subscribe_to_slot_status().await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    if let Some(EventMessage::Slot(slot)) = msg.event_message {
        println!("Slot: {}, status={}, height={}",
            slot.slot, slot.status, slot.block_height);
    }
}
```

### subscribe_to_wallet_transactions

Monitor up to 10 wallet addresses:

```rust
let wallets = vec![
    "9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM".to_string(),
    "2ojv9BAiHUrvsm9gxDe7fJSzbNZSJcxZvf8dqmWGHG8S".to_string(),
];

let mut stream = client.subscribe_to_wallet_transactions(wallets).await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
        println!("Wallet transaction received");
    }
}
```

### subscribe_to_account_updates

Monitor accounts with optional owner filtering:

```rust
let accounts = vec!["account1...".to_string()];
let owners = vec!["TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA".to_string()];

let mut stream = client.subscribe_to_account_updates(accounts, owners).await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    if let Some(EventMessage::AccountUpdate(update)) = msg.event_message {
        println!("Account: lamports={}", update.lamports);
    }
}
```

## Error Handling

```rust
use tonic::Status;

match stream.message().await {
    Ok(Some(response)) => {
        // Process message
    }
    Ok(None) => {
        // Stream ended normally
        println!("Stream closed");
    }
    Err(status) => {
        // Handle gRPC error
        eprintln!("Stream error: {:?}", status);
    }
}
```

## Environment Variables

Using `dotenv`:

```toml
[dependencies]
dotenv = "0.15"
```

```rust
use dotenv::dotenv;

dotenv().ok();
let config = ClientConfig {
    server_addr: std::env::var("SERVER_ADDRESS")?,
    token: std::env::var("AUTH_TOKEN")?,
    ..Default::default()
};
```

## Full Example

See [examples/rust](https://github.com/thorlabsDev/ThorStreamer/tree/master/examples/rust) for a complete implementation.

## Resources

- [crates.io](https://crates.io/crates/thorstreamer-grpc-client)
- [docs.rs Documentation](https://docs.rs/thorstreamer-grpc-client)
- [GitHub Repository](https://github.com/thorlabsDev/ThorStreamer)
