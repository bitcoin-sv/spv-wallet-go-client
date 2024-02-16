package walletclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
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

	t.Run("GetXPub", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, true)

		// when
		xpub, err := client.GetXPub(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Xpub, xpub)
	})

	t.Run("UpdateXPubMetadata", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, true)

		// when
		xpub, err := client.UpdateXPubMetadata(context.Background(), fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Xpub, xpub)
	})
}
