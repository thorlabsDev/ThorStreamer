package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "thor_grpc/proto"

	"github.com/mr-tron/base58"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Config represents the client configuration
type Config struct {
	ServerAddress string `json:"server_address"`
	AuthToken     string `json:"auth_token"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

func debugTransaction(tx *pb.TransactionEvent) {
	fmt.Println("\n🔍 Transaction Debug Information:")
	fmt.Printf("├─ Signature: %s\n", base58.Encode(tx.Signature))
	fmt.Printf("├─ Slot: %d\n", tx.Slot)

	if tx.Transaction == nil {
		fmt.Println("├─ ⚠️  Transaction is nil!")
		return
	}

	if tx.Transaction.Message == nil {
		fmt.Println("├─ ⚠️  Message is nil!")
		return
	}

	fmt.Printf("├─ Version: %d (%s)\n",
		tx.Transaction.Message.Version,
		getVersionString(tx.Transaction.Message.Version))

	debugHeader(tx.Transaction.Message)
	debugAccountKeys(tx.Transaction.Message)
	debugBlockhash(tx.Transaction.Message)
	debugInstructions(tx.Transaction.Message)
}

func getVersionString(version uint32) string {
	switch version {
	case 0:
		return "Legacy"
	case 1:
		return "V0"
	default:
		return fmt.Sprintf("Unknown(%d)", version)
	}
}

func debugHeader(msg *pb.Message) {
	fmt.Println("├─ Header:")
	if msg.Header == nil {
		fmt.Println("│  └─ ⚠️  Header is nil!")
		return
	}
	fmt.Printf("│  ├─ NumRequiredSignatures: %d\n", msg.Header.NumRequiredSignatures)
	fmt.Printf("│  ├─ NumReadonlySignedAccounts: %d\n", msg.Header.NumReadonlySignedAccounts)
	fmt.Printf("│  └─ NumReadonlyUnsignedAccounts: %d\n", msg.Header.NumReadonlyUnsignedAccounts)
}

func debugAccountKeys(msg *pb.Message) {
	fmt.Printf("├─ Account Keys (%d):\n", len(msg.AccountKeys))
	if len(msg.AccountKeys) == 0 {
		fmt.Println("│  └─ ⚠️  No account keys!")
		return
	}

	for i, key := range msg.AccountKeys {
		if i < 5 {
			fmt.Printf("│  ├─ [%d]: %s\n", i, base58.Encode(key))
		}
	}
	if len(msg.AccountKeys) > 5 {
		fmt.Printf("│  └─ ... and %d more keys\n", len(msg.AccountKeys)-5)
	}
}

func debugBlockhash(msg *pb.Message) {
	fmt.Println("├─ Recent Blockhash:")
	if len(msg.RecentBlockHash) == 0 {
		fmt.Println("│  └─ ⚠️  Blockhash is empty!")
		return
	}
	fmt.Printf("│  └─ %s\n", base58.Encode(msg.RecentBlockHash))
}

func debugInstructions(msg *pb.Message) {
	fmt.Printf("├─ Instructions (%d):\n", len(msg.Instructions))
	if len(msg.Instructions) == 0 {
		fmt.Println("│  └─ ⚠️  No instructions!")
		return
	}

	for i, ix := range msg.Instructions {
		if i < 3 {
			fmt.Printf("│  ├─ Instruction %d:\n", i)
			fmt.Printf("│  │  ├─ Program ID Index: %d\n", ix.ProgramIdIndex)
			fmt.Printf("│  │  ├─ Account Indexes: %v\n", ix.Accounts)
			fmt.Printf("│  │  └─ Data Length: %d bytes\n", len(ix.Data))
		}
	}
	if len(msg.Instructions) > 3 {
		fmt.Printf("│  └─ ... and %d more instructions\n", len(msg.Instructions)-3)
	}
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", config.AuthToken,
	)

	conn, err := grpc.Dial(config.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewEventPublisherClient(conn)

	fmt.Printf("🔍 Starting Transaction Debugger on %s\n", config.ServerAddress)
	fmt.Println("--------------------------------")

	stream, err := client.SubscribeToTransactions(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving transaction: %v", err)
			break
		}

		// Unmarshal into MessageWrapper first
		var msgWrapper pb.MessageWrapper
		if err := proto.Unmarshal(resp.Data, &msgWrapper); err != nil {
			log.Printf("Failed to unmarshal MessageWrapper: %v", err)
			continue
		}

		eventMessage := msgWrapper.GetEventMessage()
		if eventMessage == nil {
			log.Println("No event_message found in MessageWrapper")
			continue
		}

		// event_message is a oneof: check if it's a transaction
		txEvent, ok := eventMessage.(*pb.MessageWrapper_Transaction)
		if !ok {
			log.Println("Received a message that is not a transaction event")
			continue
		}

		if txEvent.Transaction == nil {
			log.Println("No transaction found in MessageWrapper")
			continue
		}

		debugTransaction(txEvent.Transaction)
		fmt.Println("└─ End Debug Info\n")
	}
}
