package walletclient

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bip32"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func TestGenerateTotpForContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		sut := NewWithXPriv("localhost:3001", fixtures.XPrivString)
		require.NotNil(t, sut.xPriv)

		contact := models.Contact{PubKey: fixtures.PubKey}
		// when
		pass, err := sut.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.NoError(t, err)
		require.Len(t, pass, 2)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		// given
		sut := NewWithXPub("localhost:3001", fixtures.XPubString)
		require.NotNil(t, sut.xPub)
		// when
		_, err := sut.GenerateTotpForContact(nil, 30, 2)

		// then
		require.ErrorIs(t, err, ErrClientInitNoXpriv)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut := NewWithXPriv("localhost:3001", fixtures.XPrivString)
		require.NotNil(t, sut.xPriv)

		contact := models.Contact{PubKey: "invalid-pk-format"}
		// when
		_, err := sut.GenerateTotpForContact(&contact, 30, 2)

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
		clientAlice := NewWithXPriv(server.URL, aliceKeys.XPriv())
		require.NotNil(t, clientAlice.xPriv)
		clientBob := NewWithXPriv(server.URL, bobKeys.XPriv())
		require.NotNil(t, clientBob.xPriv)

		aliceContact := &models.Contact{
			PubKey:  makeMockPKI(aliceKeys.XPub().String()),
			Paymail: "bob@example.com",
		}

		bobContact := &models.Contact{
			PubKey:  makeMockPKI(bobKeys.XPub().String()),
			Paymail: "bob@example.com",
		}

		// Generate and validate TOTP
		passcode, err := clientAlice.GenerateTotpForContact(bobContact, 3600, 6)
		require.NoError(t, err)
		result, err := clientBob.ValidateTotpForContact(aliceContact, passcode, bobContact.Paymail, 3600, 6)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		client := NewWithXPub(server.URL, "invalid_xpub")
		require.Nil(t, client.xPub)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		sut := NewWithXPriv(server.URL, fixtures.XPrivString)

		invalidContact := &models.Contact{
			PubKey:  "invalid_pub_key_format",
			Paymail: "invalid@example.com",
		}

		_, err := sut.ValidateTotpForContact(invalidContact, "123456", "someone@example.com", 3600, 6)
		require.Error(t, err)
		require.Contains(t, err.Error(), "contact's PubKey is invalid")
	})
}

func makeMockPKI(xpub string) string {
	xPub, _ := bip32.NewKeyFromString(xpub)
	var err error
	for i := 0; i < 3; i++ { //magicNumberOfInheritance is 3 -> 2+1; 2: because of the way spv-wallet stores xpubs in db; 1: to make a PKI
		xPub, err = xPub.Child(0)
		if err != nil {
			panic(err)
		}
	}

	pubKey, err := xPub.ECPubKey()
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(pubKey.SerialiseCompressed())
}
