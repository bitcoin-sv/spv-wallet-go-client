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
	ID := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 200", ID): {
			expectedResponse: transactionstest.ExpectedTransactionWithUpdatedMetadata(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_update_metadata_200.json")),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", ID): {
			expectedErr: transactionstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPatch, url, tc.responder)

			// when:
			got, err := spvWalletClient.UpdateTransactionMetadata(context.Background(), &commands.UpdateTransactionMetadata{
				ID: ID,
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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_record_201.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: transactionstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_draft_200.json")),
		},
		"HTTP POST /api/v1/transactions/drafts response: 400": {
			expectedErr: transactionstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
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
	ID := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s  response: 200", ID): {
			expectedResponse: transactionstest.ExpectedTransaction(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transaction_200.json")),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s response: 400", ID): {
			expectedErr: transactionstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PATCH /api/v1/transactions/%s str response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/transactions/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.Transaction(context.Background(), ID)

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
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/transactions_200.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: transactionstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, transactionstest.NewBadRequestSPVError()),
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
