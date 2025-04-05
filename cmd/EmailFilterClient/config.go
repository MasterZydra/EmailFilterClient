package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the structure of the config.json file
type Config struct {
	Interval         int               `json:"interval"`
	IMAP_Connections []IMAP_Connection `json:"imapConnections"`
}

// IMAP_Connection represents a single IMAP connection configuration
type IMAP_Connection struct {
	Host     string `json:"host"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ReadConfig reads and parses the config.json file
func ReadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}
