package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"example/types"
)

// LoadConfig loads the client configuration from a file
func LoadConfig(filename string) (*types.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config types.Config
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

// ValidateConfig validates the configuration
func ValidateConfig(config *types.Config) error {
	if config.ServerAddress == "" {
		return fmt.Errorf("server address is required")
	}
	if config.AuthToken == "" {
		return fmt.Errorf("auth token is required")
	}
	return nil
}

// GetLogDirectory returns the absolute path to the log directory
func GetLogDirectory(config *types.Config) (string, error) {
	if config.LogDirectory == "" {
		return os.Getwd()
	}
	return filepath.Abs(config.LogDirectory)
}
