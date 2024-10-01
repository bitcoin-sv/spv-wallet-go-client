package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	responsemodels "github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/stretchr/testify/require"
)

// TestContactActionsRouting will test routing
func TestContactActionsRouting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/contacts/" + fixtures.PaymailAddress:
			switch r.Method {
			case http.MethodPut:
				content := response.CreateContactResponse{
					Contact:        fixtures.Contact,
					AdditionalInfo: map[string]string{},
				}
				json.NewEncoder(w).Encode(content)
			case http.MethodDelete:
				json.NewEncoder(w).Encode(map[string]any{})
			case http.MethodGet:
				json.NewEncoder(w).Encode(fixtures.Contact)
			}
		case "/api/v1/contacts/" + fixtures.PaymailAddress + "/confirmation":
			switch r.Method {
			case http.MethodPost, http.MethodDelete:
				json.NewEncoder(w).Encode(map[string]string{"result": string(responsemodels.ContactNotConfirmed)})
			}
		case "/api/v1/contacts/":
			if r.Method == http.MethodGet {
				content := response.PageModel[response.Contact]{
					Content: []*response.Contact{fixtures.Contact},
				}
				json.NewEncoder(w).Encode(content)
			}
		case "/api/v1/invitations/" + fixtures.PaymailAddress + "/contacts":
			if r.Method == http.MethodPost {
				json.NewEncoder(w).Encode(map[string]any{})
			}
		case "/api/v1/invitations/" + fixtures.PaymailAddress:
			if r.Method == http.MethodDelete {
				json.NewEncoder(w).Encode(map[string]any{})
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
