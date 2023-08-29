package buxclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
)

// GetAccessKey gets the access key given by id
func (b *BuxClient) GetAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.GetAccessKey(ctx, id)
}

// GetAccessKeys gets all the access keys filtered by the metadata
func (b *BuxClient) GetAccessKeys(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.GetAccessKeys(ctx, metadataConditions)
}

// CreateAccessKey creates new access key
func (b *BuxClient) CreateAccessKey(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.CreateAccessKey(ctx, metadata)
}

// RevokeAccessKey revoked the access key given by id
func (b *BuxClient) RevokeAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, transports.ResponseError) {
	return b.transport.RevokeAccessKey(ctx, id)
}
