package status_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestStatusAPI_Status(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse bool
		expectedErr      error
	}{
		"HTTP GET /v1/admin/status response: 200": {
			expectedResponse: true,
			responder:        httpmock.NewStringResponder(http.StatusOK, "true"),
		},
		"HTTP GET /v1/admin/status response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /v1/admin/status response: 401": {
			expectedResponse: false,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized, spvwallettest.NewUnauthorizedAccessSPVError()),
		},
		"HTTP GET /v1/admin/status response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /v1/admin/status str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/v1/admin/status"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Status(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
