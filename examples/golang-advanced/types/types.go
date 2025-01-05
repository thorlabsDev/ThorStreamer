package types

import (
	"os"
	"sync"
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
	ProgramFilters map[string]bool
	IncludeVote    bool
	IncludeFailed  bool
	mu             sync.RWMutex
}

// SignatureLogger handles signature logging
type SignatureLogger struct {
	file *os.File
	mu   sync.Mutex
}

// Constants
const (
	LamportsPerSol = 1000000000 // SOL to lamports conversion rate
)
