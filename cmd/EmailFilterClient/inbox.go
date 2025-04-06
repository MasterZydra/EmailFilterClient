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

	// Note: c.Move() cannot be used here because it would change the seqNum

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

func MoveMessageToNewsletter(c *client.Client, state *State, seqNum uint32) error {
	if !state.HasNewsletterMailbox {
		return nil
	}

	newsletterFolder := "Newsletter"
	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNum)

	// Note: c.Move() cannot be used here because it would change the seqNum

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

func HasNewsletterMailbox(c *client.Client) bool {
	mailboxes := make(chan *imap.MailboxInfo, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "Newsletter", mailboxes)
	}()

	hasNewsletterMailbox := false
	for range mailboxes {
		hasNewsletterMailbox = true
	}

	if err := <-done; err != nil {
		return false
	}

	return hasNewsletterMailbox
}
