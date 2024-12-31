package utxos_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/utxos/utxostest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const utxosURL = "/api/v1/admin/utxos"

func TestUtxosAPI_UTXOs(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.UtxosPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/utxos response: 200": {
			expectedResponse: utxostest.ExpectedUtxosPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("utxostest/get_utxos_200.json"),
		},
		"HTTP GET /api/v1/admin/utxos response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/utxos response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/utxos str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, utxosURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			opts := []queries.QueryOption[filter.AdminUtxoFilter]{
				queries.QueryWithPageFilter[filter.AdminUtxoFilter](filter.Page{
					Number: 1,
					Size:   1,
					Sort:   "asc",
					SortBy: "key",
				}),
				queries.QueryWithFilter(filter.AdminUtxoFilter{
					UtxoFilter: filter.UtxoFilter{
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
			got, err := wallet.UTXOs(context.Background(), opts...)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}
