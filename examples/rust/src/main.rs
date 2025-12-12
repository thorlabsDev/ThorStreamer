//! Main entry point for the thorStreamerExample application.
//!
//! This file sets up the gRPC client, loads configuration from a JSON file,
//! and subscribes to different event streams (transactions, slot status updates,
//! wallet transactions, or account updates) based on user input. Detailed debugging
//! information is printed and logged for transaction and account events.

use std::fs::{File, OpenOptions};
use std::io::{Read, Write};
use std::{env, io};

use bs58;
use chrono::prelude::*;
use http_body::Body;
use prost::Message as ProstMessage;
use serde::{Deserialize, Serialize};
use serde_json;
use tonic::{metadata::MetadataValue, transport::Channel, Request};

// Generated protobuf modules
pub mod thor_streamer {
    pub mod types {
        tonic::include_proto!("thor_streamer.types");
    }
}

pub mod publisher {
    tonic::include_proto!("publisher");
}

// Include generated code for google.protobuf well-known types.
pub mod google {
    pub mod protobuf {
        tonic::include_proto!("google.protobuf");
    }
}

use google::protobuf::Empty as GoogleEmpty;

// The client from `publisher.proto`:
use publisher::event_publisher_client::EventPublisherClient;
use publisher::SubscribeAccountsRequest;

// Local types from events.proto
use thor_streamer::types::{Message, MessageWrapper, SlotStatusEvent, TransactionEvent, SubscribeUpdateAccountInfo};
use thor_streamer::types::message_wrapper::EventMessage;

/// Configuration settings for the application.
///
/// Contains the server address and the authorization token.
#[derive(Debug, Deserialize)]
struct Config {
    server_address: String,
    auth_token: String,
}

/// Structure for logging transaction signatures.
///
/// This log entry includes the timestamp, signature in Base58, slot number,
/// and a boolean indicating whether the transaction was successful.
#[derive(Debug, Serialize)]
struct SignatureLog {
    timestamp: String,
    signature: String,
    slot: u64,
    success: bool,
}

/// Structure for logging account updates.
///
/// This log entry includes the timestamp, account pubkey, owner, lamports,
/// and slot information.
#[derive(Debug, Serialize)]
struct AccountUpdateLog {
    timestamp: String,
    pubkey: String,
    owner: String,
    lamports: u64,
    slot: u64,
    executable: bool,
}

/// Loads the configuration file in JSON format from the same directory as the executable.
///
/// # Arguments
///
/// * `filename` - The name of the configuration file.
///
/// # Returns
///
/// * A `Config` struct if successful, or an error.
fn load_config(filename: &str) -> Result<Config, Box<dyn std::error::Error>> {
    let exe_path = env::current_exe()?;
    let exe_dir = exe_path.parent().ok_or("Failed to get executable directory")?;
    let config_path = exe_dir.join(filename);

    let mut file = File::open(&config_path)
        .map_err(|e| format!("Failed to read config file {}: {}", config_path.display(), e))?;

    let mut contents = String::new();
    file.read_to_string(&mut contents)?;

    let config: Config = serde_json::from_str(&contents)?;
    Ok(config)
}

/// Prints detailed debug information for an account update event and logs it.
///
/// The function prints account details such as pubkey, owner, lamports, executable status,
/// and logs the information in JSON format to a log file named "account_updates.log".
///
/// # Arguments
///
/// * `account` - A reference to the `SubscribeUpdateAccountInfo` to be debugged.
fn debug_account_update(account: &SubscribeUpdateAccountInfo) {
    println!("Account Update Debug Information:");
    println!("├─ Pubkey: {}", bs58::encode(&account.pubkey).into_string());
    println!("├─ Owner: {}", bs58::encode(&account.owner).into_string());
    println!("├─ Lamports: {}", account.lamports);
    println!("├─ Executable: {}", account.executable);
    println!("├─ Rent Epoch: {}", account.rent_epoch);
    println!("├─ Write Version: {}", account.write_version);
    println!("├─ Data Length: {} bytes", account.data.len());

    if let Some(sig) = &account.txn_signature {
        println!("├─ Transaction Signature: {}", bs58::encode(sig).into_string());
    } else {
        println!("├─ Transaction Signature: None");
    }

    if let Some(slot_status) = &account.slot {
        println!("├─ Slot Information:");
        println!("│  ├─ Slot: {}", slot_status.slot);
        println!("│  ├─ Parent: {}", slot_status.parent);

        let status_str = match slot_status.status {
            0 => "PROCESSED",
            1 => "CONFIRMED",
            2 => "ROOTED",
            _ => "UNKNOWN",
        };
        println!("│  ├─ Status: {}", status_str);

        if !slot_status.block_hash.is_empty() {
            println!("│  ├─ Block Hash: {}", bs58::encode(&slot_status.block_hash).into_string());
        }
        println!("│  └─ Block Height: {}", slot_status.block_height);
    } else {
        println!("└─ Slot Information: None");
    }

    // Log to file
    let exe_path = env::current_exe().expect("Failed to get executable path");
    let exe_dir = exe_path.parent().expect("Failed to get executable directory");
    let log_file_path = exe_dir.join("account_updates.log");

    let mut log_file = OpenOptions::new()
        .append(true)
        .create(true)
        .open(log_file_path)
        .expect("Failed to open log file");

    // Log JSON to file
    let timestamp = Utc::now().format("%Y-%m-%dT%H:%M:%S.000Z").to_string();
    let pubkey_base58 = bs58::encode(&account.pubkey).into_string();
    let owner_base58 = bs58::encode(&account.owner).into_string();

    let slot = account.slot.as_ref().map(|s| s.slot).unwrap_or(0);

    let account_log = AccountUpdateLog {
        timestamp,
        pubkey: pubkey_base58,
        owner: owner_base58,
        lamports: account.lamports,
        slot,
        executable: account.executable,
    };

    let json_log = serde_json::to_string(&account_log).expect("Failed to marshal log");
    writeln!(log_file, "{}", json_log).expect("Failed to write to log file");
}

/// Prints detailed debug information for a given transaction event and logs its signature.
///
/// The function prints transaction details such as signature, slot, index, vote status,
/// and logs the information in JSON format to a log file named "signatures.log".
///
/// # Arguments
///
/// * `tx` - A reference to the `TransactionEvent` to be debugged.
fn debug_transaction(tx: &TransactionEvent) {
    println!("Transaction Debug Information:");
    println!("├─ Signature: {}", bs58::encode(&tx.signature).into_string());
    println!("├─ Slot: {}", tx.slot);
    println!("├─ Index: {}", tx.index);
    println!("├─ Is Vote: {}", tx.is_vote);

    // Log to file
    let exe_path = env::current_exe().expect("Failed to get executable path");
    let exe_dir = exe_path.parent().expect("Failed to get executable directory");
    let log_file_path = exe_dir.join("signatures.log");

    let mut log_file = OpenOptions::new()
        .append(true)
        .create(true)
        .open(log_file_path)
        .expect("Failed to open log file");

    // Log JSON to file
    let timestamp = Utc::now().format("%Y-%m-%dT%H:%M:%S.000Z").to_string();
    let signature_base58 = bs58::encode(&tx.signature).into_string();

    let signature_log = SignatureLog {
        timestamp,
        signature: signature_base58,
        slot: tx.slot,
        success: !tx.transaction_status_meta.as_ref().map_or(false, |m| m.is_status_err),
    };

    let json_log = serde_json::to_string(&signature_log).expect("Failed to marshal log");
    writeln!(log_file, "{}", json_log).expect("Failed to write to log file");

    if tx.transaction.is_none() {
        println!("├─ Transaction is nil!");
        return;
    }

    let transaction = tx.transaction.as_ref().unwrap();
    if transaction.message.is_none() {
        println!("├─ Message is nil!");
        return;
    }

    let message = transaction.message.as_ref().unwrap();
    println!(
        "├─ Version: {} ({})",
        message.version,
        get_version_string(message.version)
    );

    debug_header(message);
    debug_account_keys(message);
    debug_blockhash(message);
    debug_instructions(message);
}

/// Prints header information from a transaction message.
///
/// # Arguments
///
/// * `msg` - A reference to the `Message` whose header is to be debugged.
fn debug_header(msg: &Message) {
    println!("├─ Header:");
    if msg.header.is_none() {
        println!("│  Header is nil!");
        return;
    }

    let header = msg.header.as_ref().unwrap();
    println!("│  ├─ NumRequiredSignatures: {}", header.num_required_signatures);
    println!("│  ├─ NumReadonlySignedAccounts: {}", header.num_readonly_signed_accounts);
    println!("│  └─ NumReadonlyUnsignedAccounts: {}", header.num_readonly_unsigned_accounts);
}

/// Prints account key information from a transaction message.
///
/// It prints up to the first five keys and indicates if there are more.
///
/// # Arguments
///
/// * `msg` - A reference to the `Message` containing the account keys.
fn debug_account_keys(msg: &Message) {
    println!("├─ Account Keys ({}):", msg.account_keys.len());
    if msg.account_keys.is_empty() {
        println!("│  No account keys!");
        return;
    }

    for (i, key) in msg.account_keys.iter().enumerate().take(5) {
        println!("│  ├─ [{}]: {}", i, bs58::encode(key).into_string());
    }

    if msg.account_keys.len() > 5 {
        println!("│  └─ ... and {} more keys", msg.account_keys.len() - 5);
    }
}

/// Prints the recent blockhash from a transaction message.
///
/// # Arguments
///
/// * `msg` - A reference to the `Message` containing the recent blockhash.
fn debug_blockhash(msg: &Message) {
    println!("├─ Recent Blockhash:");
    if msg.recent_block_hash.is_empty() {
        println!("│  Blockhash is empty!");
        return;
    }
    println!("│  └─ {}", bs58::encode(&msg.recent_block_hash).into_string());
}

/// Prints the instruction details from a transaction message.
///
/// It prints up to the first three instructions and indicates if there are more.
///
/// # Arguments
///
/// * `msg` - A reference to the `Message` containing the instructions.
fn debug_instructions(msg: &Message) {
    println!("├─ Instructions ({}):", msg.instructions.len());
    if msg.instructions.is_empty() {
        println!("│  No instructions!");
        return;
    }

    for (i, ix) in msg.instructions.iter().enumerate().take(3) {
        println!("│  ├─ Instruction {}:", i);
        println!("│  │  ├─ Program ID Index: {}", ix.program_id_index);
        println!("│  │  ├─ Account Indexes: {:?}", ix.accounts);
        println!("│  │  └─ Data Length: {} bytes", ix.data.len());
    }

    if msg.instructions.len() > 3 {
        println!("│  └─ ... and {} more instructions", msg.instructions.len() - 3);
    }
}

/// Returns a human-readable string for the message version.
///
/// # Arguments
///
/// * `version` - The version number of the message.
///
/// # Returns
///
/// * A string representing the version (e.g., "Legacy" for 0, "V0" for 1).
fn get_version_string(version: u32) -> String {
    match version {
        0 => "Legacy".to_string(),
        1 => "V0".to_string(),
        _ => format!("Unknown({})", version),
    }
}

/// Prints debug information for a slot status event.
///
/// # Arguments
///
/// * `slot_ev` - A reference to the `SlotStatusEvent` to be debugged.
fn debug_slot_event(slot_ev: &SlotStatusEvent) {
    println!("Slot Debug Information:");
    println!("├─ Slot: {}", slot_ev.slot);
    println!("├─ Parent: {}", slot_ev.parent);

    let status_str = match slot_ev.status {
        0 => "PROCESSED",
        1 => "CONFIRMED",
        2 => "ROOTED",
        _ => "UNKNOWN",
    };
    println!("├─ Status: {}", status_str);

    if !slot_ev.block_hash.is_empty() {
        println!("├─ Block Hash: {}", bs58::encode(&slot_ev.block_hash).into_string());
    } else {
        println!("├─ Block Hash: (empty)");
    }

    println!("└─ Block Height: {}", slot_ev.block_height);
}

/// Subscribes to all transaction events from the gRPC server and debugs each event.
///
/// # Arguments
///
/// * `client` - A gRPC client for the `EventPublisher` service.
///
/// # Returns
///
/// * `Ok(())` if the subscription and debugging complete successfully; otherwise an error.
async fn subscribe_to_transactions<T>(
    mut client: EventPublisherClient<T>
) -> Result<(), Box<dyn std::error::Error>>
where
    T: tonic::client::GrpcService<tonic::body::BoxBody>,
    T::Error: Into<Box<dyn std::error::Error + Send + Sync>>,
    T::ResponseBody: Body<Data = tonic::codegen::Bytes> + Send + 'static,
    <T::ResponseBody as Body>::Error: Into<Box<dyn std::error::Error + Send + Sync>> + Send,
{
    // Using google.protobuf.Empty defined as GoogleEmpty
    let request = Request::new(GoogleEmpty {});
    let mut stream = client.subscribe_to_transactions(request).await?.into_inner();

    while let Some(resp) = stream.message().await? {
        let msg_wrapper = MessageWrapper::decode(&resp.data[..])?;

        if let Some(EventMessage::Transaction(tx)) = msg_wrapper.event_message {
            debug_transaction(&tx);
        } else {
            println!("Received a non-transaction event");
        }
    }

    Ok(())
}

/// Subscribes to slot status events from the gRPC server and debugs each event.
///
/// # Arguments
///
/// * `client` - A gRPC client for the `EventPublisher` service.
///
/// # Returns
///
/// * `Ok(())` if the subscription and debugging complete successfully; otherwise an error.
async fn subscribe_to_slots<T>(
    mut client: EventPublisherClient<T>
) -> Result<(), Box<dyn std::error::Error>>
where
    T: tonic::client::GrpcService<tonic::body::BoxBody>,
    T::Error: Into<Box<dyn std::error::Error + Send + Sync>>,
    T::ResponseBody: Body<Data = tonic::codegen::Bytes> + Send + 'static,
    <T::ResponseBody as Body>::Error: Into<Box<dyn std::error::Error + Send + Sync>> + Send,
{
    let request = Request::new(GoogleEmpty {});
    let mut stream = client.subscribe_to_slot_status(request).await?.into_inner();

    while let Some(resp) = stream.message().await? {
        let msg_wrapper = MessageWrapper::decode(&resp.data[..])?;

        if let Some(EventMessage::Slot(slot_event)) = msg_wrapper.event_message {
            debug_slot_event(&slot_event);
        } else {
            println!("Received a non-slot event from SubscribeToSlotStatus");
        }
    }

    Ok(())
}

/// Subscribes to account update events for specified account addresses and/or owner addresses.
///
/// # Arguments
///
/// * `client` - A gRPC client for the `EventPublisher` service.
/// * `account_addresses` - A vector of account addresses (as strings) to subscribe to.
/// * `owner_addresses` - A vector of owner addresses (as strings) to filter by.
///
/// # Returns
///
/// * `Ok(())` if the subscription and debugging complete successfully; otherwise an error.
async fn subscribe_to_account_updates<T>(
    mut client: EventPublisherClient<T>,
    account_addresses: Vec<String>,
    owner_addresses: Vec<String>,
) -> Result<(), Box<dyn std::error::Error>>
where
    T: tonic::client::GrpcService<tonic::body::BoxBody>,
    T::Error: Into<Box<dyn std::error::Error + Send + Sync>>,
    T::ResponseBody: Body<Data = tonic::codegen::Bytes> + Send + 'static,
    <T::ResponseBody as Body>::Error: Into<Box<dyn std::error::Error + Send + Sync>> + Send,
{
    let request = Request::new(SubscribeAccountsRequest {
        account_address: account_addresses,
        owner_address: owner_addresses,
    });

    let mut stream = client.subscribe_to_account_updates(request).await?.into_inner();

    println!("Subscribed to account updates. Waiting for updates...");

    while let Some(resp) = stream.message().await? {
        let msg_wrapper = MessageWrapper::decode(&resp.data[..])?;

        if let Some(EventMessage::AccountUpdate(account_update)) = msg_wrapper.event_message {
            debug_account_update(&account_update);
        } else {
            println!("Received a non-account-update event");
        }
    }

    Ok(())
}

/// Entry point of the application.
///
/// Loads the configuration, establishes a gRPC client connection with authorization,
/// and subscribes to event streams based on user input.
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = load_config("config.json")?;

    // Clone the token from the configuration so it remains valid.
    let token_string = config.auth_token.clone();
    let channel = Channel::from_shared(config.server_address)?.connect().await?;

    let add_auth = move |mut req: Request<()>| {
        let meta = MetadataValue::try_from(token_string.as_str())
            .expect("Invalid token string for metadata");
        req.metadata_mut().insert("authorization", meta);
        Ok(req)
    };

    // Create a publisher client with an interceptor that adds the auth token.
    let publisher_client = EventPublisherClient::with_interceptor(channel, add_auth);

    println!("Select subscription method:");
    println!("1: All Transactions");
    println!("2: Slot Status Updates");
    println!("3: Account Updates");

    let mut choice = String::new();
    io::stdin().read_line(&mut choice)?;
    let choice: u32 = choice.trim().parse()?;

    match choice {
        1 => subscribe_to_transactions(publisher_client).await?,
        2 => subscribe_to_slots(publisher_client).await?,
        3 => {
            println!("Enter account addresses to monitor (comma-separated, or press Enter to skip):");
            let mut account_input = String::new();
            io::stdin().read_line(&mut account_input)?;
            let accounts: Vec<String> = if account_input.trim().is_empty() {
                vec![]
            } else {
                account_input.trim().split(',').map(|s| s.trim().to_string()).collect()
            };

            println!("Enter owner addresses to filter by (comma-separated, or press Enter to skip):");
            let mut owner_input = String::new();
            io::stdin().read_line(&mut owner_input)?;
            let owners: Vec<String> = if owner_input.trim().is_empty() {
                vec![]
            } else {
                owner_input.trim().split(',').map(|s| s.trim().to_string()).collect()
            };

            if accounts.is_empty() && owners.is_empty() {
                println!("At least one account address or owner address is required");
            } else {
                println!("Subscribing to account updates:");
                if !accounts.is_empty() {
                    println!("  Account addresses: {:?}", accounts);
                }
                if !owners.is_empty() {
                    println!("  Owner addresses: {:?}", owners);
                }
                subscribe_to_account_updates(publisher_client, accounts, owners).await?
            }
        }
        _ => println!("Invalid choice"),
    }

    Ok(())
}