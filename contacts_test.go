package walletclient

// import (
// 	"context"
// 	"strings"
// 	"testing"
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// 

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
// )

// // TestContactActionsRouting will test routing
// func TestContactActionsRouting(t *testing.T) {
// 	tcs := []struct {
// 		name            string
// 		route           string
// 		responsePayload string
// 		f               func(c *WalletClient) error
// 	}{
// 		{
// 			name:            "RejectContact",
// 			route:           "/contact/rejected/",
// 			responsePayload: "{}",
// 			f:               func(c *WalletClient) error { return c.RejectContact(context.Background(), fixtures.PaymailAddress) },
// 		},
// 		{
// 			name:            "AcceptContact",
// 			route:           "/contact/accepted/",
// 			responsePayload: "{}",
// 			f:               func(c *WalletClient) error { return c.AcceptContact(context.Background(), fixtures.PaymailAddress) },
// 		},
// 		{
// 			name:            "GetContacts",
// 			route:           "/contact/search/",
// 			responsePayload: "[]",
// 			f: func(c *WalletClient) error {
// 				_, err := c.GetContacts(context.Background(), nil, nil, nil)
// 				return err
// 			},
// 		},
// 		{
// 			name:            "UpsertContact",
// 			route:           "/contact/",
// 			responsePayload: "{}",
// 			f: func(c *WalletClient) error {
// 				_, err := c.UpsertContact(context.Background(), "", "", nil)
// 				return err
// 			},
// 		},
// 		{
// 			name:            "UpsertContactForPaymail",
// 			route:           "/contact/",
// 			responsePayload: "{}",
// 			f: func(c *WalletClient) error {
// 				_, err := c.UpsertContactForPaymail(context.Background(), "", "", nil, "")
// 				return err
// 			},
// 		},
// 	}

// 	for _, tc := range tcs {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// given
// 			tmq := testTransportHandler{
// 				Type:      fixtures.RequestType,
// 				Path:      tc.route,
// 				Result:    tc.responsePayload,
// 				ClientURL: fixtures.ServerURL,
// 				Client:    WithHTTPClient,
// 			}

// 			client := getTestWalletClientWithOpts(tmq, WithXPriv(fixtures.XPrivString))

// 			// when
// 			err := tc.f(client)

// 			// then
// 			assert.NoError(t, err)
// 		})
// 	}

// }

// func TestConfirmContact(t *testing.T) {
// 	t.Run("TOTP is valid - call Confirm Action", func(t *testing.T) {
// 		// given
// 		tmq := testTransportHandler{
// 			Type:      fixtures.RequestType,
// 			Path:      "/contact/confirmed/",
// 			Result:    "{}",
// 			ClientURL: fixtures.ServerURL,
// 			Client:    WithHTTPClient,
// 		}

// 		clientMaker := func(opts ...ClientOps) (*WalletClient, error) {
// 			return getTestWalletClientWithOpts(tmq, opts...), nil
// 		}

// 		alice := makeMockUser("alice", clientMaker)
// 		bob := makeMockUser("bob", clientMaker)

// 		totp, err := alice.client.GenerateTotpForContact(bob.contact, 30, 2)
// 		require.NoError(t, err)

// 		// when
// 		result := bob.client.ConfirmContact(context.Background(), alice.contact, totp, bob.paymail, 30, 2)

// 		// then
// 		require.Nil(t, result)
// 	})

// 	t.Run("TOTP is invalid -  do not call Confirm Action", func(t *testing.T) {
// 		// given
// 		tmq := testTransportHandler{
// 			Type:      fixtures.RequestType,
// 			Path:      "/unknown/",
// 			Result:    "{}",
// 			ClientURL: fixtures.ServerURL,
// 			Client:    WithHTTPClient,
// 		}

// 		clientMaker := func(opts ...ClientOps) (*WalletClient, error) {
// 			return getTestWalletClientWithOpts(tmq, opts...), nil
// 		}

// 		alice := makeMockUser("alice", clientMaker)
// 		bob := makeMockUser("bob", clientMaker)

// 		totp, err := alice.client.GenerateTotpForContact(bob.contact, 30, 2)
// 		require.NoError(t, err)

// 		//make sure the wrongTotp is not the same as the generated one
// 		wrongTotp := incrementDigits(totp) //the length should remain the same

// 		// when
// 		result := bob.client.ConfirmContact(context.Background(), alice.contact, wrongTotp, bob.paymail, 30, 2)

// 		// then
// 		require.NotNil(t, result)
// 		require.Equal(t, result.Error(), "totp is invalid")
// 	})
// }

// // incrementDigits takes a string of digits and increments each digit by 1.
// // Digits wrap around such that '9' becomes '0'.
// func incrementDigits(input string) string {
// 	var result strings.Builder

// 	for _, c := range input {
// 		if c == '9' {
// 			result.WriteRune('0')
// 		} else {
// 			result.WriteRune(c + 1)
// 		}
// 	}

// 	return result.String()
// }
