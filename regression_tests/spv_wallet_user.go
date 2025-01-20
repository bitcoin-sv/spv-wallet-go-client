package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
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
	paymailID string          //The paymail id associated with the users paymail.
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

	return &user{
		alias:  alias,
		xPriv:  xPriv,
		client: client,
	}, nil
}
