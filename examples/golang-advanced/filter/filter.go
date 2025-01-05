package filter

import (
	"strings"
	"sync"

	pb "example/proto"
	"example/types"
	"github.com/mr-tron/base58"
)

// NewFilter creates a new transaction filter
func NewFilter(config *types.Config) *Filter {
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

// Filter struct moved from types to here
type Filter struct {
	programFilters map[string]bool
	includeVote    bool
	includeFailed  bool
	mu             sync.RWMutex
}

// WantsTransaction checks if a transaction matches the filter criteria
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
	programIDs := ExtractProgramIDs(tx)

	// Check if any program in the transaction matches our filters
	for _, programID := range programIDs {
		if f.programFilters[programID] {
			return true
		}
	}

	return false
}

// ExtractProgramIDs extracts all program IDs from a transaction
func ExtractProgramIDs(tx *pb.TransactionEvent) []string {
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
