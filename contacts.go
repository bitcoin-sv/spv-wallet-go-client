package walletclient

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// UpsertContact add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user in their contacts.
func (b *WalletClient) UpsertContact(ctx context.Context, paymail, fullName string, metadata *models.Metadata) (*models.Contact, transports.ResponseError) {
	return b.transport.UpsertContact(ctx, paymail, fullName, metadata, "")
}

// UpsertContactForPaymail add or update contact. When adding a new contact, the system utilizes Paymail's PIKE capability to dispatch an invitation request, asking the counterparty to include the current user specified paymail in their contacts.
func (b *WalletClient) UpsertContactForPaymail(ctx context.Context, paymail, fullName string, metadata *models.Metadata, requesterPaymail string) (*models.Contact, transports.ResponseError) {
	return b.transport.UpsertContact(ctx, paymail, fullName, metadata, requesterPaymail)
}

// AcceptContact will accept the contact associated with the paymail
func (b *WalletClient) AcceptContact(ctx context.Context, paymail string) transports.ResponseError {
	return b.transport.AcceptContact(ctx, paymail)
}

// RejectContact will reject the contact associated with the paymail
func (b *WalletClient) RejectContact(ctx context.Context, paymail string) transports.ResponseError {
	return b.transport.RejectContact(ctx, paymail)
}

// ConfirmContact will confirm the contact associated with the paymail
func (b *WalletClient) ConfirmContact(ctx context.Context, paymail string) transports.ResponseError {
	return b.transport.ConfirmContact(ctx, paymail)
}

// GetContacts will get contacts by conditions
func (b *WalletClient) GetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *transports.QueryParams) ([]*models.Contact, transports.ResponseError) {
	return b.transport.GetContacts(ctx, conditions, metadata, queryParams)
}