package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// admin represents an administrator within the SPV Wallet ecosystem.
// It includes the administrator's private key (xPriv) and provides access
// to the SPV Wallet's AdminAPI client for managing xPub and paymail-related operations.
type admin struct {
	xPriv   string           // The extended private key for the administrator.
	client  *wallet.AdminAPI // The API client for interacting with administrative functionalities in the SPV Wallet.
	paymail string           // The paymail addresses the administrator.
	alias   string           // The alias of the administrator.
}

// setPaymail sets the admin's Paymail address to the given value.
func (a *admin) setPaymail(s string) { a.paymail = a.alias + "@" + s }

// getAccessKeysAdmin fetches all access keys for the admin.
// It accepts a context and returns a slice of access keys and an error.
// If the operation fails, the error is non-nil and contains details of the failure.
func (a *admin) getAccessKeysAdmin(ctx context.Context) ([]*response.AccessKey, error) {
	keys, err := a.client.AccessKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin access keys: %w", err)
	}
	return keys.Content, nil
}

// getContacts retrieves all contacts for the admin.
// It accepts a context and returns a slice of contacts and an error.
// If the operation fails, the error is non-nil and contains details of the failure.
func (a *admin) getContacts(ctx context.Context) ([]*response.Contact, error) {
	contactsPage, err := a.client.Contacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin contacts: %w", err)
	}
	return contactsPage.Content, nil
}

// updateContact updates an existing contact in the admin panel.
// It accepts a context, contact ID, and full name as input parameters.
// On success, it returns the updated contact and a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) updateContact(ctx context.Context, contactID, fullName string) (*response.Contact, error) {
	cmd := &commands.UpdateContact{
		ID:       contactID,
		FullName: fullName,
	}

	updatedContact, err := a.client.ContactUpdate(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to update contact with ID %s: %w", contactID, err)
	}
	return updatedContact, nil
}

// deleteContact removes a contact from the admin panel.
// It accepts a context and contact ID as input parameters.
// On success, it returns a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) deleteContact(ctx context.Context, contactID string) error {
	if err := a.client.DeleteContact(ctx, contactID); err != nil {
		return fmt.Errorf("failed to delete contact with ID %s: %w", contactID, err)
	}
	return nil
}

// acceptContactInvitation accepts an invitation to add a contact.
// It accepts a context and invitation ID as input parameters.
// On success, it returns a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) acceptContactInvitation(ctx context.Context, invitationID string) error {
	if err := a.client.AcceptInvitation(ctx, invitationID); err != nil {
		return fmt.Errorf("failed to accept contact invitation with ID %s: %w", invitationID, err)
	}
	return nil
}

// rejectContactInvitation rejects an invitation to add a contact.
// It accepts a context and invitation ID as input parameters.
// On success, it returns a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) rejectContactInvitation(ctx context.Context, invitationID string) error {
	if err := a.client.RejectInvitation(ctx, invitationID); err != nil {
		return fmt.Errorf("failed to reject contact invitation with ID %s: %w", invitationID, err)
	}
	return nil
}

// createContact adds a new contact in the admin panel.
// It accepts a context, paymail, and full name as input parameters.
// On success, it returns the created contact and a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) createContact(ctx context.Context, paymail, fullName string) (*response.Contact, error) {
	cmd := &commands.CreateContact{
		Paymail:        paymail,
		FullName:       fullName,
		CreatorPaymail: paymail,
	}

	contact, err := a.client.CreateContact(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact for paymail %s: %w", paymail, err)
	}
	return contact, nil
}

// confirmContacts confirms a contact connection between two paymails.
// It accepts a context, paymails A and B as input parameters.
// On success, it returns a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) confirmContacts(ctx context.Context, paymailA, paymailB string) error {
	cmd := &commands.ConfirmContacts{
		PaymailA: paymailA,
		PaymailB: paymailB,
	}

	if err := a.client.ConfirmContacts(ctx, cmd); err != nil {
		return fmt.Errorf("failed to confirm contact between %s and %s: %w", paymailA, paymailB, err)
	}
	return nil
}

// unconfirmContact unconfirms a contact connection between two paymails.
// It accepts a context, paymails A and B as input parameters.
// On success, it returns a nil error.
// If the operation fails, it returns a non-nil error with details of the failure.
func (a *admin) unconfirmContact(ctx context.Context, contactID string) error {
	if err := a.client.UnconfirmContact(ctx, contactID); err != nil {
		return fmt.Errorf("failed to unconfirm contact for contact id [%s]: %w", contactID, err)
	}
	return nil
}

// initAdmin initializes a new admin within the SPV Wallet ecosystem.
// It accepts the SPV Wallet API URL and the administrator's extended private key (xPriv) as input parameters.
// The function initializes the wallet's AdminAPI client using the provided xPriv,
// enabling the management of xPub and paymail-related operations.
// On success, it returns the initialized admin and a nil error.
// If the initialization fails, it returns a non-nil error with details of the failure.
func initAdmin(url, xPriv string) (*admin, error) {
	cfg := config.New(config.WithAddr(url))
	client, err := wallet.NewAdminAPIWithXPriv(cfg, xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin API: %w", err)
	}

	return &admin{xPriv: xPriv, client: client, alias: "Admin"}, nil
}
