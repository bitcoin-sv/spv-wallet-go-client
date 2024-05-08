package walletclient

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/models"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

func (wc *WalletClient) AdminGetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *transports.QueryParams) ([]*models.Contact, transports.ResponseError) {
	return wc.transport.AdminGetContacts(ctx, conditions, metadata, queryParams)
}

func (wc *WalletClient) AdminUpdateContact(ctx context.Context, id, fullName string, metadata *models.Metadata) (*models.Contact, transports.ResponseError) {
	return wc.transport.AdminUpdateContact(ctx, id, fullName, metadata)
}

func (wc *WalletClient) AdminDeleteContact(ctx context.Context, id string) transports.ResponseError {
	return wc.transport.AdminDeleteContact(ctx, id)
}

func (wc *WalletClient) AdminAcceptContact(ctx context.Context, id string) (*models.Contact, transports.ResponseError) {
	return wc.transport.AdminAcceptContact(ctx, id)
}

func (wc *WalletClient) AdminRejectContact(ctx context.Context, id string) (*models.Contact, transports.ResponseError) {
	return wc.transport.AdminRejectContact(ctx, id)
}
