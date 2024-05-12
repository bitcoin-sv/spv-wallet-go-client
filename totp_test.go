package walletclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func TestGenerateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		sut, err := NewWalletClientWithXPrivate(fixtures.XPrivString, "localhost:3001", false)
		require.NoError(t, err)

		contact := models.Contact{PubKey: fixtures.PubKey}

		// when
		pass, err := sut.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.NoError(t, err)
		require.Len(t, pass, 2)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		// given
		sut, err := NewWalletClientWithXPublic(fixtures.XPubString, "localhost:3001", false)
		require.NoError(t, err)

		// when
		_, err = sut.GenerateTotpForContact(nil, 30, 2)

		// then
		require.ErrorIs(t, err, ErrClientInitNoXpriv)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut, err := NewWalletClientWithXPrivate(fixtures.XPrivString, "localhost:3001", false)
		require.NoError(t, err)

		contact := models.Contact{PubKey: "invalid-pk-format"}

		// when
		_, err = sut.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.ErrorContains(t, err, "contact's PubKey is invalid:")

	})
}

func TestValidateTotpForContact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler could be adjusted depending on the expected API endpoints
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("123456")) // Simulate a TOTP response for any requests
	}))
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		aliceKeys, err := xpriv.Generate()
		require.NoError(t, err)
		bobKeys, err := xpriv.Generate()
		require.NoError(t, err)

		// Set up the WalletClient for Alice and Bob
		clientAlice, err := NewWalletClientWithXPrivate(aliceKeys.XPriv(), server.URL, true)
		require.NoError(t, err)
		clientBob, err := NewWalletClientWithXPrivate(bobKeys.XPriv(), server.URL, true)
		require.NoError(t, err)

		require.NoError(t, err)
		bobContact := &models.Contact{
			PubKey:  bobKeys.XPub().String(),
			Paymail: "bob@example.com",
		}

		// Generate and validate TOTP
		passcode, err := clientAlice.GenerateTotpForContact(bobContact, 3600, 6)
		require.NoError(t, err)
		result, err := clientBob.ValidateTotpForContact(bobContact, passcode, "alice@example.com", 3600, 6)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		client, err := NewWalletClientWithXPublic("invalid_xpub", server.URL, true)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		xPrivString := "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
		sut, err := NewWalletClientWithXPrivate(xPrivString, server.URL, true)
		require.NoError(t, err)

		invalidContact := &models.Contact{
			PubKey:  "invalid_pub_key_format",
			Paymail: "invalid@example.com",
		}

		_, err = sut.ValidateTotpForContact(invalidContact, "123456", "someone@example.com", 3600, 6)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contact's PubKey is invalid")
	})
}
