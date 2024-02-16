package walletclient

import (
	"context"
	"testing"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
)

// TestDestinations will test the Destinations methods
func TestDestinations(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/destination",
		Result:    fixtures.MarshallForTestHandler(fixtures.Destination),
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("GetDestinationByID", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.GetDestinationByID(context.Background(), fixtures.Destination.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("GetDestinationByAddress", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.GetDestinationByAddress(context.Background(), fixtures.Destination.Address)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("GetDestinationByLockingScript", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.GetDestinationByLockingScript(context.Background(), fixtures.Destination.LockingScript)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("GetDestinations", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/destination/search",
			Result:    fixtures.MarshallForTestHandler([]*buxmodels.Destination{fixtures.Destination}),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestWalletClient(transportHandler, false)

		// when
		destinations, err := client.GetDestinations(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destinations, []*buxmodels.Destination{fixtures.Destination})
	})

	t.Run("NewDestination", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.NewDestination(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("UpdateDestinationMetadataByID", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.UpdateDestinationMetadataByID(context.Background(), fixtures.Destination.ID, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("UpdateDestinationMetadataByAddress", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.UpdateDestinationMetadataByAddress(context.Background(), fixtures.Destination.Address, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})

	t.Run("UpdateDestinationMetadataByLockingScript", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, false)

		// when
		destination, err := client.UpdateDestinationMetadataByLockingScript(context.Background(), fixtures.Destination.LockingScript, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, destination, fixtures.Destination)
	})
}
