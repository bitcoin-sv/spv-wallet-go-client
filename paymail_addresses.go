package walletclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// NewPaymail will create a new paymail
//
// paymailAddress: The paymail address to create (e.g., example@bux.org)
func (b *WalletClient) NewPaymail(ctx context.Context, rawXPub, paymailAddress, avatar, publicName string, metadata *buxmodels.Metadata) transports.ResponseError {
	return b.transport.NewPaymail(ctx, rawXPub, paymailAddress, avatar, publicName, metadata)
}
