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
