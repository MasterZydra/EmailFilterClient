package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Blacklist represents the structure of the blacklist.json file
type Blacklist struct {
	From       []string `json:"from"`
	Newsletter []string `json:"newsletter"`
}

// ReadBlacklist reads and parses the blacklist.json file
func ReadBlacklist(filePath string) (*Blacklist, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open blacklist file: %w", err)
	}
	defer file.Close()

	var blacklist Blacklist
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&blacklist); err != nil {
		return nil, fmt.Errorf("failed to decode blacklist file: %w", err)
	}

	return &blacklist, nil
}

// IsBlacklisted checks if a given email address is in the blacklist
func IsInList(email string, list []string) bool {
	for _, entry := range list {
		if strings.HasSuffix(email, entry) {
			return true
		}
	}
	return false
}

// ComputeBlacklistHash computes a hash for the blacklist
func ComputeBlacklistHash(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
