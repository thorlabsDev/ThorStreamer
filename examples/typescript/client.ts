Pimport * as grpc from "@grpc/grpc-js";
import { publisher } from "../../protos/publisher";
import { thor_streamer } from "../../protos/events";
import * as goog from "../../protos/google/protobuf/empty";
import base58 from "bs58";
// Load configuration
const config = {
  serverAddress: "ENDPOINT",
  authToken: "AUTH_TOKEN",
};

const client = new publisher.EventPublisherClient(
  config.serverAddress,
  grpc.credentials.createInsecure()
);

const metadata = new grpc.Metadata();
metadata.add("authorization", config.authToken);

const emptyRequest = new goog.google.protobuf.Empty;

console.log(`🔍 Starting Transaction Debugger on ${config.serverAddress}`);
console.log("--------------------------------");

const transactionStream = client.SubscribeToTransactions(
  emptyRequest, metadata
);

transactionStream.on("data", (data: any) => {
  try {
    const binaryData = data.u[0];
    const msgWrapper = thor_streamer.types.MessageWrapper.deserializeBinary(binaryData);

    const eventMessage = msgWrapper.event_message;
    if (eventMessage === "transaction") {
      const txWrapper = msgWrapper.transaction;
      if (txWrapper && txWrapper.transaction) {
        const transaction = txWrapper.transaction;
        debugTransaction(transaction!);
      } else {
        console.log("No actual transaction found in TransactionEventWrapper");
      }
    } else {
      console.log("Received a message that is not a transaction event");
    }
  } catch (error) {
    console.error("Failed to deserialize MessageWrapper:", error);
  }
});

transactionStream.on("error", (error: grpc.ServiceError) => {
  console.error("Transaction stream error:", error);
});

transactionStream.on("end", () => {
  console.log("Transaction stream ended.");
});

function debugTransaction(tx: thor_streamer.types.TransactionEvent) {
  console.log("\n🔍 Transaction Debug Information:");
  console.log(`├─ Signature: ${base58.encode(tx.signature)}`);
  console.log(`├─ Slot: ${tx.slot}`);

  const transaction = tx.transaction;

  if (!transaction) {
    console.log("├─ ⚠️  Transaction is nil!");
    return;
  }

  const message = transaction.message;
  if (!message) {
    console.log("├─ ⚠️  Message is nil!");
    return;
  }

  console.log(
    `├─ Version: ${message.version} (${getVersionString(
      message.version
    )})`
  );
  debugHeader(message);
  debugAccountKeys(message);
  debugBlockhash(message);
  debugInstructions(message);

}

function getVersionString(version: number): string {
  switch (version) {
    case 0:
      return "Legacy";
    case 1:
      return "V0";
    default:
      return `Unknown(${version})`;
  }
}

function debugHeader(msg: thor_streamer.types.Message) {
  console.log("├─ Header:");
  const header = msg.header
  if (!header) {
    console.log("│  └─ ⚠️  Header is nil!");
    return;
  }
  console.log(
    `│  ├─ NumRequiredSignatures: ${header.num_required_signatures}`
  );
  console.log(
    `│  ├─ NumReadonlySignedAccounts: ${header.num_readonly_signed_accounts}`
  );
  console.log(
    `│  └─ NumReadonlyUnsignedAccounts: ${header.num_readonly_unsigned_accounts}`
  );
}

function debugAccountKeys(msg: thor_streamer.types.Message) {
  const accountKeys = msg.account_keys;
  console.log(`├─ Account Keys (${accountKeys.length}):`);
  if (accountKeys.length === 0) {
    console.log("│  └─ ⚠️  No account keys!");
    return;
  }

  accountKeys.slice(0, 5).forEach((key: Uint8Array, i: number) => {
    console.log(`│  ├─ [${i}]: ${base58.encode(key)}`);
  });
  if (accountKeys.length > 5) {
    console.log(`│  └─ ... and ${accountKeys.length - 5} more keys`);
  }
}

function debugBlockhash(msg: thor_streamer.types.Message) {
  const blockhash = msg.recent_block_hash;
  console.log("├─ Recent Blockhash:");
  if (blockhash.length === 0) {
    console.log("│  └─ ⚠️  Blockhash is empty!");
    return;
  }
  console.log(`│  └─ ${base58.encode(blockhash)}`);
}

function debugInstructions(msg: thor_streamer.types.Message) {
  const instructions = msg.instructions;
  console.log(`─ Instructions (${instructions.length}):`);
  if (instructions.length === 0) {
    console.log("│  └─ ⚠️  No instructions!");
    return;
  }

  instructions.slice(0, 3).forEach((ix: thor_streamer.types.CompiledInstruction, i: number) => {
    console.log(`│  ├─ Instruction ${i}:`);
    console.log(`│  │  ├─ Program ID Index: ${ix.program_id_index}`);
    console.log(`│  │  ├─ Account Indexes: ${ix.accounts.length}`);
    console.log(`│  │  └─ Data Length: ${ix.data.length} bytes`);
  });
  if (instructions.length > 3) {
    console.log(`│  └─ ... and ${instructions.length - 3} more instructions`);
  }
}
