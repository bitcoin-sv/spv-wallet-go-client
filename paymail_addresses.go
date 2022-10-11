package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// NewPaymail will create a new paymail
func (b *BuxClient) NewPaymail(ctx context.Context, rawXPub, paymailAddress, avatar, publicName string, metadata *bux.Metadata) error {
	return b.transport.NewPaymail(ctx, rawXPub, paymailAddress, avatar, publicName, metadata)
}
