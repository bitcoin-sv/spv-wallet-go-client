package users_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/users/userstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestXPubAPI_UpdateXPubMetadata(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
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
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP PATCH /api/v1/users/current str response: 500": {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPatch, url, tc.responder)

			// then:
			got, err := wallet.UpdateXPubMetadata(context.Background(), &commands.UpdateXPubMetadata{
				Metadata: map[string]any{
					"example_key": "example_value",
				},
			})
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestXPubAPI_XPub(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
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
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/users/current str response: 500": {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.XPub(context.Background())
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
