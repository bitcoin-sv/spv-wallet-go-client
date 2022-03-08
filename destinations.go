package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// GetDestinationByID gets the destination by id
func (b *BuxClient) GetDestinationByID(ctx context.Context, id string) (*bux.Destination, error) {
	return b.transport.GetDestinationByID(ctx, id)
}

// GetDestinationByAddress gets the destination by address
func (b *BuxClient) GetDestinationByAddress(ctx context.Context, address string) (*bux.Destination, error) {
	return b.transport.GetDestinationByAddress(ctx, address)
}

// GetDestinationByLockingScript gets the destination by locking script
func (b *BuxClient) GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*bux.Destination, error) {
	return b.transport.GetDestinationByLockingScript(ctx, lockingScript)
}

// GetDestinations gets all destinations that match the metadata filter
func (b *BuxClient) GetDestinations(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.Destination, error) {
	return b.transport.GetDestinations(ctx, metadataConditions)
}

// NewDestination create a new destination and return it
func (b *BuxClient) NewDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	return b.transport.NewDestination(ctx, metadata)
}
