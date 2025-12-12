// package: publisher
// file: publisher.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

export class StreamResponse extends jspb.Message { 
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): StreamResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StreamResponse.AsObject;
    static toObject(includeInstance: boolean, msg: StreamResponse): StreamResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StreamResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StreamResponse;
    static deserializeBinaryFromReader(message: StreamResponse, reader: jspb.BinaryReader): StreamResponse;
}

export namespace StreamResponse {
    export type AsObject = {
        data: Uint8Array | string,
    }
}

export class SubscribeAccountsRequest extends jspb.Message { 
    clearAccountAddressList(): void;
    getAccountAddressList(): Array<string>;
    setAccountAddressList(value: Array<string>): SubscribeAccountsRequest;
    addAccountAddress(value: string, index?: number): string;
    clearOwnerAddressList(): void;
    getOwnerAddressList(): Array<string>;
    setOwnerAddressList(value: Array<string>): SubscribeAccountsRequest;
    addOwnerAddress(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubscribeAccountsRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SubscribeAccountsRequest): SubscribeAccountsRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubscribeAccountsRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubscribeAccountsRequest;
    static deserializeBinaryFromReader(message: SubscribeAccountsRequest, reader: jspb.BinaryReader): SubscribeAccountsRequest;
}

export namespace SubscribeAccountsRequest {
    export type AsObject = {
        accountAddressList: Array<string>,
        ownerAddressList: Array<string>,
    }
}
