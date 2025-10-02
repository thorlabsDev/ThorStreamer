package printer

import (
	"fmt"

	"github.com/mr-tron/base58"

	pb "example/proto"
	"example/utils"
)

const (
	lamportsPerSol = 1000000000 // SOL to lamports conversion rate
)

// PrintTransaction prints basic transaction information
func PrintTransaction(tx *pb.TransactionEvent) {
	fmt.Printf("\n📜 Transaction: %s\n", base58.Encode(tx.Signature))
	fmt.Printf("├─ Slot: %d\n", tx.Slot)
	if tx.Transaction != nil {
		fmt.Printf("├─ Version: %s\n",
			utils.GetTransactionVersionString(tx.Transaction.Message.Version))
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

// PrintDetailedTransaction prints detailed transaction information
func PrintDetailedTransaction(tx *pb.TransactionEvent) {
	fmt.Println("\n👛 Wallet Transaction Details:")
	fmt.Printf("├─ Signature: %s\n", base58.Encode(tx.Signature))
	fmt.Printf("├─ Slot: %d\n", tx.Slot)
	fmt.Printf("├─ Is Vote Transaction: %v\n", tx.IsVote)

	if tx.Transaction != nil {
		fmt.Printf("├─ Transaction Version: %s\n",
			utils.GetTransactionVersionString(tx.Transaction.Message.Version))
		printTransactionDetails(tx.Transaction)
	}

	if tx.TransactionStatusMeta != nil {
		printTransactionStatusMeta(tx.TransactionStatusMeta)
	}

	fmt.Println("└─ End Transaction\n")
}

// PrintAccountUpdate prints account update information (updated for new proto)
func PrintAccountUpdate(account *pb.SubscribeUpdateAccountInfo) {
	fmt.Println("\n💳 Account Update:")
	fmt.Printf("├─ Address: %s\n", base58.Encode(account.Pubkey))
	fmt.Printf("├─ Owner: %s\n", base58.Encode(account.Owner))
	fmt.Printf("├─ Balance: %s SOL\n", utils.LamportsToSol(account.Lamports))
	fmt.Printf("├─ Executable: %v\n", account.Executable)
	fmt.Printf("├─ Rent Epoch: %v\n", account.RentEpoch)
	fmt.Printf("├─ Write Version: %v\n", account.WriteVersion)

	if account.TxnSignature != nil {
		fmt.Printf("├─ Transaction Signature: %s\n", base58.Encode(account.TxnSignature))
	}

	if account.Slot != nil {
		fmt.Printf("├─ Slot: %d\n", account.Slot.Slot)
		fmt.Printf("├─ Parent: %d\n", account.Slot.Parent)
		fmt.Printf("├─ Status: %d\n", account.Slot.Status)
		if len(account.Slot.BlockHash) > 0 {
			fmt.Printf("├─ Block Hash: %s\n", base58.Encode(account.Slot.BlockHash))
		}
		fmt.Printf("└─ Block Height: %d\n\n", account.Slot.BlockHeight)
	} else {
		fmt.Println("└─ Slot: N/A\n")
	}
}

// PrintSlotStatus prints slot status information (updated for new proto)
func PrintSlotStatus(slot *pb.SlotStatusEvent) {
	fmt.Println("\n📢 Slot Status:")
	fmt.Printf("├─ Slot: %d\n", slot.Slot)
	fmt.Printf("├─ Parent: %d\n", slot.Parent)

	statusStr := "UNKNOWN"
	switch slot.Status {
	case 0:
		statusStr = "PROCESSED"
	case 1:
		statusStr = "CONFIRMED"
	case 2:
		statusStr = "ROOTED"
	case 3:
		statusStr = "FIRST_SHRED_RECEIVED"
	case 4:
		statusStr = "COMPLETED"
	case 5:
		statusStr = "CREATED_BANK"
	case 6:
		statusStr = "DEAD"
	}

	fmt.Printf("├─ Status: %s\n", statusStr)

	if len(slot.BlockHash) > 0 {
		fmt.Printf("├─ Block Hash: %s\n", base58.Encode(slot.BlockHash))
	}
	fmt.Printf("└─ Block Height: %d\n\n", slot.BlockHeight)
}

// Private helper functions for detailed printing

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

	if msg.Version == 1 {
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
	fmt.Printf("│  ├─ Status: %s\n", utils.FormatStatus(meta.IsStatusErr, meta.ErrorInfo))
	fmt.Printf("│  ├─ Fee: %s SOL\n", utils.LamportsToSol(meta.Fee))

	printBalanceChanges(meta)
	printTokenBalances(meta)
	printInnerInstructions(meta.InnerInstructions)
	printLogMessages(meta.LogMessages)
	printRewards(meta.Rewards)
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
				utils.LamportsToSol(meta.PreBalances[i]),
				utils.LamportsToSol(meta.PostBalances[i]),
				utils.LamportsToSol(uint64(utils.Abs(change))))
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
	if balance.UiTokenAmount == nil {
		return
	}

	fmt.Printf("│  │  ├─ %s-Token Balance:\n", prefix)
	fmt.Printf("│  │  │  ├─ Account Index: %d\n", balance.AccountIndex)
	fmt.Printf("│  │  │  ├─ Mint: %s\n", balance.Mint)
	fmt.Printf("│  │  │  ├─ Owner: %s\n", balance.Owner)
	fmt.Printf("│  │  │  ├─ Amount: %s\n", balance.UiTokenAmount.UiAmountString)
	fmt.Printf("│  │  │  └─ Decimals: %d\n", balance.UiTokenAmount.Decimals)
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

func lamportsToSol(lamports uint64) string {
	return fmt.Sprintf("%.9f", float64(lamports)/float64(lamportsPerSol))
}

// Helper functions
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
