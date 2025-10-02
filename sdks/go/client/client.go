package thorclient

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"time"

	pb "github.com/thorlabsDev/ThorStreamer/sdks/go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn           *grpc.ClientConn
	eventClient    pb.EventPublisherClient
	thorClient     pb.ThorStreamerClient
	token          string
	defaultTimeout time.Duration
}

type Config struct {
	ServerAddr     string
	Token          string
	DefaultTimeout time.Duration
	MaxRetries     int
}

type TransactionStream struct {
	stream pb.EventPublisher_SubscribeToTransactionsClient
	ctx    context.Context
}

type SlotStream struct {
	stream pb.EventPublisher_SubscribeToSlotStatusClient
}

type WalletStream struct {
	stream pb.EventPublisher_SubscribeToWalletTransactionsClient
}

type AccountStream struct {
	stream pb.EventPublisher_SubscribeToAccountUpdatesClient
}

type ThorStream struct {
	stream pb.ThorStreamer_StreamUpdatesClient
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.DefaultTimeout == 0 {
		cfg.DefaultTimeout = 30 * time.Second
	}

	conn, err := grpc.NewClient(
		cfg.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024),
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:           conn,
		eventClient:    pb.NewEventPublisherClient(conn),
		thorClient:     pb.NewThorStreamerClient(conn),
		token:          cfg.Token,
		defaultTimeout: cfg.DefaultTimeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// contextWithAuth adds authentication token to context
func (c *Client) contextWithAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", c.token)
}

// SubscribeToTransactions subscribes to transaction events
func (c *Client) SubscribeToTransactions(ctx context.Context) (*TransactionStream, error) {
	authCtx := c.contextWithAuth(ctx)
	stream, err := c.eventClient.SubscribeToTransactions(authCtx, &emptypb.Empty{}) // Changed here
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &TransactionStream{stream: stream, ctx: authCtx}, nil
}

// Recv receives the next transaction message
func (ts *TransactionStream) Recv() (*pb.MessageWrapper, error) {
	resp, err := ts.stream.Recv()
	if err != nil {
		return nil, err
	}

	var wrapper pb.MessageWrapper
	if err := proto.Unmarshal(resp.Data, &wrapper); err != nil { // Changed here
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &wrapper, nil
}

// SubscribeToSlotStatus subscribes to slot status events
func (c *Client) SubscribeToSlotStatus(ctx context.Context) (*SlotStream, error) {
	authCtx := c.contextWithAuth(ctx)
	stream, err := c.eventClient.SubscribeToSlotStatus(authCtx, &emptypb.Empty{}) // Changed here
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &SlotStream{stream: stream}, nil
}

// Recv receives the next slot status message
func (ss *SlotStream) Recv() (*pb.MessageWrapper, error) {
	resp, err := ss.stream.Recv()
	if err != nil {
		return nil, err
	}

	var wrapper pb.MessageWrapper
	if err := proto.Unmarshal(resp.Data, &wrapper); err != nil { // Changed here
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &wrapper, nil
}

// WalletStream represents a wallet transaction subscription

// SubscribeToWalletTransactions subscribes to wallet transaction events
func (c *Client) SubscribeToWalletTransactions(ctx context.Context, wallets []string) (*WalletStream, error) {
	authCtx := c.contextWithAuth(ctx)
	req := &pb.SubscribeWalletRequest{WalletAddress: wallets}
	stream, err := c.eventClient.SubscribeToWalletTransactions(authCtx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &WalletStream{stream: stream}, nil
}

// Recv receives the next wallet transaction message
func (ws *WalletStream) Recv() (*pb.MessageWrapper, error) {
	resp, err := ws.stream.Recv()
	if err != nil {
		return nil, err
	}

	var wrapper pb.MessageWrapper
	if err := proto.Unmarshal(resp.Data, &wrapper); err != nil { // Changed here
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &wrapper, nil
}

// SubscribeToAccountUpdates subscribes to account update events
func (c *Client) SubscribeToAccountUpdates(ctx context.Context, accounts, owners []string) (*AccountStream, error) {
	authCtx := c.contextWithAuth(ctx)
	req := &pb.SubscribeAccountsRequest{
		AccountAddress: accounts,
		OwnerAddress:   owners,
	}
	stream, err := c.eventClient.SubscribeToAccountUpdates(authCtx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &AccountStream{stream: stream}, nil
}

// Recv receives the next account update message
func (as *AccountStream) Recv() (*pb.MessageWrapper, error) {
	resp, err := as.stream.Recv()
	if err != nil {
		return nil, err
	}

	var wrapper pb.MessageWrapper
	if err := proto.Unmarshal(resp.Data, &wrapper); err != nil { // Changed here
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &wrapper, nil
}

// SubscribeToThorUpdates subscribes to Thor update events
func (c *Client) SubscribeToThorUpdates(ctx context.Context) (*ThorStream, error) {
	stream, err := c.thorClient.StreamUpdates(ctx, &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &ThorStream{stream: stream}, nil
}

// Recv receives the next Thor update message
func (ts *ThorStream) Recv() (*pb.MessageWrapper, error) {
	return ts.stream.Recv()
}

// Helper function to check if stream is done
func IsStreamDone(err error) bool {
	return err == io.EOF || err == context.Canceled
}
