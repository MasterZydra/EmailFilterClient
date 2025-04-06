package main

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// MoveMessageToTrash moves a message to the Trash folder
func MoveMessageToTrash(c *client.Client, seqNum uint32) error {
	trashFolder := "Trash"
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	// Copy the message to the Trash folder
	if err := c.Copy(seqset, trashFolder); err != nil {
		return fmt.Errorf("failed to move message to Trash: %w", err)
	}

	// Mark the message as deleted in the current folder
	if err := c.Store(seqset, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.DeletedFlag}, nil); err != nil {
		return fmt.Errorf("failed to mark message as deleted: %w", err)
	}

	return nil
}

func MoveMessageToNewsletter(c *client.Client, seqNum uint32) error {
	newsletterFolder := "Newsletter"
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	// Check if newsletter folder exists
	_, err := c.Select(newsletterFolder, false)
	if err != nil {
		// Create newsletter folder
		err = c.Create(newsletterFolder)
		if err != nil {
			return fmt.Errorf("failed to create mailbox \"%s\": %w", newsletterFolder, err)
		}
	}

	// Copy the message to the newsletter folder
	if err := c.Copy(seqset, newsletterFolder); err != nil {
		return fmt.Errorf("failed to move message to newsletter: %w", err)
	}

	// Mark the message as deleted in the current folder
	if err := c.Store(seqset, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.DeletedFlag}, nil); err != nil {
		return fmt.Errorf("failed to mark message as deleted: %w", err)
	}

	return nil
}
