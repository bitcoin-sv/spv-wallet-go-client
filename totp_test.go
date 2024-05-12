package walletclient

// import (
// 	"encoding/hex"
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

// 	"github.com/bitcoin-sv/spv-wallet/models"
// 	"github.com/libsv/go-bk/bip32"
// 	"github.com/stretchr/testify/require"

// 	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
// 	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
// )

// func TestGenerateTotpForContact(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		// given
// 		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
// 		require.NoError(t, err)

// 		contact := models.Contact{PubKey: fixtures.PubKey}

// 		// when
// 		pass, err := sut.GenerateTotpForContact(&contact, 30, 2)

// 		// then
// 		require.NoError(t, err)
// 		require.Len(t, pass, 2)
// 	})

// 	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
// 		// given
// 		sut, err := New(WithXPub(fixtures.XPubString), WithHTTP("localhost:3001"))
// 		require.NoError(t, err)

// 		// when
// 		_, err = sut.GenerateTotpForContact(nil, 30, 2)

// 		// then
// 		require.ErrorIs(t, err, ErrClientInitNoXpriv)
// 	})

// 	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
// 		// given
// 		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
// 		require.NoError(t, err)

// 		contact := models.Contact{PubKey: "invalid-pk-format"}

// 		// when
// 		_, err = sut.GenerateTotpForContact(&contact, 30, 2)

// 		// then
// 		require.ErrorContains(t, err, "contact's PubKey is invalid:")

// 	})
// }

// func TestValidateTotpForContact(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		// given
// 		clientMaker := func(opts ...ClientOps) (*WalletClient, error) {
// 			allOptions := append(opts, WithHTTP("localhost:3001"))
// 			return New(allOptions...)
// 		}
// 		alice := makeMockUser("alice", clientMaker)
// 		bob := makeMockUser("bob", clientMaker)

// 		pass, err := alice.client.GenerateTotpForContact(bob.contact, 3600, 2)
// 		require.NoError(t, err)

// 		// when
// 		result, err := bob.client.ValidateTotpForContact(alice.contact, pass, bob.paymail, 3600, 2)

// 		// then
// 		require.NoError(t, err)
// 		require.True(t, result)
// 	})

// 	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
// 		// given
// 		sut, err := New(WithXPub(fixtures.XPubString), WithHTTP("localhost:3001"))
// 		require.NoError(t, err)

// 		// when
// 		_, err = sut.ValidateTotpForContact(nil, "", fixtures.PaymailAddress, 30, 2)

// 		// then
// 		require.ErrorIs(t, err, ErrClientInitNoXpriv)
// 	})

// 	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
// 		// given
// 		sut, err := New(WithXPriv(fixtures.XPrivString), WithHTTP("localhost:3001"))
// 		require.NoError(t, err)

// 		contact := models.Contact{PubKey: "invalid-pk-format"}

// 		// when
// 		_, err = sut.ValidateTotpForContact(&contact, "", fixtures.PaymailAddress, 30, 2)

// 		// then
// 		require.ErrorContains(t, err, "contact's PubKey is invalid:")

// 	})
// }

// type mockUser struct {
// 	contact *models.Contact
// 	client  *WalletClient
// 	paymail string
// }

// func makeMockUser(name string, clientMaker func(opts ...ClientOps) (*WalletClient, error)) mockUser {
// 	keys, _ := xpriv.Generate()
// 	paymail := name + "@example.com"
// 	client, _ := clientMaker(WithXPriv(keys.XPriv()))
// 	pki := makeMockPKI(keys.XPub().String())
// 	contact := models.Contact{PubKey: pki, Paymail: paymail}
// 	return mockUser{
// 		contact: &contact,
// 		client:  client,
// 		paymail: paymail,
// 	}
// }

// func makeMockPKI(xpub string) string {
// 	xPub, _ := bip32.NewKeyFromString(xpub)
// 	magicNumberOfInheritance := 3 //2+1; 2: because of the way spv-wallet stores xpubs in db; 1: to make a PKI
// 	var err error
// 	for i := 0; i < magicNumberOfInheritance; i++ {
// 		xPub, err = xPub.Child(0)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	pubKey, err := xPub.ECPubKey()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return hex.EncodeToString(pubKey.SerialiseCompressed())
// }
