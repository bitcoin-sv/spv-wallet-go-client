package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

// TestContactActionsRouting will test routing
func TestContactActionsRouting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/contact/rejected/"):
			if r.Method == http.MethodPatch {
				json.NewEncoder(w).Encode(map[string]string{"result": "rejected"})
			}
		case r.URL.Path == "/v1/contact/accepted/":
			if r.Method == http.MethodPost {
				json.NewEncoder(w).Encode(map[string]string{"result": "accepted"})
			}
		case r.URL.Path == "/v1/contact/search":
			if r.Method == http.MethodPost {
				content := models.PagedResponse[*models.Contact]{
					Content: []*models.Contact{fixtures.Contact},
				}
				json.NewEncoder(w).Encode(content)
			}
		case strings.HasPrefix(r.URL.Path, "/v1/contact/"):
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				json.NewEncoder(w).Encode(map[string]string{"result": "upserted"})
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewWithAccessKey(server.URL, fixtures.AccessKeyString)
	require.NotNil(t, client.accessKey)

	t.Run("RejectContact", func(t *testing.T) {
		err := client.RejectContact(context.Background(), fixtures.PaymailAddress)
		require.NoError(t, err)
	})

	t.Run("AcceptContact", func(t *testing.T) {
		err := client.AcceptContact(context.Background(), fixtures.PaymailAddress)
		require.NoError(t, err)
	})

	t.Run("GetContacts", func(t *testing.T) {
		contacts, err := client.GetContacts(context.Background(), nil, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, contacts)
	})

	t.Run("UpsertContact", func(t *testing.T) {
		contact, err := client.UpsertContact(context.Background(), "test-id", "test@paymail.com", "", nil)
		require.NoError(t, err)
		require.NotNil(t, contact)
	})

	t.Run("UpsertContactForPaymail", func(t *testing.T) {
		contact, err := client.UpsertContactForPaymail(context.Background(), "test-id", "test@paymail.com", nil, "test@paymail.com")
		require.NoError(t, err)
		require.NotNil(t, contact)
	})
}
