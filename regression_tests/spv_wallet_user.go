package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

const (
	// TOTP_DIGITS represents the number of digits in the TOTP.
	TOTP_DIGITS uint = 2
	// TOTP_PERIOD represents the period for the TOTP.
	TOTP_PERIOD uint = 1200
)

// transactionsSlice represents a slice of response.Transaction objects.
type transactionsSlice []*response.Transaction

// Has checks if a transaction with the specified ID exists in the transactions slice.
// It returns true if a transaction with the given ID is found, and false otherwise.
func (tt transactionsSlice) Has(id string) bool {
	for _, t := range tt {
		if t.ID == id {
			return true
		}
	}

	return false
}

// user represents an individual user within the SPV Wallet ecosystem.
// It includes details like the alias, private key (xPriv), public key (xPub), and paymail address.
// The user struct also utilizes the wallet's UserAPI client to interact with the SPV Wallet API
// for transaction-related operations.
type user struct {
	alias     string          // The unique alias for the user.
	xPriv     string          // The extended private key for the user.
	xPub      string          // The extended public key for the user.
	paymail   string          // The paymail address associated with the user.
	paymailID string          // The paymail id associated with the users paymail.
	client    *wallet.UserAPI // The API client for interacting with the SPV Wallet.
}

// setPaymail sets the user's Paymail address with the given domain.
func (u *user) setPaymail(domain string) { u.paymail = u.alias + "@" + domain }

// transferFunds sends a specified amount of satoshis to a recipient's paymail.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the actor has insufficient funds or the API call fails, it returns a non-nil error.
func (u *user) transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := u.balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("balance failed: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d available, %d required", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := u.client.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds to %s: %w", paymail, err)
	}

	return transaction, nil
}

// balance retrieves the current satoshi balance for given actor.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the current balance and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) balance(ctx context.Context) (uint64, error) {
	xPub, err := u.client.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve xPub: %w", err)
	}

	return xPub.CurrentBalance, nil
}

// addContact adds a new contact to the user's contact list.
// It accepts the contact's paymail, contact full name, and the requester's paymail as input parameters.
// The function uses the SPV Wallet API to add the contact to the user's contact list.
// On success, it returns the newly added contact and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) addContact(ctx context.Context, contactPaymail, contactFullName string) (*response.Contact, error) {
	resp, err := u.client.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   contactPaymail,
		FullName:         contactFullName,
		RequesterPaymail: u.paymail,
	})
	if err != nil {
		return nil, fmt.Errorf("could not add contact %s: %w", contactPaymail, err)
	}

	return resp, nil
}

// getContact retrieves the contact details for a given paymail address.
// It accepts the contact's paymail as input parameter.
// The function uses the SPV Wallet API to fetch the contact details for the given paymail.
// On success, it returns the contact details and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) getContact(ctx context.Context, contactPaymail string) (*response.Contact, error) {
	resp, err := u.client.ContactWithPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, fmt.Errorf("could not get contact %s: %w", contactPaymail, err)
	}

	return resp, nil
}

// confirmContact confirms a contact using a received TOTP.
// It accepts the contact's paymail, received TOTP, and requester's paymail as input parameters.
// The function uses the SPV Wallet API to confirm the contact using the received TOTP.
// On success, it returns a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) confirmContact(ctx context.Context, contactPaymail, receivedTotp string) error {
	contactResp, err := u.client.ContactWithPaymail(ctx, contactPaymail)
	if err != nil {
		return fmt.Errorf("failed to fetch contact %s: %w", contactPaymail, err)
	}
	if contactResp == nil {
		return fmt.Errorf("contact %s not found", contactPaymail)
	}

	contact := mapToContactModel(contactResp)
	err = u.client.ConfirmContact(ctx, contact, receivedTotp, u.paymail, TOTP_PERIOD, TOTP_DIGITS)
	if err != nil {
		return fmt.Errorf("failed to confirm contact %s: %w", contactPaymail, err)
	}

	return nil
}

// unconfirmContact unconfirms a contact by their Paymail.
// It accepts the contact's paymail as input parameter.
// The function uses the SPV Wallet API to unconfirm the contact from the user's contact list.
// On success, it returns a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) unconfirmContact(ctx context.Context, contactPaymail string) error {
	err := u.client.UnconfirmContact(ctx, contactPaymail)
	if err != nil {
		return fmt.Errorf("failed to unconfirm contact %s: %w", contactPaymail, err)
	}
	return nil
}

// removeContact removes a contact by their Paymail.
// It accepts the contact's paymail as input parameter.
// The function uses the SPV Wallet API to remove the contact from the user's contact list.
// On success, it returns a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) removeContact(ctx context.Context, contactPaymail string) error {
	err := u.client.RemoveContact(ctx, contactPaymail)
	if err != nil {
		return fmt.Errorf("failed to remove contact %s: %w", contactPaymail, err)
	}
	return nil
}

// generateTotp generates a TOTP for a contact.
// It accepts the contact's paymail as input parameter.
// The function uses the SPV Wallet API to generate a TOTP for the given contact.
// On success, it returns the generated TOTP and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) generateTotp(ctx context.Context, contactPaymail string) (string, error) {
	contactResp, err := u.client.ContactWithPaymail(ctx, contactPaymail)
	if err != nil {
		return "", fmt.Errorf("failed to fetch contact %s: %w", contactPaymail, err)
	}
	if contactResp == nil {
		return "", fmt.Errorf("contact %s not found", contactPaymail)
	}

	contact := mapToContactModel(contactResp)
	totp, err := u.client.GenerateTotpForContact(contact, TOTP_PERIOD, TOTP_DIGITS)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP for contact %s: %w", contactPaymail, err)
	}

	return totp, nil
}

// getAccessKeys fetches all access keys for the user.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of access keys and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) getAccessKeys(ctx context.Context) ([]*response.AccessKey, error) {
	keys, err := u.client.AccessKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user access keys: %w", err)
	}
	return keys.Content, nil
}

// generateAccessKey generates a new access key for the user.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the generated access key and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) generateAccessKey(ctx context.Context) (*response.AccessKey, error) {
	key, err := u.client.GenerateAccessKey(ctx, &commands.GenerateAccessKey{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate access key: %w", err)
	}
	return key, nil
}

// getAccessKeyByID retrieves a specific access key by ID.
// It accepts the access key ID as input parameter.
// The function uses the SPV Wallet API to fetch the access key with the given ID.
// On success, it returns the access key and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) getAccessKeyByID(ctx context.Context, accessKeyID string) (*response.AccessKey, error) {
	key, err := u.client.AccessKey(ctx, accessKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access key with ID %s: %w", accessKeyID, err)
	}
	return key, nil
}

// revokeAccessKey revokes an access key.
// It accepts the access key ID as input parameter.
// The function uses the SPV Wallet API to revoke the access key with the given ID.
// On success, it returns a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) revokeAccessKey(ctx context.Context, accessKeyID string) error {
	err := u.client.RevokeAccessKey(ctx, accessKeyID)
	if err != nil {
		return fmt.Errorf("failed to revoke access key with ID %s: %w", accessKeyID, err)
	}
	return nil
}

// initUser initializes a new user within the SPV Wallet ecosystem.
// It accepts the alias and SPV Wallet API URL as input parameters.
// The function generates a random pair of wallet keys (xPub, xPriv) and uses the xPriv key
// to initialize the wallet's client, enabling transaction-related operations.
// On success, it returns the initialized user and a nil error.
// If user initialization fails, it returns a non-nil error with details of the failure.
func initUser(alias, url string) (*user, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("could not generate random keys: %w", err)
	}

	client, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), keys.XPriv())
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API for alias %q: %w", alias, err)
	}

	return &user{
		alias:  alias,
		xPriv:  keys.XPriv(),
		xPub:   keys.XPub(),
		client: client,
	}, nil
}

// initUserWithXPriv initializes a new user within the SPV Wallet ecosystem.
// It accepts the alias, xPriv and SPV Wallet API URL as input parameters.
// The function nitializes the wallet's client, enabling transaction-related operations.
// On success, it returns the initialized user and a nil error.
// If user initialization fails, it returns a non-nil error with details of the failure.
func initUserWithXPriv(alias, url, xPriv string) (*user, error) {
	client, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API for alias %q: %w", alias, err)
	}

	xPub, err := walletkeys.XPubFromXPriv(xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not get xPub from xPriv: %w", err)
	}

	return &user{
		alias:  alias,
		xPriv:  xPriv,
		xPub:   xPub,
		client: client,
	}, nil
}

// mapToContactModel maps a response.Contact object to a models.Contact object.
// It accepts a response.Contact object as input parameter.
// On success, it returns the mapped models.Contact object.
// If the input parameter is nil, it returns nil.
func mapToContactModel(resp *response.Contact) *models.Contact {
	if resp == nil {
		return nil
	}
	return &models.Contact{
		ID:       resp.ID,
		FullName: resp.FullName,
		Paymail:  resp.Paymail,
		PubKey:   resp.PubKey,
		Status:   resp.Status,
	}
}
