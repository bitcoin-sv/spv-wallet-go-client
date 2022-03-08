package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// RegisterXpub registers a new xpub - admin key needed
func (b *BuxClient) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {
	return b.transport.RegisterXpub(ctx, rawXPub, metadata)
}

// GetXPub gets the current xpub
func (b *BuxClient) GetXPub(ctx context.Context) (*bux.Xpub, error) {
	return b.transport.GetXPub(ctx)
}
