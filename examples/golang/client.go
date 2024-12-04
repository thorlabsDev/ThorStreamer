package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mr-tron/base58"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "thor_grpc/proto"
)

// Config represents the client configuration
type Config struct {
	ServerAddress     string   `json:"server_address"`
	AuthToken         string   `json:"auth_token"`
	ProgramFilters    []string `json:"program_filters"`
	LogDirectory      string   `json:"log_directory"`
	IncludeVote       bool     `json:"include_vote_transactions"`
	IncludeFailed     bool     `json:"include_failed_transactions"`
	MaxRetries        int      `json:"max_retries"`
	SignatureLogFile  string   `json:"signature_log_file"`
	ChannelBufferSize int      `json:"channel_buffer_size"`
}

// Filter handles program filtering logic
type Filter struct {
	programFilters map[string]bool
	includeVote    bool
	includeFailed  bool
	mu             sync.RWMutex
}

// SignatureLogger handles signature logging
type SignatureLogger struct {
	file *os.File
	mu   sync.Mutex
}

func NewSignatureLogger(filename string) (*SignatureLogger, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open signature log file: %v", err)
	}

	return &SignatureLogger{
		file: file,
	}, nil
}

func (sl *SignatureLogger) LogSignature(signature []byte, slot uint64, programID string, success bool) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	entry := struct {
		Timestamp string `json:"timestamp"`
		Signature string `json:"signature"`
		Slot      uint64 `json:"slot"`
		ProgramID string `json:"program_id"`
		Success   bool   `json:"success"`
	}{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Signature: base58.Encode(signature),
		Slot:      slot,
		ProgramID: programID,
		Success:   success,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal signature entry: %v", err)
	}

	if _, err := sl.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write signature: %v", err)
	}

	return nil
}

func (sl *SignatureLogger) Close() error {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.file.Close()
}

func NewFilter(config *Config) *Filter {
	programFilters := make(map[string]bool)
	for _, program := range config.ProgramFilters {
		programFilters[program] = true
	}

	return &Filter{
		programFilters: programFilters,
		includeVote:    config.IncludeVote,
		includeFailed:  config.IncludeFailed,
	}
}

func (f *Filter) WantsTransaction(tx *pb.TransactionEvent) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check vote transaction filter
	if tx.IsVote && !f.includeVote {
		return false
	}

	// Check failed transaction filter
	if tx.TransactionStatusMeta != nil && tx.TransactionStatusMeta.IsStatusErr && !f.includeFailed {
		return false
	}

	// If no program filters are set, accept all transactions
	if len(f.programFilters) == 0 {
		return true
	}

	// Extract program IDs from the transaction
	programIDs := extractProgramIDs(tx)

	// Check if any program in the transaction matches our filters
	for _, programID := range programIDs {
		if f.programFilters[programID] {
			return true
		}
	}

	return false
}

func extractProgramIDs(tx *pb.TransactionEvent) []string {
	programIDs := make(map[string]bool)

	if tx.Transaction != nil && tx.Transaction.Message != nil {
		msg := tx.Transaction.Message
		// Get program IDs from account keys
		for _, key := range msg.AccountKeys {
			programIDs[base58.Encode(key)] = true
		}

		// Add program IDs from instructions
		for _, ix := range msg.Instructions {
			if int(ix.ProgramIdIndex) < len(msg.AccountKeys) {
				programID := base58.Encode(msg.AccountKeys[ix.ProgramIdIndex])
				programIDs[programID] = true
			}
		}

		// Add program IDs from loaded addresses if it's v0 transaction
		if msg.Version == 1 && msg.LoadedAddresses != nil { // v0 transaction
			for _, key := range msg.LoadedAddresses.Writable {
				programIDs[base58.Encode(key)] = true
			}
			for _, key := range msg.LoadedAddresses.Readonly {
				programIDs[base58.Encode(key)] = true
			}
		}
	}

	// Extract program IDs from logs
	if tx.TransactionStatusMeta != nil {
		for _, log_ := range tx.TransactionStatusMeta.LogMessages {
			if strings.Contains(log_, "Program ") && strings.Contains(log_, "invoke [") {
				parts := strings.Split(log_, "Program ")
				if len(parts) > 1 {
					programID := strings.Split(parts[1], " ")[0]
					programIDs[programID] = true
				}
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(programIDs))
	for programID := range programIDs {
		result = append(result, programID)
	}
	return result
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

	// Set defaults if not specified
	if config.MaxRetries == 0 {
		config.MaxRetries = 5
	}
	if config.ChannelBufferSize == 0 {
		config.ChannelBufferSize = 100
	}
	if config.SignatureLogFile == "" {
		config.SignatureLogFile = "signatures.log"
	}

	return &config, nil
}

func subscribeToFilteredTransactions(ctx context.Context, client pb.EventPublisherClient, config *Config) error {
	filter := NewFilter(config)
	logger, err := NewSignatureLogger(config.SignatureLogFile)
	if err != nil {
		return fmt.Errorf("failed to create signature logger: %v", err)
	}
	defer logger.Close()

	stream, err := client.SubscribeToTransactions(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to transactions: %v", err)
	}

	fmt.Printf("\n📡 Monitoring transactions for programs: %v\n", config.ProgramFilters)
	fmt.Printf("📝 Logging signatures to: %s\n", config.SignatureLogFile)
	fmt.Println("-------------------------------------------")

	for {
		tx, err := stream.Recv()
		if err != nil {
			return err
		}

		// Apply filter
		if !filter.WantsTransaction(tx) {
			continue
		}

		// Extract program ID for logging
		programIDs := extractProgramIDs(tx)
		primaryProgramID := "unknown"
		if len(programIDs) > 0 {
			primaryProgramID = programIDs[0]
		}

		// Log signature
		if err := logger.LogSignature(
			tx.Signature,
			tx.Slot,
			primaryProgramID,
			tx.TransactionStatusMeta != nil && !tx.TransactionStatusMeta.IsStatusErr,
		); err != nil {
			log.Printf("Failed to log signature: %v", err)
		}

		// Print transaction details
		printTransaction(tx)
	}
}

func main() {
	// Load config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", config.AuthToken)

	// Connect to the server with retry logic using config
	conn, err := connectWithRetry(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewEventPublisherClient(conn)

	for {
		fmt.Println("\nSelect subscription type:")
		fmt.Println("1. Program-filtered Transactions")
		fmt.Println("2. Account Updates")
		fmt.Println("3. Slot Status")
		fmt.Println("4. Wallet Transactions")
		fmt.Println("5. Exit")

		select {
		case <-ctx.Done():
			return
		default:
			choice := prompt("\nEnter your choice (1-5): ")

			switch choice {
			case "1":
				handleSubscription(ctx, func() error {
					return subscribeToFilteredTransactions(ctx, client, config)
				})
			case "2":
				handleSubscription(ctx, func() error {
					return subscribeToAccountUpdates(ctx, client)
				})
			case "3":
				handleSubscription(ctx, func() error {
					return subscribeToSlotStatus(ctx, client)
				})
			case "4":
				handleSubscription(ctx, func() error {
					return subscribeToWalletTransactions(ctx, client)
				})
			case "5":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
	}
}

const (
	lamportsPerSol = 1000000000 // SOL to lamports conversion rate
)

// Update the connectWithRetry function to accept the config parameter
func connectWithRetry(ctx context.Context, config *Config) (*grpc.ClientConn, error) {
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

func handleSubscription(ctx context.Context, subFunc func() error) {
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
		if err != nil && !isContextCanceled(err) {
			log.Printf("Subscription error: %v", err)
		}
	}
}

func subscribeToWalletTransactions(ctx context.Context, client pb.EventPublisherClient) error {
	// Get wallet addresses from user input
	var wallets []string
	fmt.Println("Enter up to 10 wallet addresses (one per line).")
	fmt.Println("Press Enter twice when done:")

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 10; i++ {
		fmt.Printf("Wallet %d (or blank to finish): ", i+1)
		wallet, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}

		wallet = strings.TrimSpace(wallet)
		if wallet == "" {
			break
		}
		wallets = append(wallets, wallet)
	}

	if len(wallets) == 0 {
		return fmt.Errorf("no wallet addresses provided")
	}

	fmt.Printf("\n📡 Monitoring transactions for %d wallets:\n", len(wallets))
	for i, wallet := range wallets {
		fmt.Printf("%d. %s\n", i+1, wallet)
	}
	fmt.Println("----------------------------------")

	// Subscribe to the wallets
	stream, err := client.SubscribeToWalletTransactions(ctx, &pb.SubscribeWalletRequest{
		WalletAddress: wallets,
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to wallet transactions: %v", err)
	}

	// Receive and print transactions
	for {
		tx, err := stream.Recv()
		if err != nil {
			return err
		}

		printDetailedTransaction(tx)
	}
}

func subscribeToAccountUpdates(ctx context.Context, client pb.EventPublisherClient) error {
	stream, err := client.SubscribeToAccountUpdates(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to account updates: %v", err)
	}

	fmt.Println("\n📡 Monitoring account updates...")
	fmt.Println("-------------------------------")

	for {
		account, err := stream.Recv()
		if err != nil {
			return err
		}

		printAccountUpdate(account)
	}
}

func subscribeToSlotStatus(ctx context.Context, client pb.EventPublisherClient) error {
	stream, err := client.SubscribeToSlotStatus(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to slot status: %v", err)
	}

	fmt.Println("\n📡 Monitoring slot status...")
	fmt.Println("---------------------------")

	for {
		slot, err := stream.Recv()
		if err != nil {
			return err
		}

		printSlotStatus(slot)
	}
}

// Utility functions
func prompt(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func isContextCanceled(err error) bool {
	if err == context.Canceled {
		return true
	}
	if status, ok := status.FromError(err); ok {
		return status.Code() == codes.Canceled
	}
	return false
}

func lamportsToSol(lamports uint64) string {
	return fmt.Sprintf("%.9f", float64(lamports)/float64(lamportsPerSol))
}

// Print functions for different message types
func printTransaction(tx *pb.TransactionEvent) {
	fmt.Printf("\n📜 Transaction: %s\n", base58.Encode(tx.Signature))
	fmt.Printf("├─ Slot: %d\n", tx.Slot)
	if tx.Transaction != nil {
		fmt.Printf("├─ Version: %s\n",
			getTransactionVersionString(tx.Transaction.Message.Version))
	}
	fmt.Printf("├─ Success: %v\n", !tx.TransactionStatusMeta.IsStatusErr)

	if len(tx.TransactionStatusMeta.LogMessages) > 0 {
		fmt.Println("├─ Log Messages:")
		for _, msg := range tx.TransactionStatusMeta.LogMessages {
			fmt.Printf("│  └─ %s\n", msg)
		}
	}
	fmt.Println("└─ End Transaction\n")
}

func printAccountUpdate(account *pb.UpdateAccountEvent) {
	fmt.Println("\n💳 Account Update:")
	fmt.Printf("├─ Address: %s\n", base58.Encode(account.Pubkey))
	fmt.Printf("├─ Owner: %s\n", base58.Encode(account.Owner))
	fmt.Printf("├─ Balance: %s SOL\n", lamportsToSol(account.Lamports))
	fmt.Printf("├─ Executable: %v\n", account.Executable)
	fmt.Printf("└─ Slot: %d\n\n", account.Slot)
}

func printSlotStatus(slot *pb.SlotStatusEvent) {
	fmt.Println("\n🔢 Slot Status:")
	fmt.Printf("├─ Slot: %d\n", slot.Slot)
	fmt.Printf("├─ Parent: %d\n", slot.Parent)
	fmt.Printf("└─ Status: %s\n\n", pb.SlotStatus_name[int32(slot.Status)])
}

func printDetailedTransaction(tx *pb.TransactionEvent) {
	fmt.Println("\n👛 Wallet Transaction Details:")
	fmt.Printf("├─ Signature: %s\n", base58.Encode(tx.Signature))
	fmt.Printf("├─ Slot: %d\n", tx.Slot)
	fmt.Printf("├─ Is Vote Transaction: %v\n", tx.IsVote)

	if tx.Transaction != nil {
		fmt.Printf("├─ Transaction Version: %s\n", // Added version printing
			getTransactionVersionString(tx.Transaction.Message.Version))
		printTransactionDetails(tx.Transaction)
	}

	if tx.TransactionStatusMeta != nil {
		printTransactionStatusMeta(tx.TransactionStatusMeta)
	}

	fmt.Println("└─ End Transaction\n")
}

// Helper function to convert version number to string
func getTransactionVersionString(version uint32) string {
	switch version {
	case 0:
		return "Legacy"
	case 1:
		return "V0"
	default:
		return fmt.Sprintf("Unknown Version(%d)", version)
	}
}

func printTransactionDetails(tx *pb.SanitizedTransaction) {
	fmt.Println("├─ Transaction Details:")
	fmt.Printf("│  ├─ Message Hash: %s\n", base58.Encode(tx.MessageHash))
	fmt.Printf("│  └─ Is Simple Vote: %v\n", tx.IsSimpleVoteTransaction)

	if tx.Message != nil {
		printMessage(tx.Message)
	}
}

func printMessage(msg *pb.Message) {
	if msg == nil {
		return
	}

	versionStr := "Legacy"
	if msg.Version == 1 {
		versionStr = "V0"
	}

	fmt.Printf("├─ %s Message Details:\n", versionStr)
	fmt.Printf("│  ├─ Recent Blockhash: %s\n", base58.Encode(msg.RecentBlockHash))
	printMessageHeader(msg.Header)
	printAccountKeys("│  ├─", msg.AccountKeys)

	if msg.Version == 1 { // V0 transaction specific fields
		printLoadedAddresses("│  ├─", msg.LoadedAddresses)
		printAddressTableLookups("│  ├─", msg.AddressTableLookups)
	}

	printInstructions("│  └─", msg.Instructions)
}

func printMessageHeader(header *pb.MessageHeader) {
	if header == nil {
		return
	}

	fmt.Println("│  ├─ Header Info:")
	fmt.Printf("│  │  ├─ Required Signatures: %d\n", header.NumRequiredSignatures)
	fmt.Printf("│  │  ├─ Readonly Signed: %d\n", header.NumReadonlySignedAccounts)
	fmt.Printf("│  │  └─ Readonly Unsigned: %d\n", header.NumReadonlyUnsignedAccounts)
}

func printAccountKeys(prefix string, keys [][]byte) {
	fmt.Printf("%s Account Keys: %d\n", prefix, len(keys))
	for i, key := range keys {
		fmt.Printf("│  │  ├─ [%d] %s\n", i, base58.Encode(key))
	}
}

func printLoadedAddresses(prefix string, addresses *pb.LoadedAddresses) {
	if addresses == nil {
		return
	}

	if len(addresses.Writable) > 0 {
		fmt.Printf("%s Writable Addresses:\n", prefix)
		for i, addr := range addresses.Writable {
			fmt.Printf("│  │  ├─ [%d] %s\n", i, base58.Encode(addr))
		}
	}

	if len(addresses.Readonly) > 0 {
		fmt.Printf("%s Readonly Addresses:\n", prefix)
		for i, addr := range addresses.Readonly {
			fmt.Printf("│  │  ├─ [%d] %s\n", i, base58.Encode(addr))
		}
	}
}

func printAddressTableLookups(prefix string, lookups []*pb.MessageAddressTableLookup) {
	if len(lookups) == 0 {
		return
	}

	fmt.Printf("%s Address Table Lookups:\n", prefix)
	for i, lookup := range lookups {
		fmt.Printf("│  │  ├─ Lookup %d:\n", i)
		fmt.Printf("│  │  │  ├─ Account Key: %s\n", base58.Encode(lookup.AccountKey))
		fmt.Printf("│  │  │  ├─ Writable Indexes: %s\n",
			base58.Encode(lookup.WritableIndexes))
		fmt.Printf("│  │  │  └─ Readonly Indexes: %s\n",
			base58.Encode(lookup.ReadonlyIndexes))
	}
}

func printInstructions(prefix string, instructions []*pb.CompiledInstruction) {
	fmt.Printf("%s Instructions: %d\n", prefix, len(instructions))
	for i, ix := range instructions {
		fmt.Printf("│     ├─ Instruction %d:\n", i)
		fmt.Printf("│     │  ├─ Program ID Index: %d\n", ix.ProgramIdIndex)
		fmt.Printf("│     │  ├─ Account Indexes: %v\n", ix.Accounts)
		fmt.Printf("│     │  └─ Data: %s\n", base58.Encode(ix.Data))
	}
}

func printTransactionStatusMeta(meta *pb.TransactionStatusMeta) {
	fmt.Println("├─ Status Metadata:")
	fmt.Printf("│  ├─ Status: %s\n", formatStatus(meta.IsStatusErr, meta.ErrorInfo))
	fmt.Printf("│  ├─ Fee: %s SOL\n", lamportsToSol(meta.Fee))

	printBalanceChanges(meta)
	printTokenBalances(meta)
	printInnerInstructions(meta.InnerInstructions)
	printLogMessages(meta.LogMessages)
	printRewards(meta.Rewards)
}

func formatStatus(isErr bool, errorInfo string) string {
	if isErr {
		if errorInfo != "" {
			return fmt.Sprintf("Failed (%s)", errorInfo)
		}
		return "Failed"
	}
	return "Success"
}

func printBalanceChanges(meta *pb.TransactionStatusMeta) {
	if len(meta.PreBalances) == 0 && len(meta.PostBalances) == 0 {
		return
	}

	fmt.Println("│  ├─ Balance Changes:")
	for i := 0; i < len(meta.PreBalances) && i < len(meta.PostBalances); i++ {
		change := int64(meta.PostBalances[i]) - int64(meta.PreBalances[i])
		if change != 0 {
			fmt.Printf("│  │  ├─ Account %d: %s SOL → %s SOL (Δ %s SOL)\n",
				i,
				lamportsToSol(meta.PreBalances[i]),
				lamportsToSol(meta.PostBalances[i]),
				lamportsToSol(uint64(abs(change))))
		}
	}
}

func printTokenBalances(meta *pb.TransactionStatusMeta) {
	if len(meta.PreTokenBalances) == 0 && len(meta.PostTokenBalances) == 0 {
		return
	}

	fmt.Println("│  ├─ Token Balance Changes:")

	// Print pre-token balances
	for _, balance := range meta.PreTokenBalances {
		printTokenBalance("Pre", balance)
	}

	// Print post-token balances
	for _, balance := range meta.PostTokenBalances {
		printTokenBalance("Post", balance)
	}
}

func printTokenBalance(prefix string, balance *pb.TransactionTokenBalance) {
	if balance.UiTokenAmount == nil { // Changed from UiTokenAccount
		return
	}

	fmt.Printf("│  │  ├─ %s-Token Balance:\n", prefix)
	fmt.Printf("│  │  │  ├─ Account Index: %d\n", balance.AccountIndex)
	fmt.Printf("│  │  │  ├─ Mint: %s\n", balance.Mint)
	fmt.Printf("│  │  │  ├─ Owner: %s\n", balance.Owner)
	fmt.Printf("│  │  │  ├─ Amount: %s\n", balance.UiTokenAmount.UiAmountString) // Changed from UiTokenAccount
	fmt.Printf("│  │  │  └─ Decimals: %d\n", balance.UiTokenAmount.Decimals)     // Changed from UiTokenAccount
}

func printInnerInstructions(instructions []*pb.InnerInstructions) {
	if len(instructions) == 0 {
		return
	}

	fmt.Println("│  ├─ Inner Instructions:")
	for _, inner := range instructions {
		fmt.Printf("│  │  ├─ Index: %d\n", inner.Index)
		for i, ix := range inner.Instructions {
			if ix.Instruction != nil {
				fmt.Printf("│  │  │  ├─ Inner Instruction %d:\n", i)
				fmt.Printf("│  │  │  │  ├─ Program ID Index: %d\n", ix.Instruction.ProgramIdIndex)
				fmt.Printf("│  │  │  │  ├─ Account Indexes: %v\n", ix.Instruction.Accounts)
				fmt.Printf("│  │  │  │  └─ Data: %s\n", base58.Encode(ix.Instruction.Data))
			}
			if ix.StackHeight != nil {
				fmt.Printf("│  │  │  │  └─ Stack Height: %d\n", *ix.StackHeight)
			}
		}
	}
}

func printLogMessages(messages []string) {
	if len(messages) == 0 {
		return
	}

	fmt.Println("│  ├─ Log Messages:")
	for i, log_ := range messages {
		fmt.Printf("│  │  ├─ [%d] %s\n", i, log_)
	}
}

func printRewards(rewards []*pb.Reward) {
	if len(rewards) == 0 {
		return
	}

	fmt.Println("│  └─ Rewards:")
	for _, reward := range rewards {
		fmt.Printf("│     ├─ Pubkey: %s\n", reward.Pubkey)
		fmt.Printf("│     ├─ Lamports: %d\n", reward.Lamports)
		fmt.Printf("│     ├─ Post Balance: %s SOL\n", lamportsToSol(reward.PostBalance))
		fmt.Printf("│     ├─ Reward Type: %d\n", reward.RewardType)
		fmt.Printf("│     └─ Commission: %d%%\n", reward.Commission)
	}
}

// Helper functions
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
