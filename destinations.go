package buxclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
)

// GetDestinationByID gets the destination by id
func (b *BuxClient) GetDestinationByID(ctx context.Context, id string) (*buxmodels.Destination, error) {
	return b.transport.GetDestinationByID(ctx, id)
}

// GetDestinationByAddress gets the destination by address
func (b *BuxClient) GetDestinationByAddress(ctx context.Context, address string) (*buxmodels.Destination, error) {
	return b.transport.GetDestinationByAddress(ctx, address)
}

// GetDestinationByLockingScript gets the destination by locking script
func (b *BuxClient) GetDestinationByLockingScript(ctx context.Context,
	lockingScript string,
) (*buxmodels.Destination, error) {
	return b.transport.GetDestinationByLockingScript(ctx, lockingScript)
}

// GetDestinations gets all destinations that match the metadata filter
func (b *BuxClient) GetDestinations(ctx context.Context,
	metadataConditions *buxmodels.Metadata,
) ([]*buxmodels.Destination, error) {
	return b.transport.GetDestinations(ctx, metadataConditions)
}

// NewDestination create a new destination and return it
func (b *BuxClient) NewDestination(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Destination, error) {
	return b.transport.NewDestination(ctx, metadata)
}

// UpdateDestinationMetadataByID updates the destination metadata by id
func (b *BuxClient) UpdateDestinationMetadataByID(ctx context.Context, id string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, error) {
	return b.transport.UpdateDestinationMetadataByID(ctx, id, metadata)
}

// UpdateDestinationMetadataByAddress updates the destination metadata by address
func (b *BuxClient) UpdateDestinationMetadataByAddress(ctx context.Context, address string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, error) {
	return b.transport.UpdateDestinationMetadataByAddress(ctx, address, metadata)
}

// UpdateDestinationMetadataByLockingScript updates the destination metadata by locking script
func (b *BuxClient) UpdateDestinationMetadataByLockingScript(ctx context.Context, lockingScript string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Destination, error) {
	return b.transport.UpdateDestinationMetadataByLockingScript(ctx, lockingScript, metadata)
}
