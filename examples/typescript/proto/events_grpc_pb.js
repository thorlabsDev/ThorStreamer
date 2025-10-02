// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var events_pb = require('./events_pb.js');

function serialize_thor_streamer_types_Empty(arg) {
  if (!(arg instanceof events_pb.Empty)) {
    throw new Error('Expected argument of type thor_streamer.types.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_thor_streamer_types_Empty(buffer_arg) {
  return events_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_thor_streamer_types_MessageWrapper(arg) {
  if (!(arg instanceof events_pb.MessageWrapper)) {
    throw new Error('Expected argument of type thor_streamer.types.MessageWrapper');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_thor_streamer_types_MessageWrapper(buffer_arg) {
  return events_pb.MessageWrapper.deserializeBinary(new Uint8Array(buffer_arg));
}


var ThorStreamerService = exports.ThorStreamerService = {
  streamUpdates: {
    path: '/thor_streamer.types.ThorStreamer/StreamUpdates',
    requestStream: false,
    responseStream: true,
    requestType: events_pb.Empty,
    responseType: events_pb.MessageWrapper,
    requestSerialize: serialize_thor_streamer_types_Empty,
    requestDeserialize: deserialize_thor_streamer_types_Empty,
    responseSerialize: serialize_thor_streamer_types_MessageWrapper,
    responseDeserialize: deserialize_thor_streamer_types_MessageWrapper,
  },
};

exports.ThorStreamerClient = grpc.makeGenericClientConstructor(ThorStreamerService);
