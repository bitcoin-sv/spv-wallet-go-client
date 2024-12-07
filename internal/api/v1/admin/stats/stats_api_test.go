package stats_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/stats/statstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestStatsAPI_Stats(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *models.AdminStats
		expectedErr      error
	}{
		"HTTP GET /v1/admin/stats response: 200": {
			expectedResponse: statstest.ExpectedStatsResponse(),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("statstest/get_stats_200.json")),
		},
		"HTTP GET /v1/admin/stats response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /v1/admin/stats response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /v1/admin/stats str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/v1/admin/stats"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Stats(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
