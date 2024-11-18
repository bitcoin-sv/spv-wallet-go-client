package utxos_test

import (
	"context"
	"net/http"
	"testing"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos/utxostest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUTXOAPI_UTXOs(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *queries.UtxosPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/utxos response: 200": {
			expectedResponse: utxostest.ExpectedUtxosPage(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("utxostest/get_utxos_200.json")),
		},
		"HTTP GET /api/v1/utxos response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, utxostest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/utxos str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/utxos"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.UTXOs(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
