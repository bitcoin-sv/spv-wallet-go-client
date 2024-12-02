package configs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configs/configstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
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
			expectedErr: configstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, configstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/configs/shared str response: 500": {
			expectedErr: configstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, configstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/configs/shared"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.SharedConfig(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
