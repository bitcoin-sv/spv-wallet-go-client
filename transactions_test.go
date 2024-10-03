package walletclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/transaction":
			handleTransaction(w, r)
		case "/v1/transaction/search":
			json.NewEncoder(w).Encode([]*models.Transaction{fixtures.Transaction})
		case "/v1/transaction/count":
			json.NewEncoder(w).Encode(1)
		case "/v1/transaction/record":
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(fixtures.Transaction)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewWithXPriv(server.URL, fixtures.XPrivString)
	require.NoError(t, err)
	require.NotNil(t, client.xPriv)

	t.Run("GetTransaction", func(t *testing.T) {
		tx, err := client.GetTransaction(context.Background(), fixtures.Transaction.ID)
		require.NoError(t, err)
		require.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("GetTransactions", func(t *testing.T) {
		conditions := &filter.TransactionFilter{
			Fee:        Optional(uint64(97)),
			TotalValue: Optional(uint64(6955)),
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

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		if err := json.NewEncoder(w).Encode(fixtures.Transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	case http.MethodPatch:
		var input map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "bad request"}); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}
		response := fixtures.Transaction
		if metadata, ok := input["metadata"].(map[string]interface{}); ok {
			response.Metadata = metadata
		}
		if id, ok := input["id"].(string); ok {
			response.ID = id
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
