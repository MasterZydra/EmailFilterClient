package main

import (
	"fmt"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// MoveMessageToTrash moves a message to the Trash folder
func MoveMessageToTrash(c *client.Client, seqNum uint32) {
	trashFolder := "Trash"
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	// Copy the message to the Trash folder
	if err := c.Copy(seqset, trashFolder); err != nil {
		fmt.Println("Failed to move message to Trash:", err)
		return
	}

	// Mark the message as deleted in the current folder
	if err := c.Store(seqset, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.DeletedFlag}, nil); err != nil {
		fmt.Println("Failed to mark message as deleted:", err)
	}
}
