package auth_test

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
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
	authenticator, err := auth.NewAccessKeyAuthenticator("")

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, errors.ErrEmptyAccessKey)
}

func TestAccessKeyAuthenticator_Authenticate(t *testing.T) {
	// given:
	authenticator, err := auth.NewAccessKeyAuthenticator(spvwallettest.UserPrivAccessKey)
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
	authenticator, err := auth.NewXprivAuthenticator("")

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, errors.ErrEmptyXprivKey)
}

func TestXprivAuthenitcator_Authenticate(t *testing.T) {
	// given:
	authenticator, err := auth.NewXprivAuthenticator(spvwallettest.UserXPriv)
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
	authenticator, err := auth.NewXpubOnlyAuthenticator("")

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, errors.ErrEmptyPubKey)
}

func TestXpubOnlyAuthenticator_Authenticate(t *testing.T) {
	// given:
	authenticator, err := auth.NewXpubOnlyAuthenticator(spvwallettest.UserXPriv)
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
	require.Equal(t, []string{spvwallettest.UserPubAccessKey}, h[xAuthKey])
}

func requireXpubHeaderToBeSet(t *testing.T, h http.Header) {
	require.Equal(t, []string{spvwallettest.UserXPub}, h[xAuthXPubKey])
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
