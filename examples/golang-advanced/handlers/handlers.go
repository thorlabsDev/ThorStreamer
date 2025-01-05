package handlers

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"example/filter"
	"example/logger"
	"example/printer"
	pb "example/proto"
	"example/types"
	"example/utils"
)

// SubscribeToFilteredTransactions subscribes to filtered transactions
func SubscribeToFilteredTransactions(ctx context.Context, client pb.EventPublisherClient, config *types.Config) error {
	txFilter := filter.NewFilter(config)
	sigLogger, err := logger.NewSignatureLogger(config.SignatureLogFile)
	if err != nil {
		return fmt.Errorf("failed to create signature logger: %v", err)
	}
	defer sigLogger.Close()

	stream, err := client.SubscribeToTransactions(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to transactions: %v", err)
	}

	fmt.Printf("\n📡 Monitoring transactions for programs: %v\n", config.ProgramFilters)
	fmt.Printf("📝 Logging signatures to: %s\n", config.SignatureLogFile)
	fmt.Println("-------------------------------------------")

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		var msgWrapper pb.MessageWrapper
		if err := proto.Unmarshal(resp.Data, &msgWrapper); err != nil {
			log.Printf("Failed to unmarshal MessageWrapper: %v", err)
			continue
		}

		txWrapper, ok := msgWrapper.GetEventMessage().(*pb.MessageWrapper_Transaction)
		if !ok || txWrapper.Transaction == nil || txWrapper.Transaction.Transaction == nil {
			continue
		}

		tx := txWrapper.Transaction.Transaction

		if !txFilter.WantsTransaction(tx) {
			continue
		}

		// Extract program ID for logging
		programIDs := filter.ExtractProgramIDs(tx)
		primaryProgramID := "unknown"
		if len(programIDs) > 0 {
			primaryProgramID = programIDs[0]
		}

		// Log signature
		if err := sigLogger.LogSignature(
			tx.Signature,
			tx.Slot,
			primaryProgramID,
			tx.TransactionStatusMeta != nil && !tx.TransactionStatusMeta.IsStatusErr,
		); err != nil {
			log.Printf("Failed to log signature: %v", err)
		}

		printer.PrintTransaction(tx)
	}
}

// SubscribeToWalletTransactions subscribes to wallet-specific transactions
func SubscribeToWalletTransactions(ctx context.Context, client pb.EventPublisherClient) error {
	wallets, err := utils.GetUserWallets()
	if err != nil {
		return err
	}

	fmt.Printf("\n📡 Monitoring transactions for %d wallets:\n", len(wallets))
	for i, wallet := range wallets {
		fmt.Printf("%d. %s\n", i+1, wallet)
	}
	fmt.Println("----------------------------------")

	stream, err := client.SubscribeToWalletTransactions(ctx, &pb.SubscribeWalletRequest{
		WalletAddress: wallets,
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to wallet transactions: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		var msgWrapper pb.MessageWrapper
		if err := proto.Unmarshal(resp.Data, &msgWrapper); err != nil {
			log.Printf("Failed to unmarshal MessageWrapper: %v", err)
			continue
		}

		txWrapper, ok := msgWrapper.GetEventMessage().(*pb.MessageWrapper_Transaction)
		if !ok || txWrapper.Transaction == nil || txWrapper.Transaction.Transaction == nil {
			continue
		}

		if txWrapper.Transaction.StreamType != pb.StreamType_STREAM_TYPE_WALLET {
			continue
		}

		printer.PrintDetailedTransaction(txWrapper.Transaction.Transaction)
	}
}

// SubscribeToAccountUpdates subscribes to account updates
func SubscribeToAccountUpdates(ctx context.Context, client pb.EventPublisherClient) error {
	stream, err := client.SubscribeToAccountUpdates(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to account updates: %v", err)
	}

	fmt.Println("\n📡 Monitoring account updates...")
	fmt.Println("-------------------------------")

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		var msgWrapper pb.MessageWrapper
		if err := proto.Unmarshal(resp.Data, &msgWrapper); err != nil {
			log.Printf("Failed to unmarshal MessageWrapper: %v", err)
			continue
		}

		accountEvent, ok := msgWrapper.GetEventMessage().(*pb.MessageWrapper_Account)
		if !ok || accountEvent.Account == nil {
			continue
		}

		printer.PrintAccountUpdate(accountEvent.Account)
	}
}

// SubscribeToSlotStatus subscribes to slot status updates
func SubscribeToSlotStatus(ctx context.Context, client pb.EventPublisherClient) error {
	stream, err := client.SubscribeToSlotStatus(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to slot status: %v", err)
	}

	fmt.Println("\n📡 Monitoring slot status...")
	fmt.Println("---------------------------")

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		var msgWrapper pb.MessageWrapper
		if err := proto.Unmarshal(resp.Data, &msgWrapper); err != nil {
			log.Printf("Failed to unmarshal MessageWrapper: %v", err)
			continue
		}

		slotEvent, ok := msgWrapper.GetEventMessage().(*pb.MessageWrapper_Slot)
		if !ok || slotEvent.Slot == nil {
			continue
		}

		printer.PrintSlotStatus(slotEvent.Slot)
	}
}
