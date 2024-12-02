package restyutil_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// mockAuthenticator is a mock implementation of Authenticator interface
type mockAuthenticator struct{}

// Authenticate is a mock implementation of Authenticator interface
func (m *mockAuthenticator) Authenticate(r *resty.Request) error {
	return nil
}

// TestNewHTTPClient_OnAfterResponse tests the OnAfterResponse callback of NewHTTPClient
func TestNewHTTPClient_OnAfterResponse(t *testing.T) {
	tests := map[string]struct {
		statusCode       int
		responseBody     interface{}
		expectedError    error
		expectedSPVError *models.SPVError
	}{
		"Success Response 200": {
			statusCode:    200,
			responseBody:  map[string]string{"message": "success"},
			expectedError: nil,
		},
		"Client Error 400": {
			statusCode: 400,
			responseBody: models.SPVError{
				Message:    "Invalid request",
				StatusCode: 400,
				Code:       "invalid-request",
			},
			expectedError: models.SPVError{
				Message:    "Invalid request",
				StatusCode: 400,
				Code:       "invalid-request",
			},
		},
		"Server Error 500": {
			statusCode: 500,
			responseBody: models.SPVError{
				Message:    goclienterr.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: 500,
				Code:       "internal-server-error",
			},
			expectedError: models.SPVError{
				Message:    goclienterr.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: 500,
				Code:       "internal-server-error",
			},
		},
	}

	// Mock configuration
	cfg := config.Config{
		Addr:      "http://mock-api",
		Timeout:   5,
		Transport: httpmock.DefaultTransport, // Use httpmock
	}

	// Create HTTP client with mock authenticator
	client := restyutil.NewHTTPClient(cfg, &mockAuthenticator{})
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Mock HTTP response
			httpmock.RegisterResponder(
				"GET",
				"http://mock-api/test",
				httpmock.NewJsonResponderOrPanic(tc.statusCode, tc.responseBody),
			)

			// Make request
			resp, err := client.R().Get("/test")

			// Assert errors
			require.ErrorIs(t, err, tc.expectedError)
			require.NotNil(t, resp)

		})
	}
}
