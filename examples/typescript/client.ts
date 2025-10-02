import * as grpc from "@grpc/grpc-js";
import * as fs from "fs";
import * as path from "path";
import * as readline from "readline";
import base58 from "bs58";

// Import all proto definitions from the consolidated index
import {
  EventPublisherClient,
  SubscribeWalletRequest,
  SubscribeAccountsRequest,
  MessageWrapper,
  TransactionEvent,
  SlotStatusEvent,
  SubscribeUpdateAccountInfo,
  Message,
  CompiledInstruction,
  Empty,
  IEventPublisherClient,
} from "./proto";

// ============================================================================
// Type Definitions
// ============================================================================

interface Config {
  serverAddress: string;
  authToken: string;
}

interface SignatureLog {
  timestamp: string;
  signature: string;
  slot: number;
  success: boolean;
}

interface AccountUpdateLog {
  timestamp: string;
  pubkey: string;
  owner: string;
  lamports: number;
  slot: number;
  executable: boolean;
}

interface StreamResponse {
  data: Uint8Array;
}

// ============================================================================
// Configuration and Setup
// ============================================================================

class ThorStreamerClient {
  private client: IEventPublisherClient;
  private metadata: grpc.Metadata;
  private config: Config;
  private rl: readline.Interface;
  private readonly SIGNATURES_LOG_PATH: string;
  private readonly ACCOUNT_UPDATES_LOG_PATH: string;

  constructor() {
    this.config = this.loadConfig();
    this.SIGNATURES_LOG_PATH = path.join(__dirname, "signatures.log");
    this.ACCOUNT_UPDATES_LOG_PATH = path.join(__dirname, "account_updates.log");

    // Setup readline interface
    this.rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    // Setup gRPC client
    this.client = new EventPublisherClient(
        this.config.serverAddress,
        grpc.credentials.createInsecure()
    );

    // Setup metadata with auth token
    this.metadata = new grpc.Metadata();
    this.metadata.add("authorization", this.config.authToken);
  }

  // --------------------------------------------------------------------------
  // Configuration Loading
  // --------------------------------------------------------------------------

  private loadConfig(): Config {
    try {
      const configPath = path.join(__dirname, "config.json");
      const configData = fs.readFileSync(configPath, "utf8");
      return JSON.parse(configData);
    } catch (error) {
      console.error("⚠️  Failed to load config.json:", error);
      console.log("Using default configuration values...");
      return {
        serverAddress: "http://localhost:10000",
        authToken: "your-auth-token",
      };
    }
  }

  // --------------------------------------------------------------------------
  // Utility Methods
  // --------------------------------------------------------------------------

  private question(prompt: string): Promise<string> {
    return new Promise((resolve) => {
      this.rl.question(prompt, (answer) => {
        resolve(answer);
      });
    });
  }

  private appendToLogFile(filePath: string, data: object): void {
    try {
      const jsonLine = JSON.stringify(data) + "\n";
      fs.appendFileSync(filePath, jsonLine);
    } catch (error) {
      console.error(`❌ Failed to write to log file ${filePath}:`, error);
    }
  }

  private getVersionString(version: number): string {
    switch (version) {
      case 0:
        return "Legacy";
      case 1:
        return "V0";
      default:
        return `Unknown(${version})`;
    }
  }

  private getStatusString(status: number): string {
    switch (status) {
      case 0:
        return "PROCESSED";
      case 1:
        return "CONFIRMED";
      case 2:
        return "ROOTED";
      case 3:
        return "FIRST_SHRED_RECEIVED";
      case 4:
        return "COMPLETED";
      case 5:
        return "CREATED_BANK";
      case 6:
        return "DEAD";
      default:
        return "UNKNOWN";
    }
  }

  // --------------------------------------------------------------------------
  // Debug Methods
  // --------------------------------------------------------------------------

  private debugTransaction(tx: any): void {
    const signature = tx.getSignature_asU8();
    const slot = tx.getSlot();
    const index = tx.getIndex();
    const isVote = tx.getIsVote();

    console.log("\n📋 Transaction Debug Information:");
    console.log(`├─ Signature: ${base58.encode(signature)}`);
    console.log(`├─ Slot: ${slot}`);
    console.log(`├─ Index: ${index}`);
    console.log(`├─ Is Vote: ${isVote}`);

    // Log to file
    const timestamp = new Date().toISOString();
    const signatureBase58 = base58.encode(signature);

    const statusMeta = tx.getTransactionStatusMeta();
    const isStatusErr = statusMeta?.getIsStatusErr() || false;

    const logEntry: SignatureLog = {
      timestamp,
      signature: signatureBase58,
      slot: Number(slot),
      success: !isStatusErr,
    };

    this.appendToLogFile(this.SIGNATURES_LOG_PATH, logEntry);

    const transaction = tx.getTransaction();
    if (!transaction) {
      console.log("├─ ⚠️  Transaction is nil!");
      return;
    }

    const message = transaction.getMessage();
    if (!message) {
      console.log("├─ ⚠️  Message is nil!");
      return;
    }

    const version = message.getVersion();
    console.log(
        `├─ Version: ${version} (${this.getVersionString(version)})`
    );

    this.debugHeader(message);
    this.debugAccountKeys(message);
    this.debugBlockhash(message);
    this.debugInstructions(message);
  }

  private debugAccountUpdate(account: any): void {
    const pubkey = account.getPubkey_asU8();
    const owner = account.getOwner_asU8();
    const lamports = account.getLamports();
    const executable = account.getExecutable();
    const rentEpoch = account.getRentEpoch();
    const writeVersion = account.getWriteVersion();
    const data = account.getData_asU8();
    const txnSignature = account.getTxnSignature_asU8();
    const slotInfo = account.getSlot();

    console.log("\n📋 Account Update Debug Information:");
    console.log(`├─ Pubkey: ${base58.encode(pubkey)}`);
    console.log(`├─ Owner: ${base58.encode(owner)}`);
    console.log(`├─ Lamports: ${lamports}`);
    console.log(`├─ Executable: ${executable}`);
    console.log(`├─ Rent Epoch: ${rentEpoch}`);
    console.log(`├─ Write Version: ${writeVersion}`);
    console.log(`├─ Data Length: ${data.length} bytes`);

    if (txnSignature && txnSignature.length > 0) {
      console.log(`├─ Transaction Signature: ${base58.encode(txnSignature)}`);
    } else {
      console.log("├─ Transaction Signature: None");
    }

    if (slotInfo) {
      console.log("├─ Slot Information:");
      console.log(`│  ├─ Slot: ${slotInfo.getSlot()}`);
      console.log(`│  ├─ Parent: ${slotInfo.getParent()}`);

      const statusStr = this.getStatusString(slotInfo.getStatus());
      console.log(`│  ├─ Status: ${statusStr}`);

      const blockHash = slotInfo.getBlockHash_asU8();
      if (blockHash && blockHash.length > 0) {
        console.log(`│  ├─ Block Hash: ${base58.encode(blockHash)}`);
      }
      console.log(`│  └─ Block Height: ${slotInfo.getBlockHeight()}`);
    } else {
      console.log("└─ Slot Information: None");
    }

    // Log to file
    const timestamp = new Date().toISOString();
    const pubkeyBase58 = base58.encode(pubkey);
    const ownerBase58 = base58.encode(owner);

    const slot = slotInfo ? Number(slotInfo.getSlot()) : 0;

    const logEntry: AccountUpdateLog = {
      timestamp,
      pubkey: pubkeyBase58,
      owner: ownerBase58,
      lamports: Number(lamports),
      slot,
      executable,
    };

    this.appendToLogFile(this.ACCOUNT_UPDATES_LOG_PATH, logEntry);
  }

  private debugSlotEvent(slotEv: any): void {
    const slot = slotEv.getSlot();
    const parent = slotEv.getParent();
    const status = slotEv.getStatus();
    const blockHash = slotEv.getBlockHash_asU8();
    const blockHeight = slotEv.getBlockHeight();

    console.log("\n📋 Slot Debug Information:");
    console.log(`├─ Slot: ${slot}`);
    console.log(`├─ Parent: ${parent}`);

    const statusStr = this.getStatusString(status);
    console.log(`├─ Status: ${statusStr}`);

    if (blockHash && blockHash.length > 0) {
      console.log(`├─ Block Hash: ${base58.encode(blockHash)}`);
    } else {
      console.log("├─ Block Hash: (empty)");
    }

    console.log(`└─ Block Height: ${blockHeight}`);
  }

  private debugHeader(msg: any): void {
    console.log("├─ Header:");
    const header = msg.getHeader();
    if (!header) {
      console.log("│  └─ ⚠️  Header is nil!");
      return;
    }
    console.log(`│  ├─ NumRequiredSignatures: ${header.getNumRequiredSignatures()}`);
    console.log(`│  ├─ NumReadonlySignedAccounts: ${header.getNumReadonlySignedAccounts()}`);
    console.log(`│  └─ NumReadonlyUnsignedAccounts: ${header.getNumReadonlyUnsignedAccounts()}`);
  }

  private debugAccountKeys(msg: any): void {
    const accountKeys = msg.getAccountKeysList();
    console.log(`├─ Account Keys (${accountKeys.length}):`);
    if (accountKeys.length === 0) {
      console.log("│  └─ ⚠️  No account keys!");
      return;
    }

    accountKeys.slice(0, 5).forEach((key: Uint8Array | string, i: number) => {
      const keyBytes = typeof key === 'string' ? Buffer.from(key, 'base64') : key;
      console.log(`│  ├─ [${i}]: ${base58.encode(keyBytes)}`);
    });
    if (accountKeys.length > 5) {
      console.log(`│  └─ ... and ${accountKeys.length - 5} more keys`);
    }
  }

  private debugBlockhash(msg: any): void {
    const blockhash = msg.getRecentBlockHash_asU8();
    console.log("├─ Recent Blockhash:");
    if (!blockhash || blockhash.length === 0) {
      console.log("│  └─ ⚠️  Blockhash is empty!");
      return;
    }
    console.log(`│  └─ ${base58.encode(blockhash)}`);
  }

  private debugInstructions(msg: any): void {
    const instructions = msg.getInstructionsList();
    console.log(`├─ Instructions (${instructions.length}):`);
    if (instructions.length === 0) {
      console.log("│  └─ ⚠️  No instructions!");
      return;
    }

    instructions.slice(0, 3).forEach((ix: any, i: number) => {
      console.log(`│  ├─ Instruction ${i}:`);
      console.log(`│  │  ├─ Program ID Index: ${ix.getProgramIdIndex()}`);
      const accounts = ix.getAccountsList();
      console.log(`│  │  ├─ Account Indexes: [${accounts.join(", ")}]`);
      const data = ix.getData_asU8();
      console.log(`│  │  └─ Data Length: ${data.length} bytes`);
    });
    if (instructions.length > 3) {
      console.log(`│  └─ ... and ${instructions.length - 3} more instructions`);
    }
  }

  // --------------------------------------------------------------------------
  // Subscription Methods
  // --------------------------------------------------------------------------

  private async subscribeToTransactions(): Promise<void> {
    console.log(`\n🔍 Subscribing to all transactions...`);

    const emptyRequest = new Empty();
    const transactionStream = this.client.subscribeToTransactions(emptyRequest, this.metadata);

    transactionStream.on("data", (resp: any) => {
      try {
        // Debug the response structure
        console.log("Response type:", typeof resp);
        console.log("Response keys:", Object.keys(resp));

        // Get the binary data - it should be in resp.data
        const binaryData = resp.data || resp.getData?.() || resp;

        const msgWrapper = MessageWrapper.deserializeBinary(binaryData);

        // Debug what we received
        const eventCase = msgWrapper.getEventMessageCase();
        console.log("EventMessageCase:", eventCase);
        console.log("Available cases:", MessageWrapper.EventMessageCase);

        // Check each possible case
        if (eventCase === MessageWrapper.EventMessageCase.TRANSACTION) {
          const txWrapper = msgWrapper.getTransaction();
          if (txWrapper?.hasTransaction()) {
            const transaction = txWrapper.getTransaction();
            if (transaction) {
              this.debugTransaction(transaction);
            }
          }
        } else if (eventCase === MessageWrapper.EventMessageCase.SLOT) {
          console.log("⚠️ Received SLOT event in transaction stream");
        } else if (eventCase === MessageWrapper.EventMessageCase.ACCOUNT_UPDATE) {
          console.log("⚠️ Received ACCOUNT_UPDATE event in transaction stream");
        } else {
          console.log(`⚠️ Unknown event case: ${eventCase}`);
        }
      } catch (error) {
        console.error("❌ Failed to process message:", error);
      }
    });

    transactionStream.on("error", (error: grpc.ServiceError) => {
      console.error("❌ Stream error:", error);
    });

    transactionStream.on("end", () => {
      console.log("ℹ️ Stream ended.");
    });
  }

  private async subscribeToSlots(): Promise<void> {
    console.log(`\n🔍 Subscribing to slot status updates...`);

    const emptyRequest = new Empty();
    const slotStream = this.client.subscribeToSlotStatus(emptyRequest, this.metadata);

    slotStream.on("data", (resp: any) => {
      try {
        const binaryData = resp.data || resp.getData?.() || resp;
        const msgWrapper = MessageWrapper.deserializeBinary(binaryData);

        const eventCase = msgWrapper.getEventMessageCase();

        if (eventCase === MessageWrapper.EventMessageCase.SLOT) {
          const slotEvent = msgWrapper.getSlot();
          if (slotEvent) {
            this.debugSlotEvent(slotEvent);
          }
        } else if (eventCase === MessageWrapper.EventMessageCase.TRANSACTION) {
          console.log("⚠️ Received TRANSACTION event in slot stream");
        } else if (eventCase === MessageWrapper.EventMessageCase.ACCOUNT_UPDATE) {
          console.log("⚠️ Received ACCOUNT_UPDATE event in slot stream");
        } else {
          console.log(`⚠️ Unknown event case: ${eventCase}`);
        }
      } catch (error) {
        console.error("❌ Failed to deserialize MessageWrapper:", error);
      }
    });

    slotStream.on("error", (error: grpc.ServiceError) => {
      console.error("❌ Slot stream error:", error.message);
    });

    slotStream.on("end", () => {
      console.log("ℹ️  Slot stream ended.");
    });
  }

  private async subscribeToWalletTransactions(walletAddresses: string[]): Promise<void> {
    console.log(`\n🔍 Subscribing to wallet transactions...`);
    console.log(`   Wallets: ${walletAddresses.join(", ")}`);

    const request = new SubscribeWalletRequest();
    request.setWalletAddressList(walletAddresses);

    const walletStream = this.client.subscribeToWalletTransactions(request, this.metadata);

    walletStream.on("data", (resp: any) => {
      try {
        const binaryData = resp.data || resp.getData?.() || resp;
        const msgWrapper = MessageWrapper.deserializeBinary(binaryData);

        const eventCase = msgWrapper.getEventMessageCase();

        if (eventCase === MessageWrapper.EventMessageCase.TRANSACTION) {
          const txWrapper = msgWrapper.getTransaction();
          if (txWrapper && txWrapper.hasTransaction()) {
            const transaction = txWrapper.getTransaction();
            if (transaction) {
              this.debugTransaction(transaction);
            }
          }
        } else {
          console.log(`⚠️ Received event case ${eventCase} instead of transaction`);
        }
      } catch (error) {
        console.error("❌ Failed to deserialize MessageWrapper:", error);
      }
    });

    walletStream.on("error", (error: grpc.ServiceError) => {
      console.error("❌ Wallet stream error:", error.message);
    });

    walletStream.on("end", () => {
      console.log("ℹ️  Wallet stream ended.");
    });
  }

  private async subscribeToAccountUpdates(
      accountAddresses: string[],
      ownerAddresses: string[]
  ): Promise<void> {
    console.log(`\n🔍 Subscribing to account updates...`);
    if (accountAddresses.length > 0) {
      console.log(`   Account addresses: ${accountAddresses.join(", ")}`);
    }
    if (ownerAddresses.length > 0) {
      console.log(`   Owner addresses: ${ownerAddresses.join(", ")}`);
    }

    const request = new SubscribeAccountsRequest();
    request.setAccountAddressList(accountAddresses);
    request.setOwnerAddressList(ownerAddresses);

    const accountStream = this.client.subscribeToAccountUpdates(request, this.metadata);

    console.log("⏳ Waiting for account updates...\n");

    accountStream.on("data", (resp: any) => {
      try {
        const binaryData = resp.data || resp.getData?.() || resp;
        const msgWrapper = MessageWrapper.deserializeBinary(binaryData);

        const eventCase = msgWrapper.getEventMessageCase();

        if (eventCase === MessageWrapper.EventMessageCase.ACCOUNT_UPDATE) {
          const accountUpdate = msgWrapper.getAccountUpdate();
          if (accountUpdate) {
            this.debugAccountUpdate(accountUpdate);
          }
        } else {
          console.log(`⚠️ Received event case ${eventCase} instead of account update`);
        }
      } catch (error) {
        console.error("❌ Failed to deserialize MessageWrapper:", error);
      }
    });

    accountStream.on("error", (error: grpc.ServiceError) => {
      console.error("❌ Account stream error:", error.message);
    });

    accountStream.on("end", () => {
      console.log("ℹ️  Account stream ended.");
    });
  }


  // --------------------------------------------------------------------------
  // Main Application Flow
  // --------------------------------------------------------------------------

  async run(): Promise<void> {
    console.log("╔════════════════════════════════════════╗");
    console.log("║     🚀 Thor Streamer Client v1.0       ║");
    console.log("╚════════════════════════════════════════╝");
    console.log(`\n📍 Server: ${this.config.serverAddress}`);
    console.log("────────────────────────────────────────\n");

    console.log("Select subscription method:");
    console.log("  1️⃣  All Transactions");
    console.log("  2️⃣  Slot Status Updates");
    console.log("  3️⃣  Wallet Transactions");
    console.log("  4️⃣  Account Updates");
    console.log("  0️⃣  Exit");

    const choice = await this.question("\n▶️  Enter your choice (0-4): ");

    switch (choice.trim()) {
      case "1":
        await this.subscribeToTransactions();
        break;

      case "2":
        await this.subscribeToSlots();
        break;

      case "3":
        const walletInput = await this.question("\n📝 Enter wallet addresses (comma-separated): ");
        const wallets = walletInput
            .split(",")
            .map(s => s.trim())
            .filter(s => s.length > 0);

        if (wallets.length === 0) {
          console.log("❌ No wallet addresses provided!");
          this.shutdown();
          return;
        }

        await this.subscribeToWalletTransactions(wallets);
        break;

      case "4":
        const accountInput = await this.question(
            "\n📝 Enter account addresses to monitor (comma-separated, or press Enter to skip): "
        );
        const accounts = accountInput.trim() === ""
            ? []
            : accountInput.split(",").map(s => s.trim()).filter(s => s.length > 0);

        const ownerInput = await this.question(
            "📝 Enter owner addresses to filter by (comma-separated, or press Enter to skip): "
        );
        const owners = ownerInput.trim() === ""
            ? []
            : ownerInput.split(",").map(s => s.trim()).filter(s => s.length > 0);

        if (accounts.length === 0 && owners.length === 0) {
          console.log("❌ At least one account address or owner address is required!");
          this.shutdown();
          return;
        }

        await this.subscribeToAccountUpdates(accounts, owners);
        break;

      case "0":
        console.log("\n👋 Goodbye!");
        this.shutdown();
        return;

      default:
        console.log("❌ Invalid choice!");
        this.shutdown();
        return;
    }

    // Setup graceful shutdown
    process.on("SIGINT", () => {
      console.log("\n\n🛑 Shutting down gracefully...");
      this.shutdown();
    });
  }

  private shutdown(): void {
    this.rl.close();
    process.exit(0);
  }
}

// ============================================================================
// Application Entry Point
// ============================================================================

async function main(): Promise<void> {
  try {
    const client = new ThorStreamerClient();
    await client.run();
  } catch (error) {
    console.error("💥 Fatal error:", error);
    process.exit(1);
  }
}

// Run the application
main();