package querybuilders_test

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestPageFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.Page
		expectedParams url.Values
		expectedErr    error
	}{
		"page filter: filter with only 'number' set": {
			filter: filter.Page{
				Number: 10,
			},
			expectedParams: url.Values{
				"page": []string{"10"},
			},
		},
		"page filter: filter with only 'size' set": {
			filter: filter.Page{
				Size: 20,
			},
			expectedParams: url.Values{
				"size": []string{"20"},
			},
		},
		"page filter: filter with only 'sort' set": {
			filter: filter.Page{
				Sort: "asc",
			},
			expectedParams: url.Values{
				"sort": []string{"asc"},
			},
		},
		"page filter: filter with only 'sortBy' set": {
			filter: filter.Page{
				SortBy: "key",
			},
			expectedParams: url.Values{
				"sortBy": []string{"key"},
			},
		},
		"page filter: all fields set": {
			filter: filter.Page{
				Number: 10,
				Size:   20,
				Sort:   "asc",
				SortBy: "key",
			},
			expectedParams: url.Values{
				"sortBy": []string{"key"},
				"sort":   []string{"asc"},
				"size":   []string{"20"},
				"page":   []string{"10"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			builder := querybuilders.PageFilterBuilder{
				Page: tc.filter,
			}

			// then:
			got, err := builder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
