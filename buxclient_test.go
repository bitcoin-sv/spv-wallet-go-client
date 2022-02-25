package buxclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/utils"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	adminKeyXpub     = "xprv9s21ZrQH143K4Z8JnrQ7XsYxzKbFNsAEPyHMaMU2fbMtoY1YmsJLFo3XBkg2m7e9UJLS6xvd2HjZ5WN9fQbMSGU7uXEE2pksvbQYCXswLB5"
	xPubID           = "9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36"
	testAddress      = "1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt"
	testAddress2     = "1PnRDRF517hhrFJ5VvR7QGcpQhc7qRshFA"
	xPrivString      = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	xPubString       = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	serverURL        = "https://example.com/"
	xpubJSON         = `{"data":{"xpub":{"id":"0092de4d2aafa59a71a1f90342c138e1c4f19cd1b10e2d17422b34a1d06733e0"}}}`
	txID             = "041479f86c475603fd510431cf702bc8c9849a9c350390eb86b467d82a13cc24"
	draftTxJSON      = `{"created_at":"2022-02-09T16:28:39.000639Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"id":"fe6fe12c25b81106b7332d58fe87dab7bc6e56c8c21ca45b4de05f673f3f653c","hex":"010000000141e3be4d5a3f25e11157bfdd100e7c3497b9be2b80b57eb55e5376b075e7dc5d0200000000ffffffff02e8030000000000001976a9147ff514e6ae3deb46e6644caac5cdd0bf2388906588ac170e0000000000001976a9143dbdb346aaf1c3dc501a2f8c186c3d3e8a87764588ac00000000","xpub_id":"9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36","expires_at":"2022-02-09T16:29:08.991801Z","metadata":{"testkey":"test-value"},"configuration":{"change_destinations":[{"created_at":"2022-02-09T16:28:38.997313Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"id":"252e8a915a5f05effab827a887e261a2416a76f3d3aada946a70a575c0bb76a7","xpub_id":"9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36","locking_script":"76a9143dbdb346aaf1c3dc501a2f8c186c3d3e8a87764588ac","type":"pubkeyhash","chain":1,"num":100,"address":"16dTUJwi7qT3JqzAUMcDHaVV3sB4fH85Ep","draft_id":"fe6fe12c25b81106b7332d58fe87dab7bc6e56c8c21ca45b4de05f673f3f653c"}],"change_destinations_strategy":"","change_minimum_satoshis":0,"change_number_of_destinations":0,"change_satoshis":3607,"expires_in":0,"fee":97,"fee_unit":{"satoshis":1,"bytes":2},"from_utxos":null,"inputs":[{"created_at":"2022-01-28T13:45:02.352Z","updated_at":"2022-02-09T16:28:38.993207Z","deleted_at":null,"id":"efe383eea1a6f7925afb2621b69ea9ba6bd0623e8d61827bad994f8be85161fc","transaction_id":"5ddce775b076535eb57eb5802bbeb997347c0e10ddbf5711e1253f5a4dbee341","xpub_id":"9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36","output_index":2,"satoshis":4704,"script_pub_key":"76a914c746bf0f295375cbea4a5ef25b36c84ff9801bac88ac","type":"pubkeyhash","draft_id":"fe6fe12c25b81106b7332d58fe87dab7bc6e56c8c21ca45b4de05f673f3f653c","reserved_at":"2022-02-09T16:28:38.993205Z","spending_tx_id":null,"destination":{"created_at":"2022-01-28T13:45:02.324Z","updated_at":"0001-01-01T00:00:00Z","metadata":{"client_id":"8","run":90,"run_id":"3108aa426fc7102488bb0ffd","xbench":"destination for testing"},"deleted_at":null,"id":"b8bfa56e37c90f1b25df2e571f727cfec80dd17c5d1845c4b93e21034f7f6a0b","xpub_id":"9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36","locking_script":"76a914c746bf0f295375cbea4a5ef25b36c84ff9801bac88ac","type":"pubkeyhash","chain":0,"num":212,"address":"1KAgDiUasnC7roCjQZM1XLJUpq4BYHjdp6","draft_id":""}}],"miner":"","outputs":[{"satoshis":1000,"scripts":[{"address":"1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt","satoshis":1000,"script":"76a9147ff514e6ae3deb46e6644caac5cdd0bf2388906588ac","script_type":"pubkeyhash"}],"to":"1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt","op_return":null},{"satoshis":3607,"scripts":[{"address":"16dTUJwi7qT3JqzAUMcDHaVV3sB4fH85Ep","satoshis":3607,"script":"76a9143dbdb346aaf1c3dc501a2f8c186c3d3e8a87764588ac","script_type":""}],"to":"16dTUJwi7qT3JqzAUMcDHaVV3sB4fH85Ep","op_return":null}],"send_all_to":"","sync":null},"status":"draft"}`
	destinationJSON  = `{"id":"90d10acb85f37dd009238fe7ec61a1411725825c82099bd8432fcb47ad8326ce","xpub_id":"9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36","locking_script":"76a9140e0eb4911d79e9b7683f268964f595b66fa3604588ac","type":"pubkeyhash","chain":0,"num":245,"address":"12HL5RyEy3Rt6SCwxgpiFSTigem1Pzbq22","metadata":{"test":"test value"}}}`
	transactionJSON  = `{"id":"041479f86c475603fd510431cf702bc8c9849a9c350390eb86b467d82a13cc24","created_at":"2022-01-28T13:45:01.711Z","updated_at":null,"deleted_at":null,"hex":"0100000004afcafa163824904aa3bbc403b30db56a08f29ffa53b16b1b4b4914b9bd7d7610010000006a4730440220710c2b2fe5a0ece2cbc962635d0fb6dabf95c94db0b125c3e2613cede9738666022067e9cc0f4f706c3a2781990981a50313fb0aad18c1e19a757125eec2408ecadb412103dcd8d28545c9f80af54648fcca87972d89e3e7ed7b482465dd78b62c784ad533ffffffff783452c4038c46a4d68145d829f09c70755edd8d4b3512d7d6a27db08a92a76b000000006b483045022100ee7e24859274013e748090a022bf51200ab216771b5d0d57c0d074843dfa62bd02203933c2bd2880c2f8257befff44dc19cb1f3760c6eea44fc0f8094ff94bce652a41210375680e36c45658bd9b0694a48f5756298cf95b77f50bada14ef1cba6d7ea1d3affffffff25e893beb8240ede7661c02cb959799d364711ba638eccdf12e3ce60faa2fd0f010000006b483045022100fc380099ac7f41329aaeed364b95baa390be616243b80a8ef444ae0ddc76fa3a0220644a9677d40281827fa4602269720a5a453fbe77409be40293c3f8248534e5f8412102398146eff37de36ed608b2ee917a3d4b4a424722f9a00f1b48c183322a8ef2a1ffffffff00e6f915a5a3678f01229e5c320c64755f242be6cebfac54e2f77ec5e0eec581000000006b483045022100951511f81291ac234926c866f777fe8e77bc00661031675978ddecf159cc265902207a5957dac7c89493e2b7df28741ce3291e19dc8bba4b13082c69d0f2b79c70ab4121031d674b3ad42b28f3a445e9970bd9ae8fe5d3fb89ee32452d9f6dc7916ea184bfffffffff04c7110000000000001976a91483615db3fb9b9cbbf4cd407100833511a1cb278588ac30060000000000001976a914296a5295e70697e844fb4c2113b41a501d41452e88ac96040000000000001976a914e73e21935fc48df0d1cf8b73f2e8bbd23b78244a88ac27020000000000001976a9140b2b03751813e3467a28ce916cbb102d84c6eec588ac00000000","block_hash":"","block_height":0,"fee":354,"number_of_inputs":4,"number_of_outputs":4,"total_value":6955,"metadata":{"client_id":"8","run":76,"run_id":"3108aa426fc7102488bb0ffd","xbench":"is awesome"},"output_value":1725,"direction":"incoming"}`
	transactionsJSON = `[{"id":"caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda","created_at":"2022-01-28T13:44:59.376Z","updated_at":null,"deleted_at":null,"hex":"0100000001cf4faa628ce1abdd2cfc641c948898bb7a3dbe043999236c3ea4436a0c79f5dc000000006a47304402206aeca14175e4477031970c1cda0af4d9d1206289212019b54f8e1c9272b5bac2022067c4d32086146ca77640f02a989f51b3c6738ebfa24683c4a923f647cf7f1c624121036295a81525ba33e22c6497c0b758e6a84b60d97c2d8905aa603dd364915c3a0effffffff023e030000000000001976a914f7fc6e0b05e91c3610efd0ce3f04f6502e2ed93d88ac99030000000000001976a914550e06a3aa71ba7414b53922c13f96a882bf027988ac00000000","block_hash":"","block_height":0,"fee":97,"number_of_inputs":1,"number_of_outputs":2,"total_value":733,"metadata":{"client_id":"8","run":14,"run_id":"3108aa426fc7102488bb0ffd","xbench":"is awesome"},"output_value":921,"direction":"incoming"},{"id":"5f4fd2be162769852e8bd1362bb8d815a89e137707b4985249876a7f0ebbb071","created_at":"2022-01-28T13:44:59.996Z","updated_at":null,"deleted_at":null,"hex":"01000000016c0c005d516ccd1f1029fa5b61be51a0feaee6e2b07804ceba71047e06edb2df000000006b483045022100ab020464941452dff13bf4ff40a6218825b8dc3502d7860857ee0dd9407e490402206325d24bd46c09b246ebe8493257f2b91d4157de58adfdedf42ba72d6de9aaf5412103a06808b0c597ee6c572baf4f167166e9fed4b8ca66d651d2345b12e0ae5344b3ffffffff0208020000000000001976a914c3367acfc659588393c68dae3eb435c5d0a088b988ac46120000000000001976a91492fc673e0630962068c8b7d909fbfeeb77e3ea3288ac00000000","block_hash":"","block_height":0,"fee":97,"number_of_inputs":1,"number_of_outputs":2,"total_value":423,"metadata":{"client_id":"8","run":32,"run_id":"3108aa426fc7102488bb0ffd","xbench":"is awesome"},"output_value":4678,"direction":"incoming"}]`
	accessKeyString  = `7779d24ca6f8821f225042bf55e8f80aa41b08b879b72827f51e41e6523b9cd0`
)

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func mustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

type testTransportHandler struct {
	ClientURL string
	Client    func(serverURL string, httpClient *http.Client) ClientOps
	Path      string
	Queries   []*testTransportHandlerRequest
	Result    string
	Type      string
}

type testTransportHandlerRequest struct {
	Path   string
	Result func(w http.ResponseWriter, req *http.Request)
}

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
			WithHTTP(serverURL),
		)
		require.NoError(t, err)
		assert.IsType(t, BuxClient{}, *client)
	})

	t.Run("valid xPub client", func(t *testing.T) {
		client, err := New(
			WithXPub(xPubString),
			WithHTTP(serverURL),
		)
		require.NoError(t, err)
		assert.IsType(t, BuxClient{}, *client)
	})

	t.Run("valid access keys", func(t *testing.T) {
		client, err := New(
			WithAccessKey(accessKeyString),
			WithHTTP(serverURL),
		)
		require.NoError(t, err)
		assert.IsType(t, BuxClient{}, *client)
	})

	t.Run("valid access key WIF", func(t *testing.T) {
		wifKey, _ := bitcoin.PrivateKeyToWif(accessKeyString)
		client, err := New(
			WithAccessKey(wifKey.String()),
			WithHTTP(serverURL),
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
			WithHTTP(serverURL),
		)
		err := client.SetAdminKey("")
		assert.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		err := client.SetAdminKey(xPrivString)
		assert.NoError(t, err)
	})

	t.Run("invalid with", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithAdminKey("rest"),
			WithHTTP(serverURL),
		)
		assert.Error(t, err)
	})

	t.Run("valid with", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithAdminKey(xPrivString),
			WithHTTP(serverURL),
		)
		assert.NoError(t, err)
	})
}

// TestSetDebug will test the debug setter
func TestSetDebug(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		client.SetDebug(true)
		assert.True(t, client.IsDebug())
	})

	t.Run("false", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		client.SetDebug(false)
		assert.False(t, client.IsDebug())
	})

	t.Run("false", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithDebugging(false),
			WithHTTP(serverURL),
		)
		require.NoError(t, err)
	})
}

// TestSetSignRequest will test the sign request setter
func TestSetSignRequest(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		client.SetSignRequest(true)
		assert.True(t, client.IsSignRequest())
	})

	t.Run("false", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		client.SetSignRequest(false)
		assert.False(t, client.IsSignRequest())
	})

	t.Run("false", func(t *testing.T) {
		_, err := New(
			WithXPriv(xPrivString),
			WithDebugging(false),
			WithHTTP(serverURL),
		)
		require.NoError(t, err)
	})
}

// TestDraftTransaction will test the DraftTransaction method
func TestDraftTransaction(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/transactions/new",
		Result:    draftTxJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"new_transaction":` + draftTxJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("draft transaction "+transportHandler.Type, func(t *testing.T) {
			var client = getTestBuxClient(transportHandler, false)
			config := &bux.TransactionConfig{
				Outputs: []*bux.TransactionOutput{{
					Satoshis: 1000,
					To:       testAddress,
				}},
			}
			metadata := &bux.Metadata{
				"test-key": "test-value",
			}

			draft, err := client.DraftTransaction(context.Background(), config, metadata)
			assert.NoError(t, err)
			checkDraftTransactionOutput(t, draft)
		})
	}
}

// TestRegisterXpub will test the RegisterXpub method
func TestRegisterXpub(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/xpubs",
		Result:    xpubJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"xpub":` + xpubJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("draft transaction "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, true)
			metadata := &bux.Metadata{
				"test-key": "test-value",
			}
			err := client.RegisterXpub(context.Background(), xPubString, metadata)
			assert.NoError(t, err)
		})
	}
}

// TestDraftToRecipients will test the DraftToRecipients method
func TestDraftToRecipients(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/transactions/new",
		Result:    draftTxJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"new_transaction":` + draftTxJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("draft transaction "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			recipients := []*transports.Recipients{{
				Satoshis: 1000,
				To:       testAddress,
			}}
			metadata := &bux.Metadata{
				"test-key": "test-value",
			}

			draft, err := client.DraftToRecipients(context.Background(), recipients, metadata)
			assert.NoError(t, err)
			checkDraftTransactionOutput(t, draft)
		})
	}
}

func checkDraftTransactionOutput(t *testing.T, draft *bux.DraftTransaction) {
	assert.IsType(t, bux.DraftTransaction{}, *draft)
	assert.Equal(t, xPubID, draft.XpubID)
	assert.Equal(t, bux.DraftStatusDraft, draft.Status)
	assert.Len(t, draft.Configuration.Inputs, 1)
	assert.Len(t, draft.Configuration.Outputs, 2)
	assert.Equal(t, uint64(1000), draft.Configuration.Outputs[0].Satoshis)
	assert.Equal(t, "test-value", draft.Metadata["testkey"])
}

// TestGetDestination will test the GetDestination method
func TestGetDestination(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/destinations",
		Result:    destinationJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"destination":` + destinationJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("new destination "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			destination, err := client.GetDestination(context.Background(), nil)
			assert.NoError(t, err)
			assert.IsType(t, bux.Destination{}, *destination)
			assert.Equal(t, "90d10acb85f37dd009238fe7ec61a1411725825c82099bd8432fcb47ad8326ce", destination.ID)
			assert.Equal(t, xPubID, destination.XpubID)
			assert.Equal(t, "76a9140e0eb4911d79e9b7683f268964f595b66fa3604588ac", destination.LockingScript)
			assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)
			assert.Equal(t, uint32(0), destination.Chain)
			assert.Equal(t, uint32(245), destination.Num)
			assert.Equal(t, "12HL5RyEy3Rt6SCwxgpiFSTigem1Pzbq22", destination.Address)
		})
	}
}

// TestGetTransaction will test the GetTransaction method
func TestGetTransaction(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/transaction",
		Result:    transactionJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"transaction":` + transactionJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("get transaction "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			transaction, err := client.GetTransaction(context.Background(), txID)
			assert.NoError(t, err)
			assert.IsType(t, bux.Transaction{}, *transaction)
			assert.Equal(t, txID, transaction.ID)
			assert.Equal(t, uint64(354), transaction.Fee)
			assert.Equal(t, uint32(4), transaction.NumberOfInputs)
			assert.Equal(t, uint32(4), transaction.NumberOfOutputs)
			assert.Equal(t, uint64(6955), transaction.TotalValue)
		})
	}
}

// TestGetTransactions will test the GetTransactions method
func TestGetTransactions(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/transactions",
		Result:    transactionsJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"transactions":` + transactionsJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("get transactions "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			conditions := map[string]interface{}{
				"fee": map[string]interface{}{
					"$lt": 100,
				},
				"total_value": map[string]interface{}{
					"$lt": 740,
				},
			}
			metadata := &bux.Metadata{
				"run_id": "3108aa426fc7102488bb0ffd",
			}
			transactions, err := client.GetTransactions(context.Background(), conditions, metadata)
			assert.NoError(t, err)
			assert.IsType(t, []*bux.Transaction{}, transactions)
			assert.Len(t, transactions, 2)
			assert.Equal(t, "caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda", transactions[0].ID)
			assert.Equal(t, uint64(97), transactions[0].Fee)
			assert.Equal(t, uint64(733), transactions[0].TotalValue)
			assert.Equal(t, "8", transactions[0].Metadata["client_id"])
			assert.Equal(t, "5f4fd2be162769852e8bd1362bb8d815a89e137707b4985249876a7f0ebbb071", transactions[1].ID)
			assert.Equal(t, uint64(97), transactions[1].Fee)
			assert.Equal(t, uint64(423), transactions[1].TotalValue)
			assert.Equal(t, "8", transactions[1].Metadata["client_id"])
		})
	}
}

// TestRecordTransaction will test the RecordTransaction method
func TestRecordTransaction(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type:      "http",
		Path:      "/transactions/record",
		Result:    transactionJSON,
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type:      "graphql",
		Path:      "/graphql",
		Result:    `{"data":{"transaction":` + transactionJSON + `}}`,
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("get transactions "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			hex := ""
			metadata := &bux.Metadata{
				"test-key": "test-value",
			}
			transaction, err := client.RecordTransaction(context.Background(), hex, "", metadata)
			assert.NoError(t, err)
			assert.IsType(t, bux.Transaction{}, *transaction)
			assert.Equal(t, txID, transaction.ID)
		})
	}
}

// TestSendToRecipients will test the SendToRecipients method
func TestSendToRecipients(t *testing.T) {
	transportHandlers := []testTransportHandler{{
		Type: "http",
		Queries: []*testTransportHandlerRequest{{
			Path: "/transactions/new",
			Result: func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				mustWrite(w, draftTxJSON)
			},
		}, {
			Path: "/transactions/record",
			Result: func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				mustWrite(w, transactionJSON)
			},
		}},
		ClientURL: serverURL,
		Client:    WithHTTPClient,
	}, {
		Type: "graphql",
		Queries: []*testTransportHandlerRequest{{
			Path: "/graphql",
			Result: func(w http.ResponseWriter, req *http.Request) {
				result := `{"data":{"transaction":` + transactionJSON + `}}`
				if req.ContentLength > 1000 {
					result = `{"data":{"new_transaction":` + draftTxJSON + `}}`
				}
				w.Header().Set("Content-Type", "application/json")
				mustWrite(w, result)
			},
		}},
		ClientURL: serverURL + `graphql`,
		Client:    WithGraphQLClient,
	}}

	for _, transportHandler := range transportHandlers {
		t.Run("get transactions "+transportHandler.Type, func(t *testing.T) {
			client := getTestBuxClient(transportHandler, false)

			recipients := []*transports.Recipients{{
				To:       testAddress,
				Satoshis: 1234,
			}, {
				To:       testAddress2,
				Satoshis: 4321,
			}}
			metadata := &bux.Metadata{
				"test-key": "test-value",
			}
			transaction, err := client.SendToRecipients(context.Background(), recipients, metadata)
			require.NoError(t, err)
			assert.IsType(t, bux.Transaction{}, *transaction)
			assert.Equal(t, txID, transaction.ID)
		})
	}
}

// TestFinalizeTransaction will test the FinalizeTransaction method
func TestFinalizeTransaction(t *testing.T) {

	t.Run("mock", func(t *testing.T) {
		httpclient := &http.Client{Transport: localRoundTripper{handler: http.NewServeMux()}}
		client, err := New(
			WithXPriv(xPrivString),
			WithHTTPClient(serverURL, httpclient),
		)
		require.NoError(t, err)

		var draft *bux.DraftTransaction
		err = json.Unmarshal([]byte(draftTxJSON), &draft)
		require.NoError(t, err)

		var draftHex string
		draftHex, err = client.FinalizeTransaction(draft)
		require.NoError(t, err)

		var txDraft *bt.Tx
		txDraft, err = bt.NewTxFromString(draftHex)
		require.NoError(t, err)
		assert.Len(t, txDraft.Inputs, 1)
		assert.Len(t, txDraft.GetInputs(), 1)
		assert.Len(t, txDraft.GetOutputs(), 2)
		// todo check the signature
	})
}

// TestGetTransport will test the GetTransport method
func TestGetTransport(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithHTTP(serverURL),
		)
		transport := client.GetTransport()
		assert.IsType(t, &transports.TransportHTTP{}, *transport)
	})

	t.Run("client", func(t *testing.T) {
		client, _ := New(
			WithXPriv(xPrivString),
			WithGraphQL(serverURL),
			WithAdminKey(xPrivString),
			WithDebugging(true),
			WithSignRequest(false),
		)
		transport := client.GetTransport()
		assert.IsType(t, &transports.TransportGraphQL{}, *transport)
	})
}

func getTestBuxClient(transportHandler testTransportHandler, adminKey bool) *BuxClient {
	mux := http.NewServeMux()
	if transportHandler.Queries != nil {
		for _, query := range transportHandler.Queries {
			mux.HandleFunc(query.Path, query.Result)
		}
	} else {
		mux.HandleFunc(transportHandler.Path, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			mustWrite(w, transportHandler.Result)
		})
	}
	httpclient := &http.Client{Transport: localRoundTripper{handler: mux}}

	opts := []ClientOps{
		WithXPriv(xPrivString),
		transportHandler.Client(transportHandler.ClientURL, httpclient),
	}
	if adminKey {
		opts = append(opts, WithAdminKey(adminKeyXpub))
	}

	client, _ := New(opts...)

	return client
}
