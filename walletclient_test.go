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
		client, err := NewWithXPriv(keys.XPriv(), server.URL)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, keys.XPriv(), *client.xPrivString)
		require.NotNil(t, client.httpClient)
		require.True(t, *client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithXPrivate fail", func(t *testing.T) {
		xPriv := "invalid_key"
		client, err := NewWithXPriv(xPriv, "http://example.com")
		require.Error(t, err) // Expect error due to invalid key
		require.Nil(t, client)
	})

	t.Run("NewWalletClientWithXPublic success", func(t *testing.T) {
		keys, err := xpriv.Generate()
		require.NoError(t, err)
		client, err := NewWithXPub(keys.XPub().String(), server.URL)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, keys.XPub().String(), *client.xPubString)
		require.NotNil(t, client.httpClient)
		require.False(t, *client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithXPublic fail", func(t *testing.T) {
		xpub := "invalid_key"
		client, err := NewWithXPub(xpub, server.URL)
		require.Error(t, err) // Expect error due to invalid key
		require.Nil(t, client)
	})

	t.Run("NewWalletClientWithAdminKey success", func(t *testing.T) {
		client, err := NewWithAdminKey(fixtures.XPrivString, server.URL)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, fixtures.XPrivString, *client.xPrivString)
		require.Equal(t, fixtures.XPrivString, client.adminXPriv.String())
		require.Equal(t, server.URL, *client.server)
		require.NotNil(t, client.httpClient)
		require.True(t, *client.signRequest)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithAdminKey fail", func(t *testing.T) {
		adminKey := "invalid_key"
		client, err := NewWithAdminKey(adminKey, server.URL)
		require.Error(t, err)
		require.Nil(t, client)
	})

	t.Run("NewWalletClientWithAccessKey success", func(t *testing.T) {
		// Attempt to create a new WalletClient with an access key
		accessKey := fixtures.AccessKeyString
		client, err := NewWithAccessKey(accessKey, server.URL)

		require.NoError(t, err)
		require.NotNil(t, client)

		require.Equal(t, &accessKey, client.accessKeyString)
		require.Equal(t, &server.URL, client.server)
		require.True(t, *client.signRequest)
		require.NotNil(t, client.httpClient)

		// Ensure HTTP calls can be made
		resp, err := client.httpClient.Get(server.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NewWalletClientWithAccessKey fail", func(t *testing.T) {
		accessKey := "invalid_key"
		client, err := NewWithAccessKey(accessKey, server.URL)

		require.Error(t, err)
		require.Nil(t, client)
	})
}
