package walletclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func TestNewWalletClient(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	// Test creating a client with a valid xPriv
	t.Run("NewWalletClientWithXPrivate success", func(t *testing.T) {
		keys, err := xpriv.Generate()
		require.NoError(t, err)
		client := NewWithXPriv(server.URL, keys.XPriv())
		require.NotNil(t, client.xPriv)
		require.Equal(t, keys.XPriv(), client.xPriv.String())
		require.NotNil(t, client.httpClient)
		require.True(t, client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithXPrivate fail", func(t *testing.T) {
		xPriv := "invalid_key"
		client := NewWithXPriv(xPriv, "http://example.com")
		require.Nil(t, client.xPriv)
	})

	t.Run("NewWalletClientWithXPublic success", func(t *testing.T) {
		keys, err := xpriv.Generate()
		require.NoError(t, err)
		client := NewWithXPub(server.URL, keys.XPub().String())
		require.NotNil(t, client.xPub)
		require.Equal(t, keys.XPub().String(), client.xPub.String())
		require.NotNil(t, client.httpClient)
		require.False(t, client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithXPublic fail", func(t *testing.T) {
		client := NewWithXPub(server.URL, "invalid_key")
		require.Nil(t, client.xPub)
	})

	t.Run("NewWalletClientWithAdminKey success", func(t *testing.T) {
		client := NewWithAdminKey(server.URL, fixtures.XPrivString)
		require.NotNil(t, client.adminXPriv)
		require.Nil(t, client.xPriv)
		require.Equal(t, fixtures.XPrivString, client.adminXPriv.String())
		require.Equal(t, server.URL, client.server)
		require.NotNil(t, client.httpClient)
		require.True(t, client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithAdminKey fail", func(t *testing.T) {
		client := NewWithAdminKey(server.URL, "invalid_key")
		require.Nil(t, client.adminXPriv)
	})

	t.Run("NewWalletClientWithAccessKey success", func(t *testing.T) {
		// Attempt to create a new WalletClient with an access key
		client := NewWithAccessKey(server.URL, fixtures.AccessKeyString)
		require.NotNil(t, client.accessKey)

		require.Equal(t, server.URL, client.server)
		require.True(t, client.signRequest)
		require.NotNil(t, client.httpClient)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithAccessKey fail", func(t *testing.T) {
		client := NewWithAccessKey(server.URL, "invalid_key")
		require.Nil(t, client.accessKey)
	})
}
