package buxclient

import (
	"testing"

	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	xPrivString = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	serverURL   = "https://example.com/"
)

var (
//testTxHex = "020000000165bb8d2733298b2d3b441a871868d6323c5392facf0d3eced3a6c6a17dc84c10000000006a473044022057b101e9a017cdcc333ef66a4a1e78720ae15adf7d1be9c33abec0fe56bc849d022013daa203095522039fadaba99e567ec3cf8615861d3b7258d5399c9f1f4ace8f412103b9c72aebee5636664b519e5f7264c78614f1e57fa4097ae83a3012a967b1c4b9ffffffff03e0930400000000001976a91413473d21dc9e1fb392f05a028b447b165a052d4d88acf9020000000000001976a91455decebedd9a6c2c2d32cf0ee77e2640c3955d3488ac00000000000000000c006a09446f7457616c6c657400000000"
//testTxID  = "1b52eac9d1eb0adf3ce6a56dee1c4768780b8126e288aca65dd1db32f173b853"
)

// TestNewBuxClient will test the TestNewBuxClient method
func TestNewBuxClient(t *testing.T) {
	t.Run("no keys", func(t *testing.T) {
		client, err := New()
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("empty xpriv", func(t *testing.T) {
		client, err := New(
			WithXPriv(""),
		)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("valid client", func(t *testing.T) {
		client, err := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		require.NoError(t, err)
		assert.IsType(t, BuxClient{}, *client)
	})
}

// TestSetAdminKey will test the admin key setter
func TestSetAdminKey(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		err := client.SetAdminKey("")
		assert.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		err := client.SetAdminKey(xPrivString)
		assert.NoError(t, err)
	})

	t.Run("invalid with", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithAdminKey("rest"),
			WithHTTPClient(serverURL),
		)
		assert.Error(t, err)
	})

	t.Run("valid with", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithAdminKey(xPrivString),
			WithHTTPClient(serverURL),
		)
		assert.NoError(t, err)
	})
}

// TestSetDebug will test the debug setter
func TestSetDebug(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		client.SetDebug(true)
	})

	t.Run("false", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		client.SetDebug(false)
	})

	t.Run("false", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithDebugging(false),
			WithHTTPClient(serverURL),
		)
		assert.NoError(t, err)
	})
}

/*
// DraftTransaction will test the DraftTransaction method
func TestDraftTransaction(t *testing.T) {
	t.Run("mock", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
		)

		config := bux.TransactionConfig{
			Outputs: []*bux.TransactionOutput{{
				Satoshis: 1000,
				To:       "test@handcash.io",
			}},
		}

		draft, err := client.DraftTransaction(context.Background(), &config)
		assert.ErrorIs(t, err, errors.New("test error"))
		assert.Nil(t, draft)
	})

	t.Run("mock", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
		)

		config := bux.TransactionConfig{
			Outputs: []*bux.TransactionOutput{{
				Satoshis: 1000,
				To:       "test@handcash.io",
			}},
		}

		draft, err := client.DraftTransaction(context.Background(), &config)
		assert.NoError(t, err)
		assert.IsType(t, bux.DraftTransaction{}, *draft)
		assert.Equal(t, bux.DraftStatusDraft, draft.Status)
	})
}
*/

/*
// TestFinalizeTransaction will test the FinalizeTransaction method
func TestFinalizeTransaction(t *testing.T) {
	t.Run("mock", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
		)
		draft := bux.DraftTransaction{
			Model: bux.Model{},
			TransactionBase: bux.TransactionBase{
				ID:  "transaction_id",
				Hex: "test_hex",
			},
			XpubID:    "xpub_id",
			ExpiresAt: time.Time{}.Add(30 * time.Second),
			Configuration: bux.TransactionConfig{
				Inputs: []*bux.TransactionInput{{
					Utxo:        bux.Utxo{},
					Destination: bux.Destination{},
				}},
				Outputs: nil,
			},
			Status: "draft",
		}

		_, err := client.FinalizeTransaction(&draft)
		assert.NoError(t, err)

		var txDraft *bt.Tx
		txDraft, err = bt.NewTxFromString(draft.Hex)
		assert.NoError(t, err)
		assert.Len(t, txDraft.Inputs, 1)
	})
}
*/

// TestGetTransport will test the GetTransport method
func TestGetTransport(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL),
		)
		transport := client.GetTransport()
		assert.IsType(t, &transports.TransportHTTP{}, *transport)
	})

	t.Run("graphql", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithGraphQLClient(serverURL),
			WithAdminKey(xPrivString),
			WithDebugging(true),
			WithSignRequest(false),
		)
		transport := client.GetTransport()
		assert.IsType(t, &transports.TransportGraphQL{}, *transport)
	})
}
