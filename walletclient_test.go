package walletclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
)

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func mustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

type testTransportHandler struct {
	ClientURL string
	Client    func(serverURL string, httpClient *http.Client) ClientOps
	Path      string
	Queries   []*testTransportHandlerRequest
	Result    string
	Type      string
}

type testTransportHandlerRequest struct {
	Path   string
	Result func(w http.ResponseWriter, req *http.Request)
}

// TestNewWalletClient will test the TestNewWalletClient method
func TestNewWalletClient(t *testing.T) {

	t.Run("empty xpriv", func(t *testing.T) {
		client, err := NewWalletClientWithXPrivate("", fixtures.ServerURL, false)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("invalid xpriv", func(t *testing.T) {
		client, err := NewWalletClientWithXPrivate("invalid-xpriv", fixtures.ServerURL, false)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("valid client", func(t *testing.T) {
		client, err := NewWalletClientWithXPrivate(fixtures.XPrivString, fixtures.ServerURL, false)

		require.NoError(t, err)
		assert.IsType(t, WalletClient{}, *client)
	})

	t.Run("valid xPub client", func(t *testing.T) {
		client, err := NewWalletClientWithXPublic(fixtures.XPubString, fixtures.ServerURL, false)
		require.NoError(t, err)
		assert.IsType(t, WalletClient{}, *client)
	})

	t.Run("invalid xPub client", func(t *testing.T) {
		client, err := NewWalletClientWithXPublic("invalid-xpub", fixtures.ServerURL, false)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("valid access keys", func(t *testing.T) {
		client, err := NewWalletClientWithAccessKey(fixtures.AccessKeyString, fixtures.ServerURL, false)
		require.NoError(t, err)
		assert.IsType(t, WalletClient{}, *client)
	})

	t.Run("invalid access keys", func(t *testing.T) {
		client, err := NewWalletClientWithAccessKey("invalid-access-key", fixtures.ServerURL, false)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("valid access key WIF", func(t *testing.T) {
		wifKey, _ := bitcoin.PrivateKeyToWif(fixtures.AccessKeyString)
		client, err := NewWalletClientWithAccessKey(wifKey.String(), fixtures.ServerURL, false)
		require.NoError(t, err)
		assert.IsType(t, WalletClient{}, *client)
	})
}

// TestSetSignRequest will test the sign request setter
func TestSetSignRequest(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		client, _ := NewWalletClientWithXPrivate(fixtures.XPrivString, fixtures.ServerURL, true)
		assert.True(t, client.IsSignRequest())
	})

	t.Run("false", func(t *testing.T) {
		client, _ := NewWalletClientWithXPrivate(fixtures.XPrivString, fixtures.ServerURL, false)
		assert.False(t, client.IsSignRequest())
	})
}
