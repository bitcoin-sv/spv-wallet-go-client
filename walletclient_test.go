package walletclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestNewWalletClientWithXPrivate(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	// Test creating a client with a valid xPriv
	xPriv := "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	client, err := NewWalletClientWithXPrivate(xPriv, server.URL, true)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, &xPriv, client.xPrivString)
	assert.NotNil(t, client.httpClient)
	assert.True(t, *client.signRequest)

	// Ensure HTTP calls can be made
	resp, err := client.httpClient.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestKeyInitialization(t *testing.T) {
	xPriv := "invalid_key"
	client, err := NewWalletClientWithXPrivate(xPriv, "http://example.com", true)
	assert.Error(t, err) // Expect error due to invalid key
	assert.Nil(t, client)
}
