package main

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var configJson = "./config.json"
var blacklistJson = "./blacklist.json"
var stateJson = "./state.json"

func main() {
	for {
		startTime := time.Now()

		// Load configuration
		config, err := ReadConfig(configJson)
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}

		// Load blacklist
		blacklist, err := ReadBlacklist(blacklistJson)
		if err != nil {
			fmt.Println("Error reading blacklist:", err)
			return
		}

		// Process each IMAP connection
		for _, connection := range config.IMAP_Connections {
			fmt.Printf("Processing connection: %s\n", connection.Email)
			if err := ProcessIMAPConnection(connection, blacklist); err != nil {
				fmt.Printf("Error processing connection %s: %v\n", connection.Email, err)
			}
		}

		// Calculate the time to sleep until the next interval
		elapsed := time.Since(startTime)
		sleepDuration := time.Duration(config.Interval)*time.Minute - elapsed
		if sleepDuration > 0 {
			fmt.Printf("Sleeping for %v before the next iteration...\n", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}

// ProcessIMAPConnection handles the IMAP connection, fetching and processing emails
func ProcessIMAPConnection(connection IMAP_Connection, blacklist *Blacklist) error {
	// Load the state
	state, err := LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Connect to the IMAP server
	c, err := client.DialTLS(connection.Host, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer c.Logout()

	// Login to the IMAP server
	if err := c.Login(connection.Email, connection.Password); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Select the INBOX
	inbox, err := c.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	fmt.Printf("Inbox for %s has %d messages\n", connection.Email, inbox.Messages)
	if inbox.Messages == 0 {
		fmt.Println("No messages. Skipping.")
		return nil
	}

	// Fetch and process messages
	return FetchAndProcessMessages(c, inbox.Messages, blacklist, connection.Email, state)
}

// FetchAndProcessMessages fetches and processes messages from the INBOX
func FetchAndProcessMessages(c *client.Client, totalMessages uint32, blacklist *Blacklist, email string, states *States) error {
	state := states.Find(email)
	from := state.SeqNumber + 1
	to := totalMessages
	if from > to {
		fmt.Println("No new messages to process.")
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var highestSeqNum uint32
	for msg := range messages {
		highestSeqNum = msg.SeqNum
		ProcessMessage(c, msg, blacklist)
	}

	if err := <-done; err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	// Update the state with the highest sequence number processed
	if highestSeqNum > 0 {
		state.SeqNumber = highestSeqNum
		if err := SaveState(states); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
	}

	// Expunge messages marked as deleted
	if err := c.Expunge(nil); err != nil {
		return fmt.Errorf("failed to expunge messages: %w", err)
	}

	fmt.Println("Finished processing messages.")
	return nil
}

// ProcessMessage processes a single email message
func ProcessMessage(c *client.Client, msg *imap.Message, blacklist *Blacklist) {
	for _, from := range msg.Envelope.Sender {
		if IsBlacklisted(from.Address(), blacklist) {
			fmt.Printf("Moving message from %s to Trash\n", from.Address())
			MoveMessageToTrash(c, msg.SeqNum)
		}
	}
}
