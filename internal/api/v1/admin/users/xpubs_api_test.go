package xpubs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/users/userstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
)

func TestXPubsAPI_CreateXPub(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP POST /api/v1/admin/users response: 201": {
			expectedResponse: userstest.ExpectedXPub(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusCreated, httpmock.File("userstest/post_xpub_201.json")),
		},
		"HTTP POST /api/v1/admin/users response: 400": {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/admin/users str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			got, err := wallet.CreateXPub(context.Background(), &commands.CreateUserXpub{
				Metadata: map[string]any{},
				XPub:     "",
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestXPubsAPI_XPubs(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.XPubPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/users response: 200": {
			expectedResponse: userstest.ExpectedXPubsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/get_xpubs_200.json")),
		},
		"HTTP GET /api/v1/admin/users response: 400": {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/users str response: 500": {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.XPubs(context.Background())
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
