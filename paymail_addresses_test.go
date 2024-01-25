package buxclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BuxOrg/go-buxclient/fixtures"
)

// TestPaymailAddresses will test Paymail Addresses methods
func TestPaymailAddresses(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/paymail",
		Result:    "null",
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("NewPaymail", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, false)

		// when
		err := client.NewPaymail(context.Background(), fixtures.XPubString, fixtures.PaymailAddress, "", "", fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
	})
}
