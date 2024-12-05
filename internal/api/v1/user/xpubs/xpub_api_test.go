package xpubs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/xpubs/xpubstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestXPubAPI_UpdateXPubMetadata(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP PATCH /api/v1/users/current response: 200": {
			expectedResponse: xpubstest.ExpectedUpdatedXPubMetadata(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("xpubstest/patch_xpub_metadata_200.json")),
		},
		"HTTP PATCH /api/v1/users/current response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP PATCH /api/v1/users/current str response: 500": {
			expectedErr: xpubstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, xpubstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPatch, url, tc.responder)

			// when:
			got, err := wallet.UpdateXPubMetadata(context.Background(), &commands.UpdateXPubMetadata{
				Metadata: map[string]any{"example_key": "example_value"},
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestXPubAPI_XPub(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP GET /api/v1/users/current/ response: 200": {
			expectedResponse: xpubstest.ExpectedUserXPub(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("xpubstest/get_xpub_200.json")),
		},
		"HTTP GET /api/v1/users/current response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/users/current str response: 500": {
			expectedErr: xpubstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, xpubstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.XPub(context.Background())

			//  then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
