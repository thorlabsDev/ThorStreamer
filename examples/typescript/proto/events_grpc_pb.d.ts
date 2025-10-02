// package: thor_streamer.types
// file: events.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as events_pb from "./events_pb";

interface IThorStreamerService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    streamUpdates: IThorStreamerService_IStreamUpdates;
}

interface IThorStreamerService_IStreamUpdates extends grpc.MethodDefinition<events_pb.Empty, events_pb.MessageWrapper> {
    path: "/thor_streamer.types.ThorStreamer/StreamUpdates";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<events_pb.Empty>;
    requestDeserialize: grpc.deserialize<events_pb.Empty>;
    responseSerialize: grpc.serialize<events_pb.MessageWrapper>;
    responseDeserialize: grpc.deserialize<events_pb.MessageWrapper>;
}

export const ThorStreamerService: IThorStreamerService;

export interface IThorStreamerServer {
    streamUpdates: grpc.handleServerStreamingCall<events_pb.Empty, events_pb.MessageWrapper>;
}

export interface IThorStreamerClient {
    streamUpdates(request: events_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<events_pb.MessageWrapper>;
    streamUpdates(request: events_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<events_pb.MessageWrapper>;
}

export class ThorStreamerClient extends grpc.Client implements IThorStreamerClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public streamUpdates(request: events_pb.Empty, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<events_pb.MessageWrapper>;
    public streamUpdates(request: events_pb.Empty, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<events_pb.MessageWrapper>;
}
