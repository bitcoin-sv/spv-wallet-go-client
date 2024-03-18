package walletclient

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/stretchr/testify/assert"
)

// TestAcceptContact will test the AcceptContact method
func TestAcceptContact(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/contact/accept/",
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("AcceptContact", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, true)

		// when
		err := client.AcceptContact(context.Background(), fixtures.PaymailAddress)

		// then
		assert.NoError(t, err)
	})
}

// TestRejectContact will test the RejectContact method
func TestRejectContact(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/contact/reject/",
		Result:    fixtures.MarshallForTestHandler(""),
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("RejectContact", func(t *testing.T) {
		// given
		client := getTestWalletClient(transportHandler, true)

		// when
		err := client.RejectContact(context.Background(), fixtures.PaymailAddress)

		// then
		assert.NoError(t, err)
	})
}
