package buxclient

import (
	"context"
	"testing"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/stretchr/testify/assert"

	"github.com/BuxOrg/go-buxclient/fixtures"
)

// TestAccessKeys will test the AccessKey methods
func TestAccessKeys(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/access-key",
		Result:    fixtures.MarshallForTestHandler(fixtures.AccessKey),
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("GetAccessKey", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		accessKey, err := client.GetAccessKey(context.Background(), fixtures.AccessKey.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, accessKey, fixtures.AccessKey)
	})

	t.Run("GetAccessKeys", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/access-key/search",
			Result:    fixtures.MarshallForTestHandler([]*buxmodels.AccessKey{fixtures.AccessKey}),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, true)

		// when
		accessKeys, err := client.GetAccessKeys(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, accessKeys, []*buxmodels.AccessKey{fixtures.AccessKey})
	})

	t.Run("CreateAccessKey", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		accessKey, err := client.CreateAccessKey(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, accessKey, fixtures.AccessKey)
	})

	t.Run("RevokeAccessKey", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		accessKey, err := client.RevokeAccessKey(context.Background(), fixtures.AccessKey.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, accessKey, fixtures.AccessKey)
	})
}
