package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	thorclient "github.com/thorlabsDev/ThorStreamer/sdks/go/client"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create client
	client, err := thorclient.NewClient(thorclient.Config{
		ServerAddr:     os.Getenv("SERVER_ADDRESS"),
		Token:          os.Getenv("AUTH_TOKEN"),
		DefaultTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Example 1: Subscribe to transactions
	go subscribeToTransactions(ctx, client)

	// Example 2: Subscribe to slot status
	go subscribeToSlots(ctx, client)

	// Example 3: Subscribe to wallet transactions
	go subscribeToWallets(ctx, client, []string{
		"YourWalletAddressHere",
	})

	// Example 4: Subscribe to account updates
	go subscribeToAccounts(ctx, client, []string{
		"AccountAddress1",
		"AccountAddress2",
	}, []string{
		"OwnerAddress1",
	})

	// Example 5: Subscribe to Thor updates
	go subscribeToThorUpdates(ctx, client)

	// Wait for shutdown
	<-ctx.Done()
	log.Println("Shutdown complete")
}

func subscribeToTransactions(ctx context.Context, client *thorclient.Client) {
	stream, err := client.SubscribeToTransactions(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to transactions: %v", err)
		return
	}

	log.Println("Subscribed to transactions")
	for {
		msg, err := stream.Recv()
		if err != nil {
			if thorclient.IsStreamDone(err) {
				log.Println("Transaction stream closed")
				return
			}
			log.Printf("Error receiving transaction: %v", err)
			return
		}

		if tx := msg.GetTransaction(); tx != nil {
			log.Printf("Received transaction: slot=%d, signature=%x",
				tx.Transaction.Slot,
				tx.Transaction.Signature[:8])
		}
	}
}

func subscribeToSlots(ctx context.Context, client *thorclient.Client) {
	stream, err := client.SubscribeToSlotStatus(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to slots: %v", err)
		return
	}

	log.Println("Subscribed to slots")
	for {
		msg, err := stream.Recv()
		if err != nil {
			if thorclient.IsStreamDone(err) {
				log.Println("Slot stream closed")
				return
			}
			log.Printf("Error receiving slot: %v", err)
			return
		}

		if slot := msg.GetSlot(); slot != nil {
			log.Printf("Received slot: slot=%d, status=%d, height=%d",
				slot.Slot, slot.Status, slot.BlockHeight)
		}
	}
}

func subscribeToWallets(ctx context.Context, client *thorclient.Client, wallets []string) {
	stream, err := client.SubscribeToWalletTransactions(ctx, wallets)
	if err != nil {
		log.Printf("Failed to subscribe to wallets: %v", err)
		return
	}

	log.Printf("Subscribed to %d wallets", len(wallets))
	for {
		msg, err := stream.Recv()
		if err != nil {
			if thorclient.IsStreamDone(err) {
				log.Println("Wallet stream closed")
				return
			}
			log.Printf("Error receiving wallet transaction: %v", err)
			return
		}

		if tx := msg.GetTransaction(); tx != nil {
			log.Printf("Received wallet transaction: slot=%d, signature=%x",
				tx.Transaction.Slot,
				tx.Transaction.Signature[:8])
		}
	}
}

func subscribeToAccounts(ctx context.Context, client *thorclient.Client, accounts, owners []string) {
	stream, err := client.SubscribeToAccountUpdates(ctx, accounts, owners)
	if err != nil {
		log.Printf("Failed to subscribe to accounts: %v", err)
		return
	}

	log.Printf("Subscribed to %d accounts and %d owners", len(accounts), len(owners))
	for {
		msg, err := stream.Recv()
		if err != nil {
			if thorclient.IsStreamDone(err) {
				log.Println("Account stream closed")
				return
			}
			log.Printf("Error receiving account update: %v", err)
			return
		}

		if acc := msg.GetAccountUpdate(); acc != nil {
			log.Printf("Received account update: pubkey=%x, lamports=%d",
				acc.Pubkey[:8], acc.Lamports)
		}
	}
}

func subscribeToThorUpdates(ctx context.Context, client *thorclient.Client) {
	stream, err := client.SubscribeToThorUpdates(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to Thor updates: %v", err)
		return
	}

	log.Println("Subscribed to Thor updates")
	for {
		msg, err := stream.Recv()
		if err != nil {
			if thorclient.IsStreamDone(err) {
				log.Println("Thor stream closed")
				return
			}
			log.Printf("Error receiving Thor update: %v", err)
			return
		}

		// Process based on message type
		switch {
		case msg.GetSlot() != nil:
			slot := msg.GetSlot()
			log.Printf("Thor slot: slot=%d, status=%d", slot.Slot, slot.Status)
		case msg.GetTransaction() != nil:
			tx := msg.GetTransaction()
			log.Printf("Thor transaction: slot=%d", tx.Transaction.Slot)
		case msg.GetAccountUpdate() != nil:
			acc := msg.GetAccountUpdate()
			log.Printf("Thor account: lamports=%d", acc.Lamports)
		}
	}
}
