// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var publisher_pb = require('./publisher_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_publisher_StreamResponse(arg) {
  if (!(arg instanceof publisher_pb.StreamResponse)) {
    throw new Error('Expected argument of type publisher.StreamResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_publisher_StreamResponse(buffer_arg) {
  return publisher_pb.StreamResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_publisher_SubscribeAccountsRequest(arg) {
  if (!(arg instanceof publisher_pb.SubscribeAccountsRequest)) {
    throw new Error('Expected argument of type publisher.SubscribeAccountsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_publisher_SubscribeAccountsRequest(buffer_arg) {
  return publisher_pb.SubscribeAccountsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_publisher_SubscribeWalletRequest(arg) {
  if (!(arg instanceof publisher_pb.SubscribeWalletRequest)) {
    throw new Error('Expected argument of type publisher.SubscribeWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_publisher_SubscribeWalletRequest(buffer_arg) {
  return publisher_pb.SubscribeWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var EventPublisherService = exports.EventPublisherService = {
  subscribeToTransactions: {
    path: '/publisher.EventPublisher/SubscribeToTransactions',
    requestStream: false,
    responseStream: true,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: publisher_pb.StreamResponse,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_publisher_StreamResponse,
    responseDeserialize: deserialize_publisher_StreamResponse,
  },
  subscribeToSlotStatus: {
    path: '/publisher.EventPublisher/SubscribeToSlotStatus',
    requestStream: false,
    responseStream: true,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: publisher_pb.StreamResponse,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_publisher_StreamResponse,
    responseDeserialize: deserialize_publisher_StreamResponse,
  },
  subscribeToWalletTransactions: {
    path: '/publisher.EventPublisher/SubscribeToWalletTransactions',
    requestStream: false,
    responseStream: true,
    requestType: publisher_pb.SubscribeWalletRequest,
    responseType: publisher_pb.StreamResponse,
    requestSerialize: serialize_publisher_SubscribeWalletRequest,
    requestDeserialize: deserialize_publisher_SubscribeWalletRequest,
    responseSerialize: serialize_publisher_StreamResponse,
    responseDeserialize: deserialize_publisher_StreamResponse,
  },
  // NEW method to subscribe to account updates
subscribeToAccountUpdates: {
    path: '/publisher.EventPublisher/SubscribeToAccountUpdates',
    requestStream: false,
    responseStream: true,
    requestType: publisher_pb.SubscribeAccountsRequest,
    responseType: publisher_pb.StreamResponse,
    requestSerialize: serialize_publisher_SubscribeAccountsRequest,
    requestDeserialize: deserialize_publisher_SubscribeAccountsRequest,
    responseSerialize: serialize_publisher_StreamResponse,
    responseDeserialize: deserialize_publisher_StreamResponse,
  },
};

exports.EventPublisherClient = grpc.makeGenericClientConstructor(EventPublisherService);
