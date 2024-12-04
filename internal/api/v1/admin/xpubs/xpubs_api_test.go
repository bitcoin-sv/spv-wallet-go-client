package xpubs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/xpubs/xpubstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestXPubsAPI_CreateXPub(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP POST /api/v1/admin/users response: 201": {
			expectedResponse: xpubstest.ExpectedXPub(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusCreated, httpmock.File("xpubstest/post_xpub_201.json")),
		},
		"HTTP POST /api/v1/admin/users response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/admin/users response: 500": {
			expectedErr: xpubstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, xpubstest.NewInternalServerSPVError()),
		},
		"HTTP POST /api/v1/admin/users str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := wallet.CreateXPub(context.Background(), &commands.CreateUserXpub{
				Metadata: map[string]any{},
				XPub:     "",
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
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
			expectedResponse: xpubstest.ExpectedXPubsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("xpubstest/get_xpubs_200.json")),
		},
		"HTTP GET /api/v1/admin/users response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/users response: 500": {
			expectedErr: xpubstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, xpubstest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/admin/users str response: 500": {
			expectedErr: xpubstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, xpubstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.XPubs(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
