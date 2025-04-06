package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

var configJsonPath = "./config/config.json"
var blacklistJsonPath = "./config/blacklist.json"
var stateJsonPath = "./state.json"
var logFilePath = "./log/info.log"

var blacklistHash string

func main() {
	// Define a command-line flag for the port
	port := flag.String("port", "8080", "Port for the web server")
	basicAuthPassword := flag.String("basicAuthPassword", "", "Secret key to protect web server routes")
	flag.Parse()

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	// Set log output to the file
	log.SetOutput(logFile)

	// Start the web server in a separate goroutine
	go startWebServer(*port, *basicAuthPassword)

	for {
		startTime := time.Now()

		// Load configuration
		config, err := ReadConfig(configJsonPath)
		if err != nil {
			log.Fatalf("Error reading config: %v", err)
		}

		// Load blacklist
		blacklist, err := ReadBlacklist(blacklistJsonPath)
		if err != nil {
			log.Fatalf("Error reading blacklist: %v", err)
		}
		blacklistHash, err = ComputeBlacklistHash(blacklistJsonPath)
		if err != nil {
			log.Fatalf("Error computing blacklist hash: %v", err)
		}

		// Process each IMAP connection
		for _, connection := range config.IMAP_Connections {
			log.Println("")
			log.Printf("Processing inbox: %s\n", connection.Email)
			if err := ProcessIMAPConnection(connection, blacklist); err != nil {
				log.Fatalf("Error processing connection %s: %v\n", connection.Email, err)
			}
		}

		// Calculate the time to sleep until the next interval
		elapsed := time.Since(startTime)
		sleepDuration := time.Duration(config.Interval)*time.Minute - elapsed
		if sleepDuration > 0 {
			log.Println("")
			log.Printf("Sleeping for %v before the next iteration...\n", sleepDuration)
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

	log.Printf("Inbox for %s has %d messages\n", connection.Email, inbox.Messages)
	if inbox.Messages == 0 {
		log.Println("No messages. Skipping.")
		return nil
	}

	// Fetch and process messages
	return FetchAndProcessMessages(c, inbox.Messages, blacklist, connection.Email, state)
}

// FetchAndProcessMessages fetches and processes messages from the INBOX
func FetchAndProcessMessages(c *client.Client, totalMessages uint32, blacklist *Blacklist, email string, states *States) error {
	state := states.Find(email)
	state.HasNewsletterMailbox = HasNewsletterMailbox(c)
	from := state.SeqNumber
	to := totalMessages

	if state.BlacklistHash != blacklistHash {
		from = 0
		state.BlacklistHash = blacklistHash
	} else if from > to {
		log.Println("No new messages to process.")
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var deletedMsgs uint32
	var highestSeqNum uint32
	for msg := range messages {
		highestSeqNum = msg.SeqNum
		if ProcessMessage(c, msg, blacklist, state) {
			deletedMsgs++
		}
	}
	highestSeqNum -= deletedMsgs

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

	log.Println("Finished processing messages.")
	return nil
}

// ProcessMessage processes a single email message
func ProcessMessage(c *client.Client, msg *imap.Message, blacklist *Blacklist, state *State) bool {
	for _, from := range msg.Envelope.Sender {
		if IsInList(from.Address(), blacklist.From) {
			log.Printf("Moving message from %s to Trash\n", from.Address())
			if err := MoveMessageToTrash(c, msg.SeqNum); err != nil {
				log.Printf("Error moving message to trash: %v\n", err)
			}
			return true
		}

		if state.HasNewsletterMailbox && IsInList(from.Address(), blacklist.Newsletter) {
			log.Printf("Moving message from %s to Newsletter\n", from.Address())
			if err := MoveMessageToNewsletter(c, state, msg.SeqNum); err != nil {
				log.Printf("Error moving message to newsletter: %v\n", err)
			}
			return true
		}
	}
	return false
}
