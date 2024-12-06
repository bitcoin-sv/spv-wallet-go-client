package totp_test

import (
	"testing"
	"time"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestClient_GenerateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		contact := models.Contact{PubKey: spvwallettest.PubKey}
		wc, err := totp.New(spvwallettest.UserXPriv)
		require.NoError(t, err)

		// when
		pass, err := wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.NoError(t, err)
		require.Len(t, pass, 2)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		contact := models.Contact{PubKey: "invalid-pk-format"}
		wc, err := totp.New(spvwallettest.UserXPriv)
		require.NoError(t, err)

		// when
		_, err = wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.ErrorIs(t, err, errors.ErrContactPubKeyInvalid)
	})
}

func TestClient_ValidateTotpForContact(t *testing.T) {
	cfg := config.Config{
		Addr:    spvwallettest.TestAPIAddr,
		Timeout: 5 * time.Second,
	}
	t.Run("success", func(t *testing.T) {
		// given
		clientAlice, err := client.NewUserAPIWithXPriv(cfg, spvwallettest.AliceXPriv)
		require.NoError(t, err)

		clientBob, err := client.NewUserAPIWithXPriv(cfg, spvwallettest.BobXPriv)
		require.NoError(t, err)

		// and
		aliceContact := &models.Contact{
			PubKey:  spvwallettest.MockPKI(t, spvwallettest.AliceXPub),
			Paymail: "alice@example.com",
		}

		bobContact := &models.Contact{
			PubKey:  spvwallettest.MockPKI(t, spvwallettest.BobXPub),
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
		sut, err := client.NewUserAPIWithXPriv(cfg, spvwallettest.UserXPriv)
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

	t.Run("xpriv empty", func(t *testing.T) {
		_, err := client.NewUserAPIWithXPriv(cfg, "")
		require.Error(t, err)
		require.ErrorIs(t, err, errors.ErrEmptyXprivKey)
	})
}
