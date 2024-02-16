package walletclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// GetAccessKey gets the access key given by id
func (b *WalletClient) GetAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.GetAccessKey(ctx, id)
}

// GetAccessKeys gets all the access keys filtered by the metadata
func (b *WalletClient) GetAccessKeys(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.GetAccessKeys(ctx, metadataConditions)
}

// CreateAccessKey creates new access key
func (b *WalletClient) CreateAccessKey(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.CreateAccessKey(ctx, metadata)
}

// RevokeAccessKey revoked the access key given by id
func (b *WalletClient) RevokeAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.RevokeAccessKey(ctx, id)
}
