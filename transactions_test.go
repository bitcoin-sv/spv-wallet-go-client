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

func TestTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/transaction":
			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(fixtures.Transaction)
			case http.MethodPost:
				json.NewEncoder(w).Encode(fixtures.Transaction)
			case http.MethodPatch:
				var input map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
					return
				}
				response := fixtures.Transaction
				response.Metadata = input["metadata"].(map[string]interface{})
				response.ID = input["id"].(string)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		case "/transaction/search":
			json.NewEncoder(w).Encode([]*models.Transaction{fixtures.Transaction})
		case "/transaction/count":
			json.NewEncoder(w).Encode(1)
		case "/transaction/record":
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				response := fixtures.Transaction
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewWithXPriv(fixtures.XPrivString, server.URL)
	require.NoError(t, err)

	t.Run("GetTransaction", func(t *testing.T) {
		tx, err := client.GetTransaction(context.Background(), fixtures.Transaction.ID)
		require.NoError(t, err)
		require.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("GetTransactions", func(t *testing.T) {
		conditions := map[string]interface{}{
			"fee":         map[string]interface{}{"$lt": 100},
			"total_value": map[string]interface{}{"$lt": 740},
		}
		txs, err := client.GetTransactions(context.Background(), conditions, fixtures.TestMetadata, nil)
		require.NoError(t, err)
		require.Equal(t, []*models.Transaction{fixtures.Transaction}, txs)
	})

	t.Run("GetTransactionsCount", func(t *testing.T) {
		count, err := client.GetTransactionsCount(context.Background(), nil, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
	})

	t.Run("RecordTransaction", func(t *testing.T) {
		tx, err := client.RecordTransaction(context.Background(), fixtures.Transaction.Hex, "", fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("UpdateTransactionMetadata", func(t *testing.T) {
		tx, err := client.UpdateTransactionMetadata(context.Background(), fixtures.Transaction.ID, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("SendToRecipients", func(t *testing.T) {
		recipients := []*Recipients{{
			OpReturn: fixtures.DraftTx.Configuration.Outputs[0].OpReturn,
			Satoshis: 1000,
			To:       fixtures.Destination.Address,
		}}
		tx, err := client.SendToRecipients(context.Background(), recipients, fixtures.TestMetadata)
		require.NoError(t, err)
		require.Equal(t, fixtures.Transaction, tx)
	})
}
