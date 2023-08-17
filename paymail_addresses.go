package buxclient

import (
	"context"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
)

// NewPaymail will create a new paymail
func (b *BuxClient) NewPaymail(ctx context.Context, rawXPub, paymailAddress, avatar, publicName string, metadata *buxmodels.Metadata) transports.ResponseError {
	return b.transport.NewPaymail(ctx, rawXPub, paymailAddress, avatar, publicName, metadata)
}
