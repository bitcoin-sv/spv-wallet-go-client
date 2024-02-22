package walletclient

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// NewXpub registers a new xpub - admin key needed
func (b *WalletClient) NewXpub(ctx context.Context, rawXPub string, metadata *models.Metadata) transports.ResponseError {
	return b.transport.NewXpub(ctx, rawXPub, metadata)
}

// GetXPub gets the current xpub
func (b *WalletClient) GetXPub(ctx context.Context) (*models.Xpub, transports.ResponseError) {
	return b.transport.GetXPub(ctx)
}

// UpdateXPubMetadata update the metadata of the logged in xpub
func (b *WalletClient) UpdateXPubMetadata(ctx context.Context, metadata *models.Metadata) (*models.Xpub, transports.ResponseError) {
	return b.transport.UpdateXPubMetadata(ctx, metadata)
}
