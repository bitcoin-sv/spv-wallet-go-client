package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
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
