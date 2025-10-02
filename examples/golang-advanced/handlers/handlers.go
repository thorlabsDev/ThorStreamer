package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

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

		// Check if this is a transaction event
		switch msgWrapper.EventMessage.(type) {
		case *pb.MessageWrapper_Transaction:
			txWrapper := msgWrapper.GetTransaction()
			if txWrapper == nil || txWrapper.Transaction == nil {
				continue
			}

			tx := txWrapper.Transaction

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
		default:
			// Not a transaction event, skip
			continue
		}
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

		// Check if this is a transaction event
		switch msgWrapper.EventMessage.(type) {
		case *pb.MessageWrapper_Transaction:
			txWrapper := msgWrapper.GetTransaction()
			if txWrapper == nil || txWrapper.Transaction == nil {
				continue
			}

			if txWrapper.StreamType != pb.StreamType_STREAM_TYPE_WALLET {
				continue
			}

			printer.PrintDetailedTransaction(txWrapper.Transaction)
		default:
			continue
		}
	}
}

// SubscribeToAccountUpdates subscribes to account updates
func SubscribeToAccountUpdates(ctx context.Context, client pb.EventPublisherClient) error {
	// Get account and owner addresses from user
	fmt.Println("\nEnter account addresses to monitor (comma-separated, or press Enter to skip):")
	accountInput := utils.Prompt("")
	var accountAddresses []string
	if accountInput != "" {
		for _, addr := range strings.Split(accountInput, ",") {
			accountAddresses = append(accountAddresses, strings.TrimSpace(addr))
		}
	}

	fmt.Println("Enter owner addresses to filter by (comma-separated, or press Enter to skip):")
	ownerInput := utils.Prompt("")
	var ownerAddresses []string
	if ownerInput != "" {
		for _, addr := range strings.Split(ownerInput, ",") {
			ownerAddresses = append(ownerAddresses, strings.TrimSpace(addr))
		}
	}

	if len(accountAddresses) == 0 && len(ownerAddresses) == 0 {
		return fmt.Errorf("at least one account or owner address is required")
	}

	stream, err := client.SubscribeToAccountUpdates(ctx, &pb.SubscribeAccountsRequest{
		AccountAddress: accountAddresses,
		OwnerAddress:   ownerAddresses,
	})
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

		// Check if this is an account update event
		switch msgWrapper.EventMessage.(type) {
		case *pb.MessageWrapper_AccountUpdate:
			accountUpdate := msgWrapper.GetAccountUpdate()
			if accountUpdate == nil {
				continue
			}
			printer.PrintAccountUpdate(accountUpdate)
		default:
			continue
		}
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

		// Check if this is a slot event
		switch msgWrapper.EventMessage.(type) {
		case *pb.MessageWrapper_Slot:
			slotEvent := msgWrapper.GetSlot()
			if slotEvent == nil {
				continue
			}
			printer.PrintSlotStatus(slotEvent)
		default:
			continue
		}
	}
}
