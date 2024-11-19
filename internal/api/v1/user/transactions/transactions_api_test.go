package transactions_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions/transactionstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTransactionsAPI_UpdateTransactionMetadata(t *testing.T) {
	ID := "1024"
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 200", ID): {
			expectedResponse: transactionstest.ExpectedTransactionWithUpdatedMetadata(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_update_metadata_200.json")),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", ID): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", ID): {
			expectedErr: wallet.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/transactions/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPatch, URL, tc.responder)

			// then:
			got, err := spvWalletClient.UpdateTransactionMetadata(context.Background(), &commands.UpdateTransactionMetadata{
				ID: ID,
				Metadata: querybuilders.Metadata{
					"example_key1": "example_key10_val",
					"example_key2": "example_key20_val",
				},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_RecordTransaction(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		"HTTP POST /api/v1/transactions response: 201": {
			statusCode:       http.StatusCreated,
			expectedResponse: transactionstest.ExpectedRecordTransaction(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_record_201.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/transactions str response: 500": {
			expectedErr: wallet.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			got, err := spvWalletClient.RecordTransaction(context.Background(), &commands.RecordTransaction{})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_DraftTransaction(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.DraftTransaction
		expectedErr      error
	}{
		"HTTP POST /api/v1/transactions/drafts response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: transactionstest.ExpectedDraftTransaction(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_draft_200.json")),
		},
		"HTTP POST /api/v1/transactions/drafts response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/transactions/drafts str response: 500": {
			expectedErr: wallet.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/transactions/drafts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			got, err := spvWalletClient.DraftTransaction(context.Background(), &commands.DraftTransaction{
				Config:   response.TransactionConfig{},
				Metadata: map[string]any{},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_Transaction(t *testing.T) {
	ID := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s  response: 200", ID): {
			statusCode:       http.StatusOK,
			expectedResponse: transactionstest.ExpectedTransaction(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_200.json")),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", ID): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", ID): {
			expectedErr: wallet.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/transactions/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := spvWalletClient.Transaction(context.Background(), ID)
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestTransactionsAPI_Transactions(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.PageModel[response.Transaction]
		expectedErr      error
	}{
		"HTTP GET /api/v1/transactions response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: transactionstest.ExpectedTransactionsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transactions_200.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/transactions str response: 500": {
			expectedErr: wallet.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := spvWalletClient.Transactions(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
