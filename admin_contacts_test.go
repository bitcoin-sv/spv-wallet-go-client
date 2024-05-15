package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
)

// TestAdminContactActions testing Admin contacts methods
func TestAdminContactActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/admin/contact/search" && r.Method == http.MethodPost:
			c := fixtures.Contact
			c.ID = "1"
			contacts := []*models.Contact{c}
			json.NewEncoder(w).Encode(contacts)
		case r.URL.Path == "/admin/contact/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			json.NewEncoder(w).Encode(contact)
		case r.URL.Path == "/admin/contact/1" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/admin/contact/accepted/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			contact.Status = "accepted"
			json.NewEncoder(w).Encode(contact)
		case r.URL.Path == "/admin/contact/rejected/1" && r.Method == http.MethodPatch:
			contact := fixtures.Contact
			contact.Status = "rejected"
			json.NewEncoder(w).Encode(contact)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewWithAdminKey(fixtures.XPrivString, server.URL)
	require.NoError(t, err)

	t.Run("AdminGetContacts", func(t *testing.T) {
		contacts, err := client.AdminGetContacts(context.Background(), nil, nil, nil)
		require.NoError(t, err)
		require.Equal(t, "1", contacts[0].ID)
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
		require.Equal(t, models.ContactStatus("accepted"), contact.Status)
	})

	t.Run("AdminRejectContact", func(t *testing.T) {
		contact, err := client.AdminRejectContact(context.Background(), "1")
		require.NoError(t, err)
		require.Equal(t, models.ContactStatus("rejected"), contact.Status)
	})
}
