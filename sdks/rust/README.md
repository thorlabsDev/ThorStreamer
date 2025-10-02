# ThorStreamer Rust SDK

Official Rust client library for ThorStreamer gRPC services.

## Installation

Add to your `Cargo.toml`:

```toml
[dependencies]
thor-grpc-client = "0.1.0"
tokio = { version = "1", features = ["full"] }
```

## Quick Start

```rust
use thor_grpc_client::{ClientConfig, ThorClient, parse_message};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = ClientConfig {
        server_addr: "http://your-server-address:50051".to_string(),
        token: "your-auth-token".to_string(),
        ..Default::default()
    };

    let mut client = ThorClient::new(config).await?;

    // Subscribe to transactions
    let mut stream = client.subscribe_to_transactions().await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        
        use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;
        if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
            if let Some(tx) = tx_wrapper.transaction {
                println!("Transaction: slot={}", tx.slot);
            }
        }
    }

    Ok(())
}
```

## Environment Setup

Using environment variables:

```rust
use std::env;

let config = ClientConfig {
    server_addr: env::var("SERVER_ADDRESS")
        .unwrap_or_else(|_| "http://your-server-address:50051".to_string()),
    token: env::var("AUTH_TOKEN")
        .expect("AUTH_TOKEN environment variable not set"),
    ..Default::default()
};
```

Or create a `.env` file:

```bash
SERVER_ADDRESS=http://your-server:50051
AUTH_TOKEN=your-auth-token
```

Then load it with the `dotenv` crate:

```toml
[dependencies]
dotenv = "0.15"
```

```rust
use dotenv::dotenv;

dotenv().ok();
let token = std::env::var("AUTH_TOKEN")?;
```

## Features

- **Transaction Streaming**: Real-time Solana transaction updates
- **Slot Status**: Subscribe to slot confirmations and updates
- **Wallet Tracking**: Monitor specific wallet addresses (up to 10 per request)
- **Account Updates**: Track account state changes with owner filtering
- **Async/Await**: Built on Tokio for efficient async operations
- **Type Safety**: Full protobuf type definitions with compile-time checking

## API Reference

### Creating a Client

```rust
use thor_grpc_client::{ClientConfig, ThorClient};
use std::time::Duration;

let config = ClientConfig {
    server_addr: "http://your-server-address:50051".to_string(),
    token: "your-auth-token".to_string(),
    timeout: Duration::from_secs(30),
};

let client = ThorClient::new(config).await?;
```

### Subscribe to Transactions

```rust
use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

let mut stream = client.subscribe_to_transactions().await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    
    if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
        if let Some(tx) = tx_wrapper.transaction {
            let sig_hex = tx.signature.iter()
                .take(8)
                .map(|b| format!("{:02x}", b))
                .collect::<Vec<_>>()
                .join("");
            println!("Transaction: slot={}, signature={}", tx.slot, sig_hex);
        }
    }
}
```

### Subscribe to Slot Status

```rust
let mut stream = client.subscribe_to_slot_status().await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    
    if let Some(EventMessage::Slot(slot)) = msg.event_message {
        println!("Slot: slot={}, status={}, height={}", 
            slot.slot, slot.status, slot.block_height);
    }
}
```

### Subscribe to Wallet Transactions

Monitor specific wallet addresses (max 10 per request):

```rust
let wallets = vec![
    "wallet1base58address...".to_string(),
    "wallet2base58address...".to_string(),
];

let mut stream = client.subscribe_to_wallet_transactions(wallets).await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    
    if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
        if let Some(tx) = tx_wrapper.transaction {
            println!("Wallet transaction: slot={}", tx.slot);
        }
    }
}
```

### Subscribe to Account Updates

Track account state changes with optional owner filtering:

```rust
let accounts = vec!["account1...".to_string(), "account2...".to_string()];
let owners = vec!["owner1...".to_string()];

let mut stream = client.subscribe_to_account_updates(accounts, owners).await?;

while let Some(response) = stream.message().await? {
    let msg = parse_message(&response.data)?;
    
    if let Some(EventMessage::AccountUpdate(update)) = msg.event_message {
        let pubkey_hex = update.pubkey.iter()
            .take(8)
            .map(|b| format!("{:02x}", b))
            .collect::<Vec<_>>()
            .join("");
        println!("Account: pubkey={}, lamports={}", 
            pubkey_hex, update.lamports);
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

## Examples

See the [examples directory](examples/) for complete working examples:

```bash
cargo run --example subscribe
```

## Building from Source

```bash
# Clone repository
git clone https://github.com/thorlabsDev/ThorStreamer.git
cd ThorStreamer/sdks/rust

# Build the library
cargo build --release

# Run tests
cargo test

# Run example
cargo run --example subscribe
```

## Module Information

```toml
[package]
name = "thor-grpc-client"
version = "0.1.0"
edition = "2021"

[dependencies]
tonic = "0.12.3"
prost = "0.13.5"
tokio = { version = "1", features = ["full"] }
prost-types = "0.13.5"

[build-dependencies]
tonic-build = "0.12.3"
```

## Generated Code

The SDK uses `tonic-build` to generate Rust code from protobuf definitions at compile time. Proto files are located in the `proto/` directory:

- `events.proto` - ThorStreamer message types
- `publisher.proto` - EventPublisher service definitions

## Version History

- `v0.1.0` - Initial release
    - Transaction streaming
    - Slot status updates
    - Wallet tracking
    - Account updates with owner filtering
    - Async/await support with Tokio

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/thorlabsDev/ThorStreamer/blob/master/LICENSE) file for details.

## Support

- **GitHub Issues**: [https://github.com/thorlabsDev/ThorStreamer/issues](https://github.com/thorlabsDev/ThorStreamer/issues)
- **Documentation**: [https://docs.rs/thor-grpc-client](https://docs.rs/thor-grpc-client)
- **Crates.io**: [https://crates.io/crates/thor-grpc-client](https://crates.io/crates/thor-grpc-client)
- **Main Repository**: [https://github.com/thorlabsDev/ThorStreamer](https://github.com/thorlabsDev/ThorStreamer)
