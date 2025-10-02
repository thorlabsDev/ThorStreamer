// proto/index.ts
// This file consolidates all proto exports for cleaner imports

// Import generated files
import * as events_pb from './events_pb';
import * as publisher_pb from './publisher_pb';
import * as publisher_grpc_pb from './publisher_grpc_pb';

// Import Empty from google-protobuf package
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

// Re-export everything
export * from './events_pb';
export * from './publisher_pb';
export * from './publisher_grpc_pb';

// Named exports for convenience
export const EventPublisherClient = publisher_grpc_pb.EventPublisherClient;
export const SubscribeWalletRequest = publisher_pb.SubscribeWalletRequest;
export const SubscribeAccountsRequest = publisher_pb.SubscribeAccountsRequest;
export const StreamResponse = publisher_pb.StreamResponse;

export const MessageWrapper = events_pb.MessageWrapper;
export const TransactionEvent = events_pb.TransactionEvent;
export const SlotStatusEvent = events_pb.SlotStatusEvent;
export const SubscribeUpdateAccountInfo = events_pb.SubscribeUpdateAccountInfo;
export const Message = events_pb.Message;
export const MessageHeader = events_pb.MessageHeader;
export const CompiledInstruction = events_pb.CompiledInstruction;
export const TransactionEventWrapper = events_pb.TransactionEventWrapper;
export const SlotStatus = events_pb.SlotStatus;
export const SanitizedTransaction = events_pb.SanitizedTransaction;
export const TransactionStatusMeta = events_pb.TransactionStatusMeta;

// Export the correct Empty from google-protobuf
export { Empty };

// Type aliases for better TypeScript support
export type IEventPublisherClient = InstanceType<typeof EventPublisherClient>;
export type IMessageWrapper = events_pb.MessageWrapper;
export type ITransactionEvent = events_pb.TransactionEvent;
export type ISlotStatusEvent = events_pb.SlotStatusEvent;
export type ISubscribeUpdateAccountInfo = events_pb.SubscribeUpdateAccountInfo;
export type IMessage = events_pb.Message;