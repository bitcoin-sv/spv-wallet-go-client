package walletclient

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContactActionsRouting will test routing
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
		{
			name:            "ConfirmContact",
			route:           "/contact/confirmed/",
			responsePayload: "{}",
			f: func(c *WalletClient) error {
				contact := models.Contact{PubKey: fixtures.PubKey}

				passcode, _ := c.GenerateTotpForContact(&contact, 30, 2)
				return c.ConfirmContact(context.Background(), &contact, passcode, 30, 2)
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

func TestConfirmContact(t *testing.T) {
	t.Run("TOTP is valid - call Confirm Action", func(t *testing.T) {
		// given
		tmq := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/contact/confirmed/",
			Result:    "{}",
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}

		client := getTestWalletClientWithOpts(tmq, WithXPriv(fixtures.XPrivString))

		contact := &models.Contact{PubKey: fixtures.PubKey}
		totp, err := client.GenerateTotpForContact(contact, 30, 2)
		require.NoError(t, err)

		// when
		result := client.ConfirmContact(context.Background(), contact, totp, 30, 2)

		// then
		require.Nil(t, result)
	})

	t.Run("TOTP is invalid -  do not call Confirm Action", func(t *testing.T) {
		// given
		tmq := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/unknown/",
			Result:    "{}",
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}

		client := getTestWalletClientWithOpts(tmq, WithXPriv(fixtures.XPrivString))

		alice := &models.Contact{PubKey: "034252e5359a1de3b8ec08e6c29b80594e88fb47e6ae9ce65ee5a94f0d371d2cde"}
		a_totp, err := client.GenerateTotpForContact(alice, 30, 2)
		require.NoError(t, err)

		bob := &models.Contact{PubKey: "02dde493752f7bc89822ed8a13e0e4aa04550c6c4430800e4be1e5e5c2556cf65b"}

		// when
		result := client.ConfirmContact(context.Background(), bob, a_totp, 30, 2)

		// then
		require.NotNil(t, result)
		require.Equal(t, result.Error(), "totp is invalid")
	})
}
