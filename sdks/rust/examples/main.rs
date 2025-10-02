// examples/main.rs
use thor_grpc_client::{ClientConfig, ThorClient, parse_message};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load .env file
    dotenv::dotenv().ok();

    env_logger::init();

    let config = ClientConfig {
        server_addr: std::env::var("SERVER_ADDRESS")
            .unwrap_or_else(|_| "http://localhost:50051".to_string()),
        token: std::env::var("AUTH_TOKEN")
            .unwrap_or_else(|_| String::new()),
        ..Default::default()
    };


    let client = ThorClient::new(config).await?;

    println!("Choose a subscription:");
    println!("1. Transactions");
    println!("2. Slots");
    println!("3. Wallets");
    println!("4. Accounts");

    let mut input = String::new();
    std::io::stdin().read_line(&mut input)?;

    match input.trim() {
        "1" => subscribe_transactions(client).await?,
        "2" => subscribe_slots(client).await?,
        "3" => subscribe_wallets(client).await?,
        "4" => subscribe_accounts(client).await?,
        _ => println!("Invalid choice"),
    }

    Ok(())
}

async fn subscribe_transactions(mut client: ThorClient) -> Result<(), Box<dyn std::error::Error>> {
    use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

    println!("Subscribed to transactions");
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
                println!("Received transaction: slot={}, signature={}", tx.slot, sig_hex);
            }
        }
    }

    Ok(())
}

async fn subscribe_slots(mut client: ThorClient) -> Result<(), Box<dyn std::error::Error>> {
    use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

    println!("Subscribed to slots");
    let mut stream = client.subscribe_to_slot_status().await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        if let Some(EventMessage::Slot(slot)) = msg.event_message {
            println!("Received slot: slot={}, status={}, height={}",
                     slot.slot, slot.status, slot.block_height);
        }
    }

    Ok(())
}

async fn subscribe_wallets(mut client: ThorClient) -> Result<(), Box<dyn std::error::Error>> {
    use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

    let wallets = vec!["YourWalletAddressHere".to_string()];
    println!("Subscribed to {} wallets", wallets.len());
    let mut stream = client.subscribe_to_wallet_transactions(wallets).await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        if let Some(EventMessage::Transaction(tx_wrapper)) = msg.event_message {
            if let Some(tx) = tx_wrapper.transaction {
                let sig_hex = tx.signature.iter()
                    .take(8)
                    .map(|b| format!("{:02x}", b))
                    .collect::<Vec<_>>()
                    .join("");
                println!("Received wallet transaction: slot={}, signature={}", tx.slot, sig_hex);
            }
        }
    }

    Ok(())
}

async fn subscribe_accounts(mut client: ThorClient) -> Result<(), Box<dyn std::error::Error>> {
    use thor_grpc_client::proto::thor_streamer::types::message_wrapper::EventMessage;

    let accounts = vec!["AccountAddress1".to_string()];
    let owners = vec!["OwnerAddress1".to_string()];
    println!("Subscribed to {} accounts and {} owners", accounts.len(), owners.len());
    let mut stream = client.subscribe_to_account_updates(accounts, owners).await?;

    while let Some(response) = stream.message().await? {
        let msg = parse_message(&response.data)?;
        if let Some(EventMessage::AccountUpdate(update)) = msg.event_message {
            let pubkey_hex = update.pubkey.iter()
                .take(8)
                .map(|b| format!("{:02x}", b))
                .collect::<Vec<_>>()
                .join("");
            println!("Received account update: pubkey={}, lamports={}",
                     pubkey_hex, update.lamports);
        }
    }

    Ok(())
}