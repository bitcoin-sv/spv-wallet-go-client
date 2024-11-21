package configs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestConfigsAPI_SharedConfig_APIResponses(t *testing.T) {
	tests := map[string]struct {
		statusCode       int
		expectedResponse *response.SharedConfig
		expectedErr      error
		responder        httpmock.Responder
	}{
		"HTTP GET /api/v1/configs/shared response: 200": {
			expectedResponse: &response.SharedConfig{
				PaymailDomains: []string{"john.test.4chain.space"},
				ExperimentalFeatures: map[string]bool{
					"pikeContactsEnabled": true,
					"pikePaymentEnabled":  true,
				},
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("configstest/response_200_status_code.json")),
		},
		"HTTP GET /api/v1/configs/shared response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, &models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			}),
		},
		"HTTP GET /api/v1/configs/shared str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/configs/shared"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.SharedConfig(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
