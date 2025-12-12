package main

import (
	"context"
	pb "example/proto"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/metadata"

	"example/client"
	"example/config"
	"example/handlers"
	"example/utils"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nReceived shutdown signal. Cleaning up...")
		cancel()
	}()

	// Add authentication token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", cfg.AuthToken)

	// Connect to the server with retry logic
	conn, err := client.ConnectWithRetry(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	eventClient := pb.NewEventPublisherClient(conn)

	// Main menu loop
	for {
		fmt.Println("\nSelect subscription type:")
		fmt.Println("1. Program-filtered Transactions")
		fmt.Println("2. Account Updates")
		fmt.Println("3. Slot Status")
		fmt.Println("4. Exit")

		select {
		case <-ctx.Done():
			return
		default:
			choice := utils.Prompt("\nEnter your choice (1-4): ")

			switch choice {
			case "1":
				client.HandleSubscription(ctx, func() error {
					return handlers.SubscribeToFilteredTransactions(ctx, eventClient, cfg)
				})
			case "2":
				client.HandleSubscription(ctx, func() error {
					return handlers.SubscribeToAccountUpdates(ctx, eventClient)
				})
			case "3":
				client.HandleSubscription(ctx, func() error {
					return handlers.SubscribeToSlotStatus(ctx, eventClient)
				})
			case "4":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
	}
}
