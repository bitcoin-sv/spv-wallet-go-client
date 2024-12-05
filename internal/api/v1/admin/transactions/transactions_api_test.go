package transactions_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/transactions/transactionstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTransactionsAPI_Transaction(t *testing.T) {
	id := "1024"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Transaction
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/admin/transactions/%s  response: 200", id): {
			expectedResponse: transactionstest.ExpectedTransaction(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/get_transaction_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/transactions/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/transactions/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/transactions/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/transactions/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVAdminAPI(t)
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
		"HTTP GET /api/v1/admin/transactions response: 200": {
			expectedResponse: transactionstest.ExpectedTransactionsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionstest/get_transactions_200.json")),
		},
		"HTTP GET /api/v1/admin/transactions response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/transactions response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/admin/transactions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.Transactions(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
