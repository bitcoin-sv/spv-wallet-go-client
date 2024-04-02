package walletclient

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/stretchr/testify/assert"
)

// TestRejectContact will test the RejectContact method
func TestContactActionsRouting(t *testing.T) {
	tcs := []struct {
		name            string
		route           string
		responsePayload string
		f               func(c *WalletClient) error
	}{
		{
			name:            "RejectContact",
			route:           "/contact/rejected/",
			responsePayload: "{}",
			f:               func(c *WalletClient) error { return c.RejectContact(context.Background(), fixtures.PaymailAddress) },
		},
		{
			name:            "AcceptContact",
			route:           "/contact/accepted/",
			responsePayload: "{}",
			f:               func(c *WalletClient) error { return c.AcceptContact(context.Background(), fixtures.PaymailAddress) },
		},
		{
			name:            "ConfirmContact",
			route:           "/contact/confirmed/",
			responsePayload: "{}",
			f:               func(c *WalletClient) error { return c.ConfirmContact(context.Background(), fixtures.PaymailAddress) },
		},
		{
			name:            "GetContacts",
			route:           "/contact/search/",
			responsePayload: "[]",
			f: func(c *WalletClient) error {
				_, err := c.GetContacts(context.Background(), nil, nil, nil)
				return err
			},
		},
		{
			name:            "UpsertContact",
			route:           "/contact/",
			responsePayload: "{}",
			f: func(c *WalletClient) error {
				_, err := c.UpsertContact(context.Background(), "", "", nil)
				return err
			},
		},
		{
			name:            "UpsertContactForPaymail",
			route:           "/contact/",
			responsePayload: "{}",
			f: func(c *WalletClient) error {
				_, err := c.UpsertContactForPaymail(context.Background(), "", "", nil, "")
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

			client := getTestWalletClient(tmq, true)

			// when
			err := tc.f(client)

			// then
			assert.NoError(t, err)
		})
	}

}
