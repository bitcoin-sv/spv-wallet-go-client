package utxos_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos/utxostest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUTXOAPI_UTXOs(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.UtxosPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/utxos response: 200": {
			expectedResponse: utxostest.ExpectedUtxosPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("utxostest/get_utxos_200.json")),
		},
		"HTTP GET /api/v1/utxos response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/utxos str response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/utxos"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.UTXOs(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
