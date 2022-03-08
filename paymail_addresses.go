package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// RegisterPaymail registers a new paymail
func (b *BuxClient) RegisterPaymail(ctx context.Context, rawXPub, paymailAddress string, metadata *bux.Metadata) error {
	return b.transport.RegisterPaymail(ctx, rawXPub, paymailAddress, metadata)
}
