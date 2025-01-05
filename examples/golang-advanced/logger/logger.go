package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mr-tron/base58"
)

// SignatureLogger handles signature logging
type SignatureLogger struct {
	file *os.File
	mu   sync.Mutex
}

// NewSignatureLogger creates a new signature logger
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

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	Slot      uint64 `json:"slot"`
	ProgramID string `json:"program_id"`
	Success   bool   `json:"success"`
}

// LogSignature logs a transaction signature
func (sl *SignatureLogger) LogSignature(signature []byte, slot uint64, programID string, success bool) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	entry := LogEntry{
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

// Close closes the signature logger
func (sl *SignatureLogger) Close() error {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.file.Close()
}
