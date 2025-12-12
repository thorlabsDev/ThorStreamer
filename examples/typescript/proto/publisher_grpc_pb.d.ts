// package: publisher
// file: publisher.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as publisher_pb from "./publisher_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

interface IEventPublisherService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    subscribeToTransactions: IEventPublisherService_ISubscribeToTransactions;
    subscribeToSlotStatus: IEventPublisherService_ISubscribeToSlotStatus;
    subscribeToAccountUpdates: IEventPublisherService_ISubscribeToAccountUpdates;
}

interface IEventPublisherService_ISubscribeToTransactions extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, publisher_pb.StreamResponse> {
    path: "/publisher.EventPublisher/SubscribeToTransactions";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<publisher_pb.StreamResponse>;
    responseDeserialize: grpc.deserialize<publisher_pb.StreamResponse>;
}
interface IEventPublisherService_ISubscribeToSlotStatus extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, publisher_pb.StreamResponse> {
    path: "/publisher.EventPublisher/SubscribeToSlotStatus";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<publisher_pb.StreamResponse>;
    responseDeserialize: grpc.deserialize<publisher_pb.StreamResponse>;
}
interface IEventPublisherService_ISubscribeToAccountUpdates extends grpc.MethodDefinition<publisher_pb.SubscribeAccountsRequest, publisher_pb.StreamResponse> {
    path: "/publisher.EventPublisher/SubscribeToAccountUpdates";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<publisher_pb.SubscribeAccountsRequest>;
    requestDeserialize: grpc.deserialize<publisher_pb.SubscribeAccountsRequest>;
    responseSerialize: grpc.serialize<publisher_pb.StreamResponse>;
    responseDeserialize: grpc.deserialize<publisher_pb.StreamResponse>;
}

export const EventPublisherService: IEventPublisherService;

export interface IEventPublisherServer {
    subscribeToTransactions: grpc.handleServerStreamingCall<google_protobuf_empty_pb.Empty, publisher_pb.StreamResponse>;
    subscribeToSlotStatus: grpc.handleServerStreamingCall<google_protobuf_empty_pb.Empty, publisher_pb.StreamResponse>;
    subscribeToAccountUpdates: grpc.handleServerStreamingCall<publisher_pb.SubscribeAccountsRequest, publisher_pb.StreamResponse>;
}

export interface IEventPublisherClient {
    subscribeToTransactions(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    subscribeToTransactions(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    subscribeToSlotStatus(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    subscribeToSlotStatus(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    subscribeToAccountUpdates(request: publisher_pb.SubscribeAccountsRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    subscribeToAccountUpdates(request: publisher_pb.SubscribeAccountsRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
}

export class EventPublisherClient extends grpc.Client implements IEventPublisherClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public subscribeToTransactions(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    public subscribeToTransactions(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    public subscribeToSlotStatus(request: google_protobuf_empty_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    public subscribeToSlotStatus(request: google_protobuf_empty_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    public subscribeToAccountUpdates(request: publisher_pb.SubscribeAccountsRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
    public subscribeToAccountUpdates(request: publisher_pb.SubscribeAccountsRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<publisher_pb.StreamResponse>;
}
