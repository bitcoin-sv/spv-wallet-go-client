package configs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestConfigsAPI_SharedConfig_APIResponses(t *testing.T) {
	tests := map[string]struct {
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
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("configstest/response_200_status_code.json")),
		},
		"HTTP GET /api/v1/configs/shared response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, &models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			}),
		},
		"HTTP GET /api/v1/configs/shared str response: 500": {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/configs/shared"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.SharedConfig(context.Background())
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
