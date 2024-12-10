import * as grpc from "@grpc/grpc-js";
import { EventPublisherClient } from "./protos/publisher_grpc_pb";
import { TransactionEvent, MessageWrapper, Message } from "./protos/events_pb";
import { Empty } from "google-protobuf/google/protobuf/empty_pb";
import base58 from "bs58";

// Load configuration
const config = {
  serverAddress: "ENDPOINT",
  authToken: "AUTH_TOKEN",
};

const client = new EventPublisherClient(
  config.serverAddress,
  grpc.credentials.createInsecure()
);

const metadata = new grpc.Metadata();
metadata.add("authorization", config.authToken);

const emptyRequest = new Empty();

console.log(`🔍 Starting Transaction Debugger on ${config.serverAddress}`);
console.log("--------------------------------");

const transactionStream = client.subscribeToTransactions(
  emptyRequest,
  metadata
);

transactionStream.on("data", (data: any) => {
  try {
    const binaryData = data.u[0];
    const msgWrapper = MessageWrapper.deserializeBinary(binaryData);

    const eventMessage = msgWrapper.getEventMessageCase();
    if (eventMessage === MessageWrapper.EventMessageCase.TRANSACTION) {
      const txWrapper = msgWrapper.getTransaction();
      if (txWrapper && txWrapper.getTransaction()) {
        const transaction = txWrapper.getTransaction();
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

function debugTransaction(tx: TransactionEvent) {
  console.log("\n🔍 Transaction Debug Information:");
  console.log(`├─ Signature: ${base58.encode(tx.getSignature_asU8())}`);
  console.log(`├─ Slot: ${tx.getSlot()}`);

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

  console.log(
    `├─ Version: ${message.getVersion()} (${getVersionString(
      message.getVersion()
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

function debugHeader(msg: any) {
  console.log("├─ Header:");
  const header = msg.getHeader();
  if (!header) {
    console.log("│  └─ ⚠️  Header is nil!");
    return;
  }
  console.log(
    `│  ├─ NumRequiredSignatures: ${header.getNumRequiredSignatures()}`
  );
  console.log(
    `│  ├─ NumReadonlySignedAccounts: ${header.getNumReadonlySignedAccounts()}`
  );
  console.log(
    `│  └─ NumReadonlyUnsignedAccounts: ${header.getNumReadonlyUnsignedAccounts()}`
  );
}

function debugAccountKeys(msg: Message) {
  const accountKeys = msg.getAccountKeysList_asU8();
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

function debugBlockhash(msg: any) {
  const blockhash = msg.getRecentBlockHash_asU8();
  console.log("├─ Recent Blockhash:");
  if (blockhash.length === 0) {
    console.log("│  └─ ⚠️  Blockhash is empty!");
    return;
  }
  console.log(`│  └─ ${base58.encode(blockhash)}`);
}

function debugInstructions(msg: any) {
  const instructions = msg.getInstructionsList();
  console.log(`─ Instructions (${instructions.length}):`);
  if (instructions.length === 0) {
    console.log("│  └─ ⚠️  No instructions!");
    return;
  }

  instructions.slice(0, 3).forEach((ix: any, i: number) => {
    console.log(`│  ├─ Instruction ${i}:`);
    console.log(`│  │  ├─ Program ID Index: ${ix.getProgramIdIndex()}`);
    console.log(`│  │  ├─ Account Indexes: ${ix.getAccountsList().length}`);
    console.log(`│  │  └─ Data Length: ${ix.getData_asU8().length} bytes`);
  });
  if (instructions.length > 3) {
    console.log(`│  └─ ... and ${instructions.length - 3} more instructions`);
  }
}
