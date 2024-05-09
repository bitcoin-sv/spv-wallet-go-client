package walletclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
)

func TestAdminContacts(t *testing.T) {
	tcs := []struct {
		name            string
		route           string
		responsePayload string
		f               func(c *WalletClient) error
	}{
		{
			name:            "AdminRejectContact",
			route:           "admin/contact/rejected/",
			responsePayload: "{}",
			f: func(c *WalletClient) error {
				_, err := c.AdminRejectContact(context.Background(), "admin")
				return err
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// given
			tmq := testTransportHandler{
				Type:      fixtures.RequestType,
				Path:      tc.route,
				Result:    tc.responsePayload,
				ClientURL: fixtures.ServerURL,
				Client:    WithHTTPClient,
			}

			client := getTestWalletClientWithOpts(tmq, WithXPriv(fixtures.XPrivString))

			// when
			err := tc.f(client)

			// then
			assert.NoError(t, err)
		})
	}

}
