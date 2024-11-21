package totp_test

import (
	"testing"
	"time"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestClient_GenerateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		contact := models.Contact{PubKey: clienttest.PubKey}
		wc := totp.New(clienttest.ExtendedKey(t))

		// when
		pass, err := wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.NoError(t, err)
		require.Len(t, pass, 2)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		contact := models.Contact{PubKey: "invalid-pk-format"}
		wc := totp.New(clienttest.ExtendedKey(t))

		// when
		_, err := wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.ErrorIs(t, err, errors.ErrContactPubKeyInvalid)
	})
}

func TestClient_ValidateTotpForContact(t *testing.T) {
	cfg := client.Config{
		Addr:    clienttest.TestAPIAddr,
		Timeout: 5 * time.Second,
	}
	t.Run("success", func(t *testing.T) {
		// given
		clientAlice, err := client.NewWithXPriv(cfg, clienttest.AliceXPriv)
		require.NoError(t, err)

		clientBob, err := client.NewWithXPriv(cfg, clienttest.BobXPriv)
		require.NoError(t, err)

		// and
		aliceContact := &models.Contact{
			PubKey:  clienttest.MockPKI(t, clienttest.AliceXPub),
			Paymail: "alice@example.com",
		}

		bobContact := &models.Contact{
			PubKey:  clienttest.MockPKI(t, clienttest.BobXPub),
			Paymail: "bob@example.com",
		}

		// when
		passcode, err := clientAlice.GenerateTotpForContact(bobContact, 3600, 6)

		// then
		require.NoError(t, err)

		// when
		err = clientBob.ValidateTotpForContact(aliceContact, passcode, bobContact.Paymail, 3600, 6)

		// then
		require.NoError(t, err)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut, err := client.NewWithXPriv(cfg, clienttest.UserXPriv)
		require.NoError(t, err)

		// and
		invalidContact := &models.Contact{
			PubKey:  "invalid_pub_key_format",
			Paymail: "invalid@example.com",
		}

		// when
		err = sut.ValidateTotpForContact(invalidContact, "123456", "someone@example.com", 3600, 6)

		// when
		require.Contains(t, err.Error(), "contact's PubKey is invalid")
	})
}
