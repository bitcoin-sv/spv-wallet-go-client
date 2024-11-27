package querybuilders_test

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestQueryParamsFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.QueryParams
		expectedParams url.Values
		expectedErr    error
	}{
		"query params: filter with only 'page' field set": {
			filter: filter.QueryParams{
				Page: 10,
			},
			expectedParams: url.Values{
				"page": []string{"10"},
			},
		},
		"query params: filter with only 'page size' field set": {
			filter: filter.QueryParams{
				PageSize: 20,
			},
			expectedParams: url.Values{
				"size": []string{"20"},
			},
		},
		"query params: filter with only 'order by' field set": {
			filter: filter.QueryParams{
				OrderByField: "value1",
			},
			expectedParams: url.Values{
				"sortBy": []string{"value1"},
			},
		},
		"query params: filter with only 'sort by' field set": {
			filter: filter.QueryParams{
				SortDirection: "asc",
			},
			expectedParams: url.Values{
				"sort": []string{"asc"},
			},
		},
		"query params: all fields set": {
			filter: filter.QueryParams{
				Page:          10,
				PageSize:      20,
				OrderByField:  "value1",
				SortDirection: "asc",
			},
			expectedParams: url.Values{
				"page":   []string{"10"},
				"size":   []string{"20"},
				"sortBy": []string{"value1"},
				"sort":   []string{"asc"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := querybuilders.QueryParamsFilterBuilder{
				QueryParamsFilter: tc.filter,
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
