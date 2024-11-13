package auth_test

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

const (
	xAuthKey          = "X-Auth-Key"
	xAuthXPubKey      = "X-Auth-Xpub"
	xAuthHashKey      = "X-Auth-Hash"
	xAuthNonceKey     = "X-Auth-Nonce"
	xAuthTimeKey      = "X-Auth-Time"
	xAuthSignatureKey = "X-Auth-Signature"
)

func TestAccessKeyAuthenitcator_NewWithNilAccessKey(t *testing.T) {
	// when:
	authenticator, err := auth.NewAccessKeyAuthenticator(nil)

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrEcPrivateKey)
}

func TestAccessKeyAuthenticator_Authenticate(t *testing.T) {
	// given:
	key := clienttest.PrivateKey(t)
	authenticator, err := auth.NewAccessKeyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXAuthHeaderToBeSet(t, req.Header)
	requireSignatureHeadersToBeSet(t, req.Header)
}

func TestXprivAuthenitcator_NewWithNilXpriv(t *testing.T) {
	// when:
	authenticator, err := auth.NewXprivAuthenticator(nil)

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrBip32ExtendedKey)
}

func TestXprivAuthenitcator_Authenticate(t *testing.T) {
	// given:
	key := clienttest.ExtendedKey(t)
	authenticator, err := auth.NewXprivAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXpubHeaderToBeSet(t, req.Header)
	requireSignatureHeadersToBeSet(t, req.Header)
}

func TestXpubOnlyAuthenticator_NewWithNilXpub(t *testing.T) {
	// when:
	authenticator, err := auth.NewXpubOnlyAuthenticator(nil)

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrBip32ExtendedKey)
}

func TestXpubOnlyAuthenticator_Authenticate(t *testing.T) {
	// given:
	key := clienttest.ExtendedKey(t)

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXpubHeaderToBeSet(t, req.Header)
}

func requireXAuthHeaderToBeSet(t *testing.T, h http.Header) {
	require.Equal(t, []string{clienttest.UserPubAccessKey}, h[xAuthKey])
}

func requireXpubHeaderToBeSet(t *testing.T, h http.Header) {
	require.Equal(t, []string{clienttest.UserXPub}, h[xAuthXPubKey])
}

func requireSignatureHeadersToBeSet(t *testing.T, h http.Header) {
	expected := []string{
		xAuthHashKey,
		xAuthNonceKey,
		xAuthTimeKey,
		xAuthSignatureKey,
	}

	actual := make([]string, 0, len(expected))
	for k := range h {
		actual = append(actual, k)
	}
	require.Subset(t, actual, expected)
}
