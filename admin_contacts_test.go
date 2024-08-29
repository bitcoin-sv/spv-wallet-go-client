package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	responsemodels "github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/stretchr/testify/require"
)

// TestAdminContactActions testing Admin contacts methods
func TestAdminContactActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/admin/contact/search" && r.Method == http.MethodPost:
			c := fixtures.OldContact
			c.ID = "1"
			content := models.PagedResponse[*models.Contact]{
				Content: []*models.Contact{c},
			}
			json.NewEncoder(w).Encode(content)
		case r.URL.Path == "/v1/admin/contact/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			json.NewEncoder(w).Encode(contact)
		case r.URL.Path == "/v1/admin/contact/1" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/v1/admin/contact/accepted/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			contact.Status = responsemodels.ContactNotConfirmed
			json.NewEncoder(w).Encode(contact)
		case r.URL.Path == "/v1/admin/contact/rejected/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			contact.Status = responsemodels.ContactRejected
			json.NewEncoder(w).Encode(contact)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewWithAdminKey(server.URL, fixtures.XPrivString)
	require.NotNil(t, client.adminXPriv)

	t.Run("AdminGetContacts", func(t *testing.T) {
		contacts, err := client.AdminGetContacts(context.Background(), nil, nil, nil)
		require.NoError(t, err)
		require.Equal(t, "1", contacts.Content[0].ID)
	})

	t.Run("AdminUpdateContact", func(t *testing.T) {
		contact, err := client.AdminUpdateContact(context.Background(), "1", "Jane Doe", nil)
		require.NoError(t, err)
		require.Equal(t, "Test User", contact.FullName)
	})

	t.Run("AdminDeleteContact", func(t *testing.T) {
		err := client.AdminDeleteContact(context.Background(), "1")
		require.NoError(t, err)
	})

	t.Run("AdminAcceptContact", func(t *testing.T) {
		contact, err := client.AdminAcceptContact(context.Background(), "1")
		require.NoError(t, err)
		require.Equal(t, responsemodels.ContactNotConfirmed, contact.Status)
	})

	t.Run("AdminRejectContact", func(t *testing.T) {
		contact, err := client.AdminRejectContact(context.Background(), "1")
		require.NoError(t, err)
		require.Equal(t, responsemodels.ContactRejected, contact.Status)
	})
}
