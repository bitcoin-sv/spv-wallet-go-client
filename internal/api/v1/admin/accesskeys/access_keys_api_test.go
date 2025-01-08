package accesskeys_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/accesskeys/accesskeystest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const accessKeysURL = "/api/v1/admin/users/keys"

func TestAccessKeyAPI_AccessKeys(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.AccessKeyPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/users/keys response: 200": {
			expectedResponse: accesskeystest.ExpectedAccessKeyPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("accesskeystest/get_access_keys_200.json"),
		},
		"HTTP GET /api/v1/admin/users/keys response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/users/keys response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/users/keys str response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
	}

	url := testutils.FullAPIURL(t, accessKeysURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			opts := []queries.QueryOption[filter.AdminAccessKeyFilter]{
				queries.QueryWithPageFilter[filter.AdminAccessKeyFilter](filter.Page{
					Number: 1,
					Size:   1,
					Sort:   "asc",
					SortBy: "key",
				}),
				queries.QueryWithFilter(filter.AdminAccessKeyFilter{
					AccessKeyFilter: filter.AccessKeyFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
						},
					},
				}),
			}
			params := "page=1&size=1&sort=asc&sortBy=key&includeDeleted=true"
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponderWithQuery(http.MethodGet, url, params, tc.responder)

			// when:
			got, err := wallet.AccessKeys(context.Background(), opts...)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
