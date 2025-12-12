# TypeScript SDK

TypeScript/Node.js client for ThorStreamer using gRPC-js.

## Installation

```bash
npm install @grpc/grpc-js google-protobuf
```

## Proto Generation

Generate TypeScript definitions from proto files:

```bash
# Clone the repository
git clone https://github.com/thorlabsDev/ThorStreamer.git
cd ThorStreamer/examples/typescript

# Install dependencies
npm install

# Generate proto files
npm run proto:gen
npm run proto:ts
```

## Quick Start

```typescript
import * as grpc from '@grpc/grpc-js';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { EventPublisherClient } from './proto/publisher_grpc_pb';
import { StreamResponse } from './proto/publisher_pb';

const client = new EventPublisherClient(
    process.env.SERVER_ADDRESS!,
    grpc.credentials.createInsecure()
);

const metadata = new grpc.Metadata();
metadata.set('authorization', process.env.AUTH_TOKEN!);

const stream = client.subscribeToTransactions(new Empty(), metadata);

stream.on('data', (response: StreamResponse) => {
    console.log('Transaction:', response.getData().length, 'bytes');
});

stream.on('error', (err) => {
    console.error('Stream error:', err);
});

stream.on('end', () => {
    console.log('Stream ended');
});
```

## API Reference

### Creating a Client

```typescript
import * as grpc from '@grpc/grpc-js';
import { EventPublisherClient } from './proto/publisher_grpc_pb';

const client = new EventPublisherClient(
    'your-server:50051',
    grpc.credentials.createInsecure()
);

const metadata = new grpc.Metadata();
metadata.set('authorization', 'your-token');
```

### subscribeToTransactions

```typescript
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

const stream = client.subscribeToTransactions(new Empty(), metadata);

stream.on('data', (response: StreamResponse) => {
    const data = response.getData_asU8();
    // Deserialize transaction data
    console.log('Transaction received:', data.length, 'bytes');
});
```

### subscribeToSlotStatus

```typescript
const stream = client.subscribeToSlotStatus(new Empty(), metadata);

stream.on('data', (response: StreamResponse) => {
    const data = response.getData_asU8();
    // Deserialize slot status
    console.log('Slot update:', data.length, 'bytes');
});
```

### subscribeToAccountUpdates

Monitor accounts with optional owner filtering:

```typescript
import { SubscribeAccountsRequest } from './proto/publisher_pb';

const request = new SubscribeAccountsRequest();
request.setAccountAddressList(['account1...', 'account2...']);
request.setOwnerAddressList(['TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA']);

const stream = client.subscribeToAccountUpdates(request, metadata);

stream.on('data', (response: StreamResponse) => {
    console.log('Account update received');
});
```

## Error Handling

```typescript
stream.on('error', (err: grpc.ServiceError) => {
    console.error('Error code:', err.code);
    console.error('Error message:', err.message);

    // Implement reconnection logic
    if (err.code === grpc.status.UNAVAILABLE) {
        // Server unavailable, retry with backoff
    }
});

stream.on('end', () => {
    console.log('Stream ended normally');
});
```

## Promise Wrapper

Wrap streams in Promises for async/await usage:

```typescript
async function subscribeTransactions(): Promise<void> {
    return new Promise((resolve, reject) => {
        const stream = client.subscribeToTransactions(new Empty(), metadata);

        stream.on('data', (response: StreamResponse) => {
            processTransaction(response.getData_asU8());
        });

        stream.on('error', reject);
        stream.on('end', resolve);
    });
}

// Usage
try {
    await subscribeTransactions();
} catch (err) {
    console.error('Stream failed:', err);
}
```

## Full Example

See [examples/typescript](https://github.com/thorlabsDev/ThorStreamer/tree/master/examples/typescript) for a complete implementation with proto generation scripts.

## Build Scripts

```json
{
  "scripts": {
    "proto:gen": "Generate JavaScript from protos",
    "proto:ts": "Generate TypeScript definitions",
    "proto:clean": "Remove generated files",
    "build": "tsc",
    "start": "node dist/client.js",
    "dev": "ts-node src/client.ts"
  }
}
```

## Resources

- [GitHub Repository](https://github.com/thorlabsDev/ThorStreamer)
- [@grpc/grpc-js Documentation](https://www.npmjs.com/package/@grpc/grpc-js)
