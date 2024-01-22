package buxclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BuxOrg/go-buxclient/fixtures"
)

// TestXpub will test the Xpub methods
func TestXpub(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/xpub",
		Result:    fixtures.MarshallForTestHandler(fixtures.Xpub),
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("NewXpub", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		err := client.NewXpub(context.Background(), fixtures.XPubString, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
	})

	t.Run("GetXPub", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		xpub, err := client.GetXPub(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Xpub, xpub)
	})

	t.Run("UpdateXPubMetadata", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, true)

		// when
		xpub, err := client.UpdateXPubMetadata(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Xpub, xpub)
	})
}
