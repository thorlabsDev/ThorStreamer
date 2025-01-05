package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example/types"
)

// ConnectWithRetry establishes a connection to the gRPC server with retry logic
func ConnectWithRetry(ctx context.Context, config *types.Config) (*grpc.ClientConn, error) {
	var retryDelay = time.Second

	for i := 0; i < config.MaxRetries; i++ {
		conn, err := grpc.DialContext(ctx, config.ServerAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err == nil {
			return conn, nil
		}

		log.Printf("Connection attempt %d/%d failed: %v", i+1, config.MaxRetries, err)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelay):
			retryDelay *= 2 // Exponential backoff
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts", config.MaxRetries)
}

// HandleSubscription manages a subscription and its error handling
func HandleSubscription(ctx context.Context, subFunc func() error) {
	subscriptionCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handle subscription errors
	errCh := make(chan error, 1)
	go func() {
		errCh <- subFunc()
	}()

	// Wait for either context cancellation or subscription error
	select {
	case <-subscriptionCtx.Done():
		return
	case err := <-errCh:
		if err != nil && !IsContextCanceled(err) {
			log.Printf("Subscription error: %v", err)
		}
	}
}

// IsContextCanceled checks if an error is due to context cancellation
func IsContextCanceled(err error) bool {
	if err == context.Canceled {
		return true
	}
	if status, ok := status.FromError(err); ok {
		return status.Code() == codes.Canceled
	}
	return false
}
