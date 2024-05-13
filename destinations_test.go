package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
)

func TestDestinations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendJSONResponse := func(data interface{}) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(data); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		switch {
		case r.URL.Path == "/destination/address/"+fixtures.Destination.Address && r.Method == http.MethodGet:
			sendJSONResponse(fixtures.Destination)
		case r.URL.Path == "/destination/lockingScript/"+fixtures.Destination.LockingScript && r.Method == http.MethodGet:
			sendJSONResponse(fixtures.Destination)
		case r.URL.Path == "/destination/search" && r.Method == http.MethodPost:
			sendJSONResponse([]*models.Destination{fixtures.Destination})
		case r.URL.Path == "/destination" && r.Method == http.MethodGet:
			sendJSONResponse(fixtures.Destination)
		case r.URL.Path == "/destination" && r.Method == http.MethodPatch:
			sendJSONResponse(fixtures.Destination)
		case r.URL.Path == "/destination" && r.Method == http.MethodPost:
			sendJSONResponse(fixtures.Destination)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewWalletClientWithAccessKey(fixtures.AccessKeyString, server.URL, true)
	require.NoError(t, err)

	t.Run("GetDestinationByID", func(t *testing.T) {
		destination, err := client.GetDestinationByID(context.Background(), fixtures.Destination.ID)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("GetDestinationByAddress", func(t *testing.T) {
		destination, err := client.GetDestinationByAddress(context.Background(), fixtures.Destination.Address)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("GetDestinationByLockingScript", func(t *testing.T) {
		destination, err := client.GetDestinationByLockingScript(context.Background(), fixtures.Destination.LockingScript)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("GetDestinations", func(t *testing.T) {
		destinations, err := client.GetDestinations(context.Background(), fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, []*models.Destination{fixtures.Destination}, destinations)
	})

	t.Run("NewDestination", func(t *testing.T) {
		destination, err := client.NewDestination(context.Background(), fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("UpdateDestinationMetadataByID", func(t *testing.T) {
		destination, err := client.UpdateDestinationMetadataByID(context.Background(), fixtures.Destination.ID, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("UpdateDestinationMetadataByAddress", func(t *testing.T) {
		destination, err := client.UpdateDestinationMetadataByAddress(context.Background(), fixtures.Destination.Address, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})

	t.Run("UpdateDestinationMetadataByLockingScript", func(t *testing.T) {
		destination, err := client.UpdateDestinationMetadataByLockingScript(context.Background(), fixtures.Destination.LockingScript, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Destination, destination)
	})
}
