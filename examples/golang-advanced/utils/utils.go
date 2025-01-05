package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example/types"
)

// LamportsToSol converts lamports to SOL
func LamportsToSol(lamports uint64) string {
	return fmt.Sprintf("%.9f", float64(lamports)/float64(types.LamportsPerSol))
}

// Prompt displays a message and returns user input
func Prompt(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// GetUserWallets prompts the user for wallet addresses
func GetUserWallets() ([]string, error) {
	var wallets []string
	fmt.Println("Enter up to 10 wallet addresses (one per line).")
	fmt.Println("Press Enter twice when done:")

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 10; i++ {
		fmt.Printf("Wallet %d (or blank to finish): ", i+1)
		wallet, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading input: %v", err)
		}

		wallet = strings.TrimSpace(wallet)
		if wallet == "" {
			break
		}
		wallets = append(wallets, wallet)
	}

	if len(wallets) == 0 {
		return nil, fmt.Errorf("no wallet addresses provided")
	}

	return wallets, nil
}

// Abs returns the absolute value of an int64
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// GetTransactionVersionString returns a string representation of a transaction version
func GetTransactionVersionString(version uint32) string {
	switch version {
	case 0:
		return "Legacy"
	case 1:
		return "V0"
	default:
		return fmt.Sprintf("Unknown Version(%d)", version)
	}
}

// FormatStatus formats a transaction status with optional error info
func FormatStatus(isErr bool, errorInfo string) string {
	if isErr {
		if errorInfo != "" {
			return fmt.Sprintf("Failed (%s)", errorInfo)
		}
		return "Failed"
	}
	return "Success"
}
