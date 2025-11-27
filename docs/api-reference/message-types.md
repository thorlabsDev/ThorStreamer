# Message Types

## Response Wrappers

### StreamResponse
Standard response for EventPublisher methods.

```protobuf
message StreamResponse {
  bytes data = 1;  // Serialized event data
}
```

### MessageWrapper
Unified wrapper for ThorStreamer service containing typed events.

```protobuf
message MessageWrapper {
  oneof event_message {
    SubscribeUpdateAccountInfo account_update = 1;
    SlotStatusEvent slot = 2;
    TransactionEventWrapper transaction = 3;
  }
}
```

---

## Transaction Events

### TransactionEventWrapper

```protobuf
message TransactionEventWrapper {
  StreamType stream_type = 1;
  TransactionEvent transaction = 2;
}

enum StreamType {
  STREAM_TYPE_UNSPECIFIED = 0;
  STREAM_TYPE_FILTERED = 1;    // From program filters
  STREAM_TYPE_WALLET = 2;      // From wallet subscription
  STREAM_TYPE_ACCOUNT = 3;     // From account subscription
}
```

### TransactionEvent

```protobuf
message TransactionEvent {
  uint64 slot = 1;                    // Slot number
  bytes signature = 2;                 // 64-byte signature
  uint64 index = 3;                   // Index within block
  bool is_vote = 4;                   // Vote transaction flag
  SanitizedTransaction transaction = 5;
  TransactionStatusMeta transaction_status_meta = 6;
}
```

### SanitizedTransaction

```protobuf
message SanitizedTransaction {
  Message message = 1;
  bytes message_hash = 2;              // 32 bytes
  repeated bytes signatures = 3;       // 64-byte signatures
  bool is_simple_vote_transaction = 4;
}
```

### TransactionStatusMeta

```protobuf
message TransactionStatusMeta {
  bool is_status_err = 1;              // Success/failure
  uint64 fee = 2;                      // Fee in lamports
  repeated uint64 pre_balances = 3;
  repeated uint64 post_balances = 4;
  repeated InnerInstructions inner_instructions = 5;
  repeated string log_messages = 6;
  repeated TransactionTokenBalance pre_token_balances = 7;
  repeated TransactionTokenBalance post_token_balances = 8;
  repeated Reward rewards = 9;
  string error_info = 10;
}
```

---

## Account Events

### SubscribeUpdateAccountInfo

```protobuf
message SubscribeUpdateAccountInfo {
  bytes pubkey = 1;              // Account public key
  uint64 lamports = 2;           // Balance in lamports
  bytes owner = 3;               // Owner program
  bool executable = 4;           // Executable flag
  uint64 rent_epoch = 5;
  bytes data = 6;                // Account data
  uint64 write_version = 7;
  optional bytes txn_signature = 8;
  optional SlotStatus slot = 9;
}
```

---

## Slot Events

### SlotStatusEvent

```protobuf
message SlotStatusEvent {
  uint64 slot = 1;
  uint64 parent = 2;
  int32 status = 3;              // 0=processed, 1=confirmed, 2=rooted
  bytes block_hash = 4;          // 32 bytes
  uint64 block_height = 5;
}
```

---

## Supporting Types

### Message (Transaction)

```protobuf
message Message {
  uint32 version = 1;                  // 0=legacy, 1=v0
  MessageHeader header = 2;
  bytes recent_block_hash = 3;         // 32 bytes
  repeated bytes account_keys = 4;     // 32-byte keys
  repeated CompiledInstruction instructions = 5;
  repeated MessageAddressTableLookup address_table_lookups = 6;  // v0 only
  LoadedAddresses loaded_addresses = 7;                          // v0 only
  repeated bool is_writable = 8;
}
```

### CompiledInstruction

```protobuf
message CompiledInstruction {
  uint32 program_id_index = 1;
  bytes data = 2;
  repeated uint32 accounts = 3;
}
```

### TransactionTokenBalance

```protobuf
message TransactionTokenBalance {
  uint32 account_index = 1;
  string mint = 2;
  UiTokenAmount ui_token_amount = 3;
  string owner = 4;
}

message UiTokenAmount {
  double ui_amount = 1;
  uint32 decimals = 2;
  string amount = 3;
  string ui_amount_string = 4;
}
```
