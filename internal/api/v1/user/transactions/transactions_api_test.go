package transactions_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions/transactionstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	transactionsURL     = "/api/v1/transactions"
	transactionDraftURL = "/api/v1/transactions/drafts"
)

func TestTransactionsAPI_SendToRecipients(t *testing.T) {
	drafTransactionURL := testutils.FullAPIURL(t, transactionDraftURL)
	recordTransactionURL := testutils.FullAPIURL(t, transactionsURL)
	opReturn := &response.OpReturn{StringParts: []string{"hello", "world"}}

	t.Run("SendToRecipients success", func(t *testing.T) {
		// given:
		wallet, transport := testutils.GivenSPVUserAPI(t)
		transport.RegisterResponder(http.MethodPost, drafTransactionURL, testutils.NewJSONFileResponderWithStatusOK("transactionstest/transaction_draft_with_hex_200.json"))
		transport.RegisterResponder(http.MethodPost, recordTransactionURL, testutils.NewJSONFileResponderWithStatusOK("transactionstest/transaction_send_to_recipients_200.json"))
		ctx := context.Background()

		// when:
		result, err := wallet.SendToRecipients(ctx, &commands.SendToRecipients{
			Recipients: []*commands.Recipients{
				{
					OpReturn: opReturn,
				},
			},
		})

		// then:
		require.ErrorIs(t, err, nil)
		require.Equal(t, transactionstest.ExpectedSendToRecipientsTransaction(t), result)
	})

	t.Run("SendToRecipients - DraftToRecipients error", func(t *testing.T) {
		// given:
		wallet, transport := testutils.GivenSPVUserAPI(t)
		transport.RegisterResponder(http.MethodPost, drafTransactionURL, testutils.NewBadRequestSPVErrorResponder())
		ctx := context.Background()

		// when:
		result, err := wallet.SendToRecipients(ctx, &commands.SendToRecipients{
			Recipients: []*commands.Recipients{
				{
					OpReturn: opReturn,
				},
			},
		})

		// then:
		require.ErrorIs(t, err, testutils.NewBadRequestSPVError())
		require.Nil(t, result)
	})

	t.Run("SendToRecipients - FinalizeTransaction error", func(t *testing.T) {
		// given:
		wallet, transport := testutils.GivenSPVUserAPI(t)
		transport.RegisterResponder(http.MethodPost, drafTransactionURL, testutils.NewJSONBodyResponderWithStatusOK(transactionstest.ExpectedDraftTransactionWithWrongHex(t)))
		ctx := context.Background()

		// when:
		result, err := wallet.SendToRecipients(ctx, &commands.SendToRecipients{
			Recipients: []*commands.Recipients{
				{
					OpReturn: opReturn,
				},
			},
		})

		// then:
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("SendToRecipients - RecordTransaction error", func(t *testing.T) {
		// given:
		wallet, transport := testutils.GivenSPVUserAPI(t)
		transport.RegisterResponder(http.MethodPost, drafTransactionURL, testutils.NewJSONFileResponderWithStatusOK("transactionstest/transaction_draft_with_hex_200.json"))
		transport.RegisterResponder(http.MethodPost, recordTransactionURL, testutils.NewBadRequestSPVErrorResponder())
		ctx := context.Background()

		// when:
		result, err := wallet.SendToRecipients(ctx, &commands.SendToRecipients{
			Recipients: []*commands.Recipients{
				{
					OpReturn: opReturn,
				},
			},
		})

		// then:
		require.ErrorIs(t, err, testutils.NewBadRequestSPVError())
		require.Nil(t, result)
	})
}

func TestTransactionsAPI_FinalizeTransaction(t *testing.T) {
	tests := map[string]struct {
		draft       *response.DraftTransaction
		expectedHex string
		expectedErr error
	}{
		"Finalize Transaction with proper draft": {
			draft:       transactionstest.ExpectedDraftTransactionWithHex(t),
			expectedHex: "01000000014c037d55e72d2ee6a95ff67bd758c4cee9c7545bb4d72ba77584152fcfa07012000000006b483045022100a01c25ad9a306f747d90a6d0e815795416ee1f004f865b0653ae3eb2939f42d90220110d994aa99f10533d2566317f55cab838b40f333bf4cdf30c82246461c31fef412102af82c4f5cac25cb5062364937c5e2286094b709610e60b7997b6715784dbf91effffffff0200000000000000000e006a0568656c6c6f05776f726c6408000000000000001976a914702cef80a7039a1aebb70dc05ce1e439646fa33788ac00000000",
		},
		"Finalize Transaction fail to parse hex": {
			draft:       transactionstest.ExpectedDraftTransactionWithWrongHex(t),
			expectedErr: errors.ErrFailedToParseHex,
		},
		"Finalize Transaction fail to prepare locking script": {
			draft:       transactionstest.ExpectedDraftTransactionWithWrongLockingScript(t),
			expectedErr: errors.ErrCreateLockingScript,
		},
		"Finalize Transaction fail to add inputs to transaction": {
			draft:       transactionstest.ExpectedDraftTransactionWithWrongInputs(t),
			expectedErr: errors.ErrAddInputsToTransaction,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			//given:
			wallet, _ := testutils.GivenSPVUserAPI(t)

			//when:
			hex, err := wallet.FinalizeTransaction(tc.draft)

			//then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedHex, hex)
		})
	}
}

func TestTransactionsAPI_UpdateTransactionMetadata(t *testing.T) {
	id := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 200", id): {
			expectedResponse: transactionstest.ExpectedTransactionWithUpdatedMetadata(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("transactionstest/patch_transaction_update_metadata_200.json"),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, transactionsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPatch, url, tc.responder)

			// when:
			got, err := wallet.UpdateTransactionMetadata(context.Background(), &commands.UpdateTransactionMetadata{
				ID: id,
				Metadata: queryparams.Metadata{
					"example_key1": "example_key10_val",
					"example_key2": "example_key20_val",
				},
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_RecordTransaction(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		"HTTP POST /api/v1/transactions response: 201": {
			expectedResponse: transactionstest.ExpectedRecordTransaction(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("transactionstest/post_transaction_record_201.json"),
		},
		"HTTP POST /api/v1/transactions response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/transactions response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/transactions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, transactionsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := wallet.RecordTransaction(context.Background(), &commands.RecordTransaction{})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_DraftTransaction(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.DraftTransaction
		expectedErr      error
	}{
		"HTTP POST /api/v1/transactions/drafts response: 200": {
			expectedResponse: transactionstest.ExpectedDraftTransaction(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("transactionstest/post_transaction_draft_200.json"),
		},
		"HTTP POST /api/v1/transactions/drafts response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/transactions/drafts response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/transactions/drafts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, transactionDraftURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := wallet.DraftTransaction(context.Background(), &commands.DraftTransaction{
				Config:   response.TransactionConfig{},
				Metadata: map[string]any{},
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_Transaction(t *testing.T) {
	id := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s  response: 200", id): {
			expectedResponse: transactionstest.ExpectedTransaction(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("transactionstest/get_transaction_200.json"),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, transactionsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.Transaction(context.Background(), id)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_Transactions(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PageModel[response.Transaction]
		expectedErr      error
	}{
		"HTTP GET /api/v1/transactions response: 200": {
			expectedResponse: transactionstest.ExpectedTransactionsPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("transactionstest/get_transactions_200.json"),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/transactions response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/transactions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, transactionsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			opts := []queries.QueryOption[filter.TransactionFilter]{
				queries.QueryWithPageFilter[filter.TransactionFilter](filter.Page{
					Number: 1,
					Size:   1,
					Sort:   "asc",
					SortBy: "key",
				}),
				queries.QueryWithFilter(filter.TransactionFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
					},
				}),
			}
			params := "page=1&size=1&sort=asc&sortBy=key&includeDeleted=true"
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponderWithQuery(http.MethodGet, url, params, tc.responder)

			// when:
			got, err := wallet.Transactions(context.Background(), opts...)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
