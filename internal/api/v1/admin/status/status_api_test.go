package status_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const statusURL = "/v1/admin/status"

func TestStatusAPI_Status(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse bool
		expectedErr      error
	}{
		"HTTP GET /v1/admin/status response: 200": {
			expectedResponse: true,
			responder:        testutils.NewStringResponderStatusOK("true"),
		},
		"HTTP GET /v1/admin/status response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /v1/admin/status response: 401": {
			expectedResponse: false,
			responder:        testutils.NewUnauthorizedAccessSPVErrorResponder(),
		},
		"HTTP GET /v1/admin/status response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /v1/admin/status str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, statusURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Status(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
