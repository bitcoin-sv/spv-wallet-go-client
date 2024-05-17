// Package walletclient here we are testing walletclient public methods
package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

// TestAccessKeys will test the AccessKey methods
func TestAccessKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/access-key":
			if r.Method == http.MethodGet {
				json.NewEncoder(w).Encode(fixtures.AccessKey)
			} else if r.Method == http.MethodPost {
				json.NewEncoder(w).Encode(fixtures.AccessKey)
			} else if r.Method == http.MethodDelete {
				json.NewEncoder(w).Encode(fixtures.AccessKey)
			}
		case "/v1/access-key/search":
			json.NewEncoder(w).Encode([]*models.AccessKey{fixtures.AccessKey})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewWithAccessKey(server.URL, fixtures.AccessKeyString)
	require.NotNil(t, client.accessKey)

	t.Run("GetAccessKey", func(t *testing.T) {
		accessKey, err := client.GetAccessKey(context.Background(), fixtures.AccessKey.ID)
		require.NoError(t, err)
		require.Equal(t, fixtures.AccessKey, accessKey)
	})

	t.Run("GetAccessKeys", func(t *testing.T) {
		accessKeys, err := client.GetAccessKeys(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, []*models.AccessKey{fixtures.AccessKey}, accessKeys)
	})

	t.Run("CreateAccessKey", func(t *testing.T) {
		accessKey, err := client.CreateAccessKey(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, fixtures.AccessKey, accessKey)
	})

	t.Run("RevokeAccessKey", func(t *testing.T) {
		accessKey, err := client.RevokeAccessKey(context.Background(), fixtures.AccessKey.ID)
		require.NoError(t, err)
		require.Equal(t, fixtures.AccessKey, accessKey)
	})
}
