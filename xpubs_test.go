package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

type xpub struct {
	CurrentBalance uint64           `json:"current_balance"`
	Metadata       *models.Metadata `json:"metadata"`
}

func TestXpub(t *testing.T) {
	var update bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var response xpub
		// Check path and method to customize the response
		switch {
		case r.URL.Path == "/api/v1/users/current":
			metadata := &models.Metadata{"key": "value"}
			if update {
				metadata = &models.Metadata{"updated": "info"}
			}
			response = xpub{
				CurrentBalance: 1234,
				Metadata:       metadata,
			}
		}
		respBytes, _ := json.Marshal(response)
		w.Write(respBytes)
	}))
	defer server.Close()
	keys, err := xpriv.Generate()
	require.NoError(t, err)

	client := NewWithXPriv(server.URL, keys.XPriv())
	require.NotNil(t, client.xPriv)

	t.Run("GetXPub", func(t *testing.T) {
		xpub, err := client.GetXPub(context.Background())
		require.NoError(t, err)
		require.NotNil(t, xpub)
		require.Equal(t, uint64(1234), xpub.CurrentBalance)
		require.Equal(t, "value", xpub.Metadata["key"])
	})

	t.Run("UpdateXPubMetadata", func(t *testing.T) {
		update = true
		metadata := map[string]any{"updated": "info"}
		xpub, err := client.UpdateXPubMetadata(context.Background(), metadata)
		require.NoError(t, err)
		require.NotNil(t, xpub)
		require.Equal(t, "info", xpub.Metadata["updated"])
	})
}
