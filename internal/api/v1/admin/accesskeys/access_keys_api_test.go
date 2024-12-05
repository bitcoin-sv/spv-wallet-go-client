package accesskeys_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/accesskeys/accesskeystest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyAPI_AccessKeys(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.AccessKeyPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/users/keys response: 200": {
			expectedResponse: accesskeystest.ExpectedAccessKeyPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("accesskeystest/get_access_keys_200.json")),
		},
		"HTTP GET /api/v1/admin/users/keys response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/users/keys response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/admin/users/keys str response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/users/keys"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// when:
			got, err := wallet.AccessKeys(context.Background(), queries.AdminAccessKeyQueryWithPageFilter(filter.Page{
				Size: 1,
			}))

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
