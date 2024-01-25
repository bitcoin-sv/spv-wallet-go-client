package buxclient

import (
	"context"
	"net/http"
	"testing"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/libsv/go-bt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BuxOrg/go-buxclient/fixtures"
	"github.com/BuxOrg/go-buxclient/transports"
)

// TestTransactions will test the Transaction methods
func TestTransactions(t *testing.T) {
	t.Run("GetTransaction", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction",
			Result:    fixtures.MarshallForTestHandler(fixtures.Transaction),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)

		// when
		tx, err := client.GetTransaction(context.Background(), fixtures.Transaction.ID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("GetTransactions", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction/search",
			Result:    fixtures.MarshallForTestHandler([]*buxmodels.Transaction{fixtures.Transaction}),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)
		conditions := map[string]interface{}{
			"fee": map[string]interface{}{
				"$lt": 100,
			},
			"total_value": map[string]interface{}{
				"$lt": 740,
			},
		}

		// when
		txs, err := client.GetTransactions(context.Background(), conditions, fixtures.TestMetadata, nil)

		// then
		assert.NoError(t, err)
		assert.Equal(t, []*buxmodels.Transaction{fixtures.Transaction}, txs)
	})

	t.Run("GetTransactionsCount", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction/count",
			Result:    "1",
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)
		conditions := map[string]interface{}{
			"fee": map[string]interface{}{
				"$lt": 100,
			},
			"total_value": map[string]interface{}{
				"$lt": 740,
			},
		}

		// when
		count, err := client.GetTransactionsCount(context.Background(), conditions, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("RecordTransaction", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction/record",
			Result:    fixtures.MarshallForTestHandler(fixtures.Transaction),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)

		// when
		tx, err := client.RecordTransaction(context.Background(), fixtures.Transaction.Hex, "", fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("UpdateTransactionMetadata", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction",
			Result:    fixtures.MarshallForTestHandler(fixtures.Transaction),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)

		// when
		tx, err := client.UpdateTransactionMetadata(context.Background(), fixtures.Transaction.ID, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("SendToRecipients", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type: fixtures.RequestType,
			Queries: []*testTransportHandlerRequest{
				{
					Path: "/transaction/record",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, fixtures.MarshallForTestHandler(fixtures.Transaction))
					},
				},
				{
					Path: "/transaction",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, fixtures.MarshallForTestHandler(fixtures.DraftTx))
					},
				},
			},
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)
		recipients := []*transports.Recipients{{
			OpReturn: fixtures.DraftTx.Configuration.Outputs[0].OpReturn,
			Satoshis: 1000,
			To:       fixtures.Destination.Address,
		}}

		// when
		tx, err := client.SendToRecipients(context.Background(), recipients, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.Transaction, tx)
	})

	t.Run("SendToRecipients - nil draft", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type: fixtures.RequestType,
			Queries: []*testTransportHandlerRequest{
				{
					Path: "/transaction/record",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, fixtures.MarshallForTestHandler(fixtures.Transaction))
					},
				},
				{
					Path: "/transaction",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, "nil")
					},
				},
			},
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)
		recipients := []*transports.Recipients{{
			OpReturn: fixtures.DraftTx.Configuration.Outputs[0].OpReturn,
			Satoshis: 1000,
			To:       fixtures.Destination.Address,
		}}

		// when
		tx, err := client.SendToRecipients(context.Background(), recipients, fixtures.TestMetadata)

		// then
		assert.Error(t, err)
		assert.Nil(t, tx)
	})

	t.Run("SendToRecipients - FinalizeTransaction error", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type: fixtures.RequestType,
			Queries: []*testTransportHandlerRequest{
				{
					Path: "/transaction/record",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, fixtures.MarshallForTestHandler(fixtures.Transaction))
					},
				},
				{
					Path: "/transaction",
					Result: func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						mustWrite(w, fixtures.MarshallForTestHandler(buxmodels.DraftTransaction{}))
					},
				},
			},
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)
		recipients := []*transports.Recipients{{
			OpReturn: fixtures.DraftTx.Configuration.Outputs[0].OpReturn,
			Satoshis: 1000,
			To:       fixtures.Destination.Address,
		}}

		// when
		tx, err := client.SendToRecipients(context.Background(), recipients, fixtures.TestMetadata)

		// then
		assert.Error(t, err)
		assert.Nil(t, tx)
	})

	t.Run("FinalizeTransaction", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/transaction/record",
			Result:    fixtures.MarshallForTestHandler(fixtures.Transaction),
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)

		// when
		signedHex, err := client.FinalizeTransaction(fixtures.DraftTx)

		txDraft, btErr := bt.NewTxFromString(signedHex)
		require.NoError(t, btErr)

		// then
		assert.NoError(t, err)
		assert.Len(t, txDraft.Inputs, len(fixtures.DraftTx.Configuration.Inputs))
		assert.Len(t, txDraft.Outputs, len(fixtures.DraftTx.Configuration.Outputs))
	})

	t.Run("UnreserveUtxos", func(t *testing.T) {
		// given
		transportHandler := testTransportHandler{
			Type:      fixtures.RequestType,
			Path:      "/utxo/unreserve",
			ClientURL: fixtures.ServerURL,
			Client:    WithHTTPClient,
		}
		client := getTestBuxClient(transportHandler, false)

		// when
		err := client.UnreserveUtxos(context.Background(), fixtures.DraftTx.Configuration.Outputs[0].PaymailP4.ReferenceID)

		// then
		assert.NoError(t, err)
	})
}

// TestDraftTransactions will test the DraftTransaction methods
func TestDraftTransactions(t *testing.T) {
	transportHandler := testTransportHandler{
		Type:      fixtures.RequestType,
		Path:      "/transaction",
		Result:    fixtures.MarshallForTestHandler(fixtures.DraftTx),
		ClientURL: fixtures.ServerURL,
		Client:    WithHTTPClient,
	}

	t.Run("DraftToRecipients", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, false)

		recipients := []*transports.Recipients{{
			OpReturn: fixtures.DraftTx.Configuration.Outputs[0].OpReturn,
			Satoshis: 1000,
			To:       fixtures.Destination.Address,
		}}

		// when
		draft, err := client.DraftToRecipients(context.Background(), recipients, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.DraftTx, draft)
	})

	t.Run("DraftTransaction", func(t *testing.T) {
		// given
		client := getTestBuxClient(transportHandler, false)

		// when
		draft, err := client.DraftTransaction(context.Background(), &fixtures.DraftTx.Configuration, fixtures.TestMetadata)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixtures.DraftTx, draft)
	})
}
