package transactions_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions/transactionstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTransactionsAPI_UpdateTransactionMetadata(t *testing.T) {
	id := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 200", id): {
			expectedResponse: transactionstest.ExpectedTransactionWithUpdatedMetadata(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/patch_transaction_update_metadata_200.json")),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPatch, url, tc.responder)

			// when:
			got, err := spvWalletClient.UpdateTransactionMetadata(context.Background(), &commands.UpdateTransactionMetadata{
				ID: id,
				Metadata: querybuilders.Metadata{
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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/post_transaction_record_201.json")),
		},
		"HTTP POST /api/v1/transactions response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/transactions response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP POST /api/v1/transactions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := spvWalletClient.RecordTransaction(context.Background(), &commands.RecordTransaction{})

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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/post_transaction_draft_200.json")),
		},
		"HTTP POST /api/v1/transactions/drafts response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/transactions/drafts response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP POST /api/v1/transactions/drafts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions/drafts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := spvWalletClient.DraftTransaction(context.Background(), &commands.DraftTransaction{
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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/get_transaction_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/transactions/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.Transaction(context.Background(), id)

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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/get_transactions_200.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/transactions response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/transactions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.Transactions(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
