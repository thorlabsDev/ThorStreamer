// package: thor_streamer.types
// file: events.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class Empty extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Empty.AsObject;
    static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Empty;
    static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
    export type AsObject = {
    }
}

export class SlotStatus extends jspb.Message { 
    getSlot(): number;
    setSlot(value: number): SlotStatus;
    getParent(): number;
    setParent(value: number): SlotStatus;
    getStatus(): number;
    setStatus(value: number): SlotStatus;
    getBlockHash(): Uint8Array | string;
    getBlockHash_asU8(): Uint8Array;
    getBlockHash_asB64(): string;
    setBlockHash(value: Uint8Array | string): SlotStatus;
    getBlockHeight(): number;
    setBlockHeight(value: number): SlotStatus;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SlotStatus.AsObject;
    static toObject(includeInstance: boolean, msg: SlotStatus): SlotStatus.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SlotStatus, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SlotStatus;
    static deserializeBinaryFromReader(message: SlotStatus, reader: jspb.BinaryReader): SlotStatus;
}

export namespace SlotStatus {
    export type AsObject = {
        slot: number,
        parent: number,
        status: number,
        blockHash: Uint8Array | string,
        blockHeight: number,
    }
}

export class SubscribeUpdateAccountInfo extends jspb.Message { 
    getPubkey(): Uint8Array | string;
    getPubkey_asU8(): Uint8Array;
    getPubkey_asB64(): string;
    setPubkey(value: Uint8Array | string): SubscribeUpdateAccountInfo;
    getLamports(): number;
    setLamports(value: number): SubscribeUpdateAccountInfo;
    getOwner(): Uint8Array | string;
    getOwner_asU8(): Uint8Array;
    getOwner_asB64(): string;
    setOwner(value: Uint8Array | string): SubscribeUpdateAccountInfo;
    getExecutable(): boolean;
    setExecutable(value: boolean): SubscribeUpdateAccountInfo;
    getRentEpoch(): number;
    setRentEpoch(value: number): SubscribeUpdateAccountInfo;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): SubscribeUpdateAccountInfo;
    getWriteVersion(): number;
    setWriteVersion(value: number): SubscribeUpdateAccountInfo;

    hasTxnSignature(): boolean;
    clearTxnSignature(): void;
    getTxnSignature(): Uint8Array | string;
    getTxnSignature_asU8(): Uint8Array;
    getTxnSignature_asB64(): string;
    setTxnSignature(value: Uint8Array | string): SubscribeUpdateAccountInfo;

    hasSlot(): boolean;
    clearSlot(): void;
    getSlot(): SlotStatus | undefined;
    setSlot(value?: SlotStatus): SubscribeUpdateAccountInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubscribeUpdateAccountInfo.AsObject;
    static toObject(includeInstance: boolean, msg: SubscribeUpdateAccountInfo): SubscribeUpdateAccountInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubscribeUpdateAccountInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubscribeUpdateAccountInfo;
    static deserializeBinaryFromReader(message: SubscribeUpdateAccountInfo, reader: jspb.BinaryReader): SubscribeUpdateAccountInfo;
}

export namespace SubscribeUpdateAccountInfo {
    export type AsObject = {
        pubkey: Uint8Array | string,
        lamports: number,
        owner: Uint8Array | string,
        executable: boolean,
        rentEpoch: number,
        data: Uint8Array | string,
        writeVersion: number,
        txnSignature: Uint8Array | string,
        slot?: SlotStatus.AsObject,
    }
}

export class ThorAccountsRequest extends jspb.Message { 
    clearAccountAddressList(): void;
    getAccountAddressList(): Array<string>;
    setAccountAddressList(value: Array<string>): ThorAccountsRequest;
    addAccountAddress(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ThorAccountsRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ThorAccountsRequest): ThorAccountsRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ThorAccountsRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ThorAccountsRequest;
    static deserializeBinaryFromReader(message: ThorAccountsRequest, reader: jspb.BinaryReader): ThorAccountsRequest;
}

export namespace ThorAccountsRequest {
    export type AsObject = {
        accountAddressList: Array<string>,
    }
}

export class SlotStatusEvent extends jspb.Message { 
    getSlot(): number;
    setSlot(value: number): SlotStatusEvent;
    getParent(): number;
    setParent(value: number): SlotStatusEvent;
    getStatus(): number;
    setStatus(value: number): SlotStatusEvent;
    getBlockHash(): Uint8Array | string;
    getBlockHash_asU8(): Uint8Array;
    getBlockHash_asB64(): string;
    setBlockHash(value: Uint8Array | string): SlotStatusEvent;
    getBlockHeight(): number;
    setBlockHeight(value: number): SlotStatusEvent;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SlotStatusEvent.AsObject;
    static toObject(includeInstance: boolean, msg: SlotStatusEvent): SlotStatusEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SlotStatusEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SlotStatusEvent;
    static deserializeBinaryFromReader(message: SlotStatusEvent, reader: jspb.BinaryReader): SlotStatusEvent;
}

export namespace SlotStatusEvent {
    export type AsObject = {
        slot: number,
        parent: number,
        status: number,
        blockHash: Uint8Array | string,
        blockHeight: number,
    }
}

export class MessageHeader extends jspb.Message { 
    getNumRequiredSignatures(): number;
    setNumRequiredSignatures(value: number): MessageHeader;
    getNumReadonlySignedAccounts(): number;
    setNumReadonlySignedAccounts(value: number): MessageHeader;
    getNumReadonlyUnsignedAccounts(): number;
    setNumReadonlyUnsignedAccounts(value: number): MessageHeader;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageHeader.AsObject;
    static toObject(includeInstance: boolean, msg: MessageHeader): MessageHeader.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageHeader, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageHeader;
    static deserializeBinaryFromReader(message: MessageHeader, reader: jspb.BinaryReader): MessageHeader;
}

export namespace MessageHeader {
    export type AsObject = {
        numRequiredSignatures: number,
        numReadonlySignedAccounts: number,
        numReadonlyUnsignedAccounts: number,
    }
}

export class CompiledInstruction extends jspb.Message { 
    getProgramIdIndex(): number;
    setProgramIdIndex(value: number): CompiledInstruction;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): CompiledInstruction;
    clearAccountsList(): void;
    getAccountsList(): Array<number>;
    setAccountsList(value: Array<number>): CompiledInstruction;
    addAccounts(value: number, index?: number): number;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CompiledInstruction.AsObject;
    static toObject(includeInstance: boolean, msg: CompiledInstruction): CompiledInstruction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CompiledInstruction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CompiledInstruction;
    static deserializeBinaryFromReader(message: CompiledInstruction, reader: jspb.BinaryReader): CompiledInstruction;
}

export namespace CompiledInstruction {
    export type AsObject = {
        programIdIndex: number,
        data: Uint8Array | string,
        accountsList: Array<number>,
    }
}

export class LoadedAddresses extends jspb.Message { 
    clearWritableList(): void;
    getWritableList(): Array<Uint8Array | string>;
    getWritableList_asU8(): Array<Uint8Array>;
    getWritableList_asB64(): Array<string>;
    setWritableList(value: Array<Uint8Array | string>): LoadedAddresses;
    addWritable(value: Uint8Array | string, index?: number): Uint8Array | string;
    clearReadonlyList(): void;
    getReadonlyList(): Array<Uint8Array | string>;
    getReadonlyList_asU8(): Array<Uint8Array>;
    getReadonlyList_asB64(): Array<string>;
    setReadonlyList(value: Array<Uint8Array | string>): LoadedAddresses;
    addReadonly(value: Uint8Array | string, index?: number): Uint8Array | string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LoadedAddresses.AsObject;
    static toObject(includeInstance: boolean, msg: LoadedAddresses): LoadedAddresses.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LoadedAddresses, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LoadedAddresses;
    static deserializeBinaryFromReader(message: LoadedAddresses, reader: jspb.BinaryReader): LoadedAddresses;
}

export namespace LoadedAddresses {
    export type AsObject = {
        writableList: Array<Uint8Array | string>,
        readonlyList: Array<Uint8Array | string>,
    }
}

export class MessageAddressTableLookup extends jspb.Message { 
    getAccountKey(): Uint8Array | string;
    getAccountKey_asU8(): Uint8Array;
    getAccountKey_asB64(): string;
    setAccountKey(value: Uint8Array | string): MessageAddressTableLookup;
    getWritableIndexes(): Uint8Array | string;
    getWritableIndexes_asU8(): Uint8Array;
    getWritableIndexes_asB64(): string;
    setWritableIndexes(value: Uint8Array | string): MessageAddressTableLookup;
    getReadonlyIndexes(): Uint8Array | string;
    getReadonlyIndexes_asU8(): Uint8Array;
    getReadonlyIndexes_asB64(): string;
    setReadonlyIndexes(value: Uint8Array | string): MessageAddressTableLookup;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageAddressTableLookup.AsObject;
    static toObject(includeInstance: boolean, msg: MessageAddressTableLookup): MessageAddressTableLookup.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageAddressTableLookup, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageAddressTableLookup;
    static deserializeBinaryFromReader(message: MessageAddressTableLookup, reader: jspb.BinaryReader): MessageAddressTableLookup;
}

export namespace MessageAddressTableLookup {
    export type AsObject = {
        accountKey: Uint8Array | string,
        writableIndexes: Uint8Array | string,
        readonlyIndexes: Uint8Array | string,
    }
}

export class Message extends jspb.Message { 
    getVersion(): number;
    setVersion(value: number): Message;

    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): MessageHeader | undefined;
    setHeader(value?: MessageHeader): Message;
    getRecentBlockHash(): Uint8Array | string;
    getRecentBlockHash_asU8(): Uint8Array;
    getRecentBlockHash_asB64(): string;
    setRecentBlockHash(value: Uint8Array | string): Message;
    clearAccountKeysList(): void;
    getAccountKeysList(): Array<Uint8Array | string>;
    getAccountKeysList_asU8(): Array<Uint8Array>;
    getAccountKeysList_asB64(): Array<string>;
    setAccountKeysList(value: Array<Uint8Array | string>): Message;
    addAccountKeys(value: Uint8Array | string, index?: number): Uint8Array | string;
    clearInstructionsList(): void;
    getInstructionsList(): Array<CompiledInstruction>;
    setInstructionsList(value: Array<CompiledInstruction>): Message;
    addInstructions(value?: CompiledInstruction, index?: number): CompiledInstruction;
    clearAddressTableLookupsList(): void;
    getAddressTableLookupsList(): Array<MessageAddressTableLookup>;
    setAddressTableLookupsList(value: Array<MessageAddressTableLookup>): Message;
    addAddressTableLookups(value?: MessageAddressTableLookup, index?: number): MessageAddressTableLookup;

    hasLoadedAddresses(): boolean;
    clearLoadedAddresses(): void;
    getLoadedAddresses(): LoadedAddresses | undefined;
    setLoadedAddresses(value?: LoadedAddresses): Message;
    clearIsWritableList(): void;
    getIsWritableList(): Array<boolean>;
    setIsWritableList(value: Array<boolean>): Message;
    addIsWritable(value: boolean, index?: number): boolean;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Message.AsObject;
    static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Message;
    static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
    export type AsObject = {
        version: number,
        header?: MessageHeader.AsObject,
        recentBlockHash: Uint8Array | string,
        accountKeysList: Array<Uint8Array | string>,
        instructionsList: Array<CompiledInstruction.AsObject>,
        addressTableLookupsList: Array<MessageAddressTableLookup.AsObject>,
        loadedAddresses?: LoadedAddresses.AsObject,
        isWritableList: Array<boolean>,
    }
}

export class SanitizedTransaction extends jspb.Message { 

    hasMessage(): boolean;
    clearMessage(): void;
    getMessage(): Message | undefined;
    setMessage(value?: Message): SanitizedTransaction;
    getMessageHash(): Uint8Array | string;
    getMessageHash_asU8(): Uint8Array;
    getMessageHash_asB64(): string;
    setMessageHash(value: Uint8Array | string): SanitizedTransaction;
    clearSignaturesList(): void;
    getSignaturesList(): Array<Uint8Array | string>;
    getSignaturesList_asU8(): Array<Uint8Array>;
    getSignaturesList_asB64(): Array<string>;
    setSignaturesList(value: Array<Uint8Array | string>): SanitizedTransaction;
    addSignatures(value: Uint8Array | string, index?: number): Uint8Array | string;
    getIsSimpleVoteTransaction(): boolean;
    setIsSimpleVoteTransaction(value: boolean): SanitizedTransaction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SanitizedTransaction.AsObject;
    static toObject(includeInstance: boolean, msg: SanitizedTransaction): SanitizedTransaction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SanitizedTransaction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SanitizedTransaction;
    static deserializeBinaryFromReader(message: SanitizedTransaction, reader: jspb.BinaryReader): SanitizedTransaction;
}

export namespace SanitizedTransaction {
    export type AsObject = {
        message?: Message.AsObject,
        messageHash: Uint8Array | string,
        signaturesList: Array<Uint8Array | string>,
        isSimpleVoteTransaction: boolean,
    }
}

export class TransactionEvent extends jspb.Message { 
    getSlot(): number;
    setSlot(value: number): TransactionEvent;
    getSignature(): Uint8Array | string;
    getSignature_asU8(): Uint8Array;
    getSignature_asB64(): string;
    setSignature(value: Uint8Array | string): TransactionEvent;
    getIndex(): number;
    setIndex(value: number): TransactionEvent;
    getIsVote(): boolean;
    setIsVote(value: boolean): TransactionEvent;

    hasTransaction(): boolean;
    clearTransaction(): void;
    getTransaction(): SanitizedTransaction | undefined;
    setTransaction(value?: SanitizedTransaction): TransactionEvent;

    hasTransactionStatusMeta(): boolean;
    clearTransactionStatusMeta(): void;
    getTransactionStatusMeta(): TransactionStatusMeta | undefined;
    setTransactionStatusMeta(value?: TransactionStatusMeta): TransactionEvent;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TransactionEvent.AsObject;
    static toObject(includeInstance: boolean, msg: TransactionEvent): TransactionEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TransactionEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TransactionEvent;
    static deserializeBinaryFromReader(message: TransactionEvent, reader: jspb.BinaryReader): TransactionEvent;
}

export namespace TransactionEvent {
    export type AsObject = {
        slot: number,
        signature: Uint8Array | string,
        index: number,
        isVote: boolean,
        transaction?: SanitizedTransaction.AsObject,
        transactionStatusMeta?: TransactionStatusMeta.AsObject,
    }
}

export class InnerInstruction extends jspb.Message { 

    hasInstruction(): boolean;
    clearInstruction(): void;
    getInstruction(): CompiledInstruction | undefined;
    setInstruction(value?: CompiledInstruction): InnerInstruction;

    hasStackHeight(): boolean;
    clearStackHeight(): void;
    getStackHeight(): number | undefined;
    setStackHeight(value: number): InnerInstruction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InnerInstruction.AsObject;
    static toObject(includeInstance: boolean, msg: InnerInstruction): InnerInstruction.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InnerInstruction, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InnerInstruction;
    static deserializeBinaryFromReader(message: InnerInstruction, reader: jspb.BinaryReader): InnerInstruction;
}

export namespace InnerInstruction {
    export type AsObject = {
        instruction?: CompiledInstruction.AsObject,
        stackHeight?: number,
    }
}

export class InnerInstructions extends jspb.Message { 
    getIndex(): number;
    setIndex(value: number): InnerInstructions;
    clearInstructionsList(): void;
    getInstructionsList(): Array<InnerInstruction>;
    setInstructionsList(value: Array<InnerInstruction>): InnerInstructions;
    addInstructions(value?: InnerInstruction, index?: number): InnerInstruction;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InnerInstructions.AsObject;
    static toObject(includeInstance: boolean, msg: InnerInstructions): InnerInstructions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InnerInstructions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InnerInstructions;
    static deserializeBinaryFromReader(message: InnerInstructions, reader: jspb.BinaryReader): InnerInstructions;
}

export namespace InnerInstructions {
    export type AsObject = {
        index: number,
        instructionsList: Array<InnerInstruction.AsObject>,
    }
}

export class UiTokenAmount extends jspb.Message { 
    getUiAmount(): number;
    setUiAmount(value: number): UiTokenAmount;
    getDecimals(): number;
    setDecimals(value: number): UiTokenAmount;
    getAmount(): string;
    setAmount(value: string): UiTokenAmount;
    getUiAmountString(): string;
    setUiAmountString(value: string): UiTokenAmount;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UiTokenAmount.AsObject;
    static toObject(includeInstance: boolean, msg: UiTokenAmount): UiTokenAmount.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UiTokenAmount, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UiTokenAmount;
    static deserializeBinaryFromReader(message: UiTokenAmount, reader: jspb.BinaryReader): UiTokenAmount;
}

export namespace UiTokenAmount {
    export type AsObject = {
        uiAmount: number,
        decimals: number,
        amount: string,
        uiAmountString: string,
    }
}

export class TransactionTokenBalance extends jspb.Message { 
    getAccountIndex(): number;
    setAccountIndex(value: number): TransactionTokenBalance;
    getMint(): string;
    setMint(value: string): TransactionTokenBalance;

    hasUiTokenAmount(): boolean;
    clearUiTokenAmount(): void;
    getUiTokenAmount(): UiTokenAmount | undefined;
    setUiTokenAmount(value?: UiTokenAmount): TransactionTokenBalance;
    getOwner(): string;
    setOwner(value: string): TransactionTokenBalance;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TransactionTokenBalance.AsObject;
    static toObject(includeInstance: boolean, msg: TransactionTokenBalance): TransactionTokenBalance.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TransactionTokenBalance, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TransactionTokenBalance;
    static deserializeBinaryFromReader(message: TransactionTokenBalance, reader: jspb.BinaryReader): TransactionTokenBalance;
}

export namespace TransactionTokenBalance {
    export type AsObject = {
        accountIndex: number,
        mint: string,
        uiTokenAmount?: UiTokenAmount.AsObject,
        owner: string,
    }
}

export class Reward extends jspb.Message { 
    getPubkey(): string;
    setPubkey(value: string): Reward;
    getLamports(): number;
    setLamports(value: number): Reward;
    getPostBalance(): number;
    setPostBalance(value: number): Reward;
    getRewardType(): number;
    setRewardType(value: number): Reward;
    getCommission(): number;
    setCommission(value: number): Reward;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Reward.AsObject;
    static toObject(includeInstance: boolean, msg: Reward): Reward.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Reward, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Reward;
    static deserializeBinaryFromReader(message: Reward, reader: jspb.BinaryReader): Reward;
}

export namespace Reward {
    export type AsObject = {
        pubkey: string,
        lamports: number,
        postBalance: number,
        rewardType: number,
        commission: number,
    }
}

export class TransactionStatusMeta extends jspb.Message { 
    getIsStatusErr(): boolean;
    setIsStatusErr(value: boolean): TransactionStatusMeta;
    getFee(): number;
    setFee(value: number): TransactionStatusMeta;
    clearPreBalancesList(): void;
    getPreBalancesList(): Array<number>;
    setPreBalancesList(value: Array<number>): TransactionStatusMeta;
    addPreBalances(value: number, index?: number): number;
    clearPostBalancesList(): void;
    getPostBalancesList(): Array<number>;
    setPostBalancesList(value: Array<number>): TransactionStatusMeta;
    addPostBalances(value: number, index?: number): number;
    clearInnerInstructionsList(): void;
    getInnerInstructionsList(): Array<InnerInstructions>;
    setInnerInstructionsList(value: Array<InnerInstructions>): TransactionStatusMeta;
    addInnerInstructions(value?: InnerInstructions, index?: number): InnerInstructions;
    clearLogMessagesList(): void;
    getLogMessagesList(): Array<string>;
    setLogMessagesList(value: Array<string>): TransactionStatusMeta;
    addLogMessages(value: string, index?: number): string;
    clearPreTokenBalancesList(): void;
    getPreTokenBalancesList(): Array<TransactionTokenBalance>;
    setPreTokenBalancesList(value: Array<TransactionTokenBalance>): TransactionStatusMeta;
    addPreTokenBalances(value?: TransactionTokenBalance, index?: number): TransactionTokenBalance;
    clearPostTokenBalancesList(): void;
    getPostTokenBalancesList(): Array<TransactionTokenBalance>;
    setPostTokenBalancesList(value: Array<TransactionTokenBalance>): TransactionStatusMeta;
    addPostTokenBalances(value?: TransactionTokenBalance, index?: number): TransactionTokenBalance;
    clearRewardsList(): void;
    getRewardsList(): Array<Reward>;
    setRewardsList(value: Array<Reward>): TransactionStatusMeta;
    addRewards(value?: Reward, index?: number): Reward;
    getErrorInfo(): string;
    setErrorInfo(value: string): TransactionStatusMeta;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TransactionStatusMeta.AsObject;
    static toObject(includeInstance: boolean, msg: TransactionStatusMeta): TransactionStatusMeta.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TransactionStatusMeta, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TransactionStatusMeta;
    static deserializeBinaryFromReader(message: TransactionStatusMeta, reader: jspb.BinaryReader): TransactionStatusMeta;
}

export namespace TransactionStatusMeta {
    export type AsObject = {
        isStatusErr: boolean,
        fee: number,
        preBalancesList: Array<number>,
        postBalancesList: Array<number>,
        innerInstructionsList: Array<InnerInstructions.AsObject>,
        logMessagesList: Array<string>,
        preTokenBalancesList: Array<TransactionTokenBalance.AsObject>,
        postTokenBalancesList: Array<TransactionTokenBalance.AsObject>,
        rewardsList: Array<Reward.AsObject>,
        errorInfo: string,
    }
}

export class TransactionEventWrapper extends jspb.Message { 
    getStreamType(): StreamType;
    setStreamType(value: StreamType): TransactionEventWrapper;

    hasTransaction(): boolean;
    clearTransaction(): void;
    getTransaction(): TransactionEvent | undefined;
    setTransaction(value?: TransactionEvent): TransactionEventWrapper;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TransactionEventWrapper.AsObject;
    static toObject(includeInstance: boolean, msg: TransactionEventWrapper): TransactionEventWrapper.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TransactionEventWrapper, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TransactionEventWrapper;
    static deserializeBinaryFromReader(message: TransactionEventWrapper, reader: jspb.BinaryReader): TransactionEventWrapper;
}

export namespace TransactionEventWrapper {
    export type AsObject = {
        streamType: StreamType,
        transaction?: TransactionEvent.AsObject,
    }
}

export class MessageWrapper extends jspb.Message { 

    hasAccountUpdate(): boolean;
    clearAccountUpdate(): void;
    getAccountUpdate(): SubscribeUpdateAccountInfo | undefined;
    setAccountUpdate(value?: SubscribeUpdateAccountInfo): MessageWrapper;

    hasSlot(): boolean;
    clearSlot(): void;
    getSlot(): SlotStatusEvent | undefined;
    setSlot(value?: SlotStatusEvent): MessageWrapper;

    hasTransaction(): boolean;
    clearTransaction(): void;
    getTransaction(): TransactionEventWrapper | undefined;
    setTransaction(value?: TransactionEventWrapper): MessageWrapper;

    getEventMessageCase(): MessageWrapper.EventMessageCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageWrapper.AsObject;
    static toObject(includeInstance: boolean, msg: MessageWrapper): MessageWrapper.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageWrapper, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageWrapper;
    static deserializeBinaryFromReader(message: MessageWrapper, reader: jspb.BinaryReader): MessageWrapper;
}

export namespace MessageWrapper {
    export type AsObject = {
        accountUpdate?: SubscribeUpdateAccountInfo.AsObject,
        slot?: SlotStatusEvent.AsObject,
        transaction?: TransactionEventWrapper.AsObject,
    }

    export enum EventMessageCase {
        EVENT_MESSAGE_NOT_SET = 0,
        ACCOUNT_UPDATE = 1,
        SLOT = 2,
        TRANSACTION = 3,
    }

}

export enum StreamType {
    STREAM_TYPE_UNSPECIFIED = 0,
    STREAM_TYPE_FILTERED = 1,
    STREAM_TYPE_WALLET = 2,
    STREAM_TYPE_ACCOUNT = 3,
}
