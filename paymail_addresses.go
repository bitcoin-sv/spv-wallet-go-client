package walletclient

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// NewPaymail will create a new paymail
//
// paymailAddress: The paymail address to create (e.g., example@bux.org)
func (b *WalletClient) NewPaymail(ctx context.Context, rawXPub, paymailAddress, avatar, publicName string, metadata *models.Metadata) transports.ResponseError {
	return b.transport.NewPaymail(ctx, rawXPub, paymailAddress, avatar, publicName, metadata)
}
