package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// GetAccessKey gets the access key given by id
func (b *BuxClient) GetAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	return b.transport.GetAccessKey(ctx, id)
}

// GetAccessKeys gets all the access keys filtered by the metadata
func (b *BuxClient) GetAccessKeys(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.AccessKey, error) {
	return b.transport.GetAccessKeys(ctx, metadataConditions)
}

// CreateAccessKey creates new access key
func (b *BuxClient) CreateAccessKey(ctx context.Context, metadata *bux.Metadata) (*bux.AccessKey, error) {
	return b.transport.CreateAccessKey(ctx, metadata)
}

// RevokeAccessKey revoked the access key given by id
func (b *BuxClient) RevokeAccessKey(ctx context.Context, id string) (*bux.AccessKey, error) {
	return b.transport.RevokeAccessKey(ctx, id)
}
