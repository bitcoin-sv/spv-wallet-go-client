package users_test

import (
	"context"
	"net/http"
	"testing"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/users/userstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestXPubAPI_UpdateXPubMetadata(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP PATCH /api/v1/users/current response: 200": {
			expectedResponse: userstest.ExpectedUpdatedXPubMetadata(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/patch_xpub_metadata_200.json")),
		},
		"HTTP PATCH /api/v1/users/current response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP PATCH /api/v1/users/current str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPatch, URL, tc.responder)

			// then:
			got, err := wallet.UpdateXPubMetadata(context.Background(), &commands.UpdateXPubMetadata{
				Metadata: map[string]any{
					"example_key": "example_value",
				},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestXPubAPI_XPub(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP GET /api/v1/users/current/ response: 200": {
			expectedResponse: userstest.ExpectedUserXPub(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/get_xpub_200.json")),
		},
		"HTTP GET /api/v1/users/current response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/users/current str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := clienttest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.XPub(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
