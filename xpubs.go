package buxclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
)

// NewXpub registers a new xpub - admin key needed
func (b *BuxClient) NewXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) transports.ResponseError {
	return b.transport.NewXpub(ctx, rawXPub, metadata)
}

// GetXPub gets the current xpub
func (b *BuxClient) GetXPub(ctx context.Context) (*buxmodels.Xpub, transports.ResponseError) {
	return b.transport.GetXPub(ctx)
}

// UpdateXPubMetadata update the metadata of the logged in xpub
func (b *BuxClient) UpdateXPubMetadata(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Xpub, transports.ResponseError) {
	return b.transport.UpdateXPubMetadata(ctx, metadata)
}
