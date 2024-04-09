package walletclient

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestGenerateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
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
		sut, err := New(WithXPub(fixtures.XPubString), WithHTTP("localhost:3001"))
		require.NoError(t, err)

		// when
		_, err = sut.GenerateTotpForContact(nil, 30, 2)

		// then
		require.ErrorIs(t, err, ErrClientInitNoXpriv)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
		require.NoError(t, err)

		contact := models.Contact{PubKey: "invalid-pk-format"}

		// when
		_, err = sut.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.ErrorContains(t, err, "contact's PubKey is invalid:")

	})
}

func TestValidateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
		require.NoError(t, err)

		contact := models.Contact{PubKey: fixtures.PubKey}
		pass, err := sut.GenerateTotpForContact(&contact, 30, 2)
		require.NoError(t, err)

		// when
		result, err := sut.ValidateTotpForContact(&contact, pass, 30, 2)

		// then
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		// given
		sut, err := New(WithXPub(fixtures.XPubString), WithHTTP("localhost:3001"))
		require.NoError(t, err)

		// when
		_, err = sut.ValidateTotpForContact(nil, "", 30, 2)

		// then
		require.ErrorIs(t, err, ErrClientInitNoXpriv)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
		require.NoError(t, err)

		contact := models.Contact{PubKey: "invalid-pk-format"}

		// when
		_, err = sut.ValidateTotpForContact(&contact, "", 30, 2)

		// then
		require.ErrorContains(t, err, "contact's PubKey is invalid:")

	})
}
