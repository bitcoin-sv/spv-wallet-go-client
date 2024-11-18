package querybuilders_test

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders/querybuilderstest"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Build(t *testing.T) {
	type filters struct {
		MetadataFilter querybuilders.Metadata
		ModelFilter    filter.ModelFilter
		PageFilter     filter.Page
	}
	tests := map[string]struct {
		filters        filters
		expectedParams url.Values
		expectedErr    error
		builder        querybuilders.FilterQueryBuilder
	}{
		"query bilder: empty filters": {
			filters:        filters{},
			expectedParams: make(url.Values),
		},
		"query builder: URL values with page filter-only": {
			filters: filters{
				PageFilter: filter.Page{
					Number: 10,
					Size:   20,
					SortBy: "id",
					Sort:   "asc",
				},
			},
			expectedParams: url.Values{
				"page":   []string{"10"},
				"size":   []string{"20"},
				"sortBy": []string{"id"},
				"sort":   []string{"asc"},
			},
		},
		"query builder: URL values with metadata filter-only": {
			expectedParams: url.Values{
				"metadata[key1]": []string{"value1"},
				"metadata[key2]": []string{"1024"},
			},
			filters: filters{
				MetadataFilter: querybuilders.Metadata{
					"key1": "value1",
					"key2": 1024,
				},
			},
		},
		"query builder: URL values with all filters set": {
			filters: filters{
				PageFilter: filter.Page{
					Number: 10,
					Size:   20,
					Sort:   "asc",
					SortBy: "id",
				},
				ModelFilter: filter.ModelFilter{
					IncludeDeleted: querybuilderstest.Ptr(true),
					CreatedRange: &filter.TimeRange{
						From: querybuilderstest.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
						To:   querybuilderstest.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					UpdatedRange: &filter.TimeRange{
						From: querybuilderstest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
						To:   querybuilderstest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
					},
				},
				MetadataFilter: querybuilders.Metadata{
					"key1": "value1",
					"key2": 1024,
				},
			},
			expectedParams: url.Values{
				"page":               []string{"10"},
				"size":               []string{"20"},
				"sortBy":             []string{"id"},
				"sort":               []string{"asc"},
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":     []string{"value1"},
				"metadata[key2]":     []string{"1024"},
			},
		},
		"query builder: injected dependency filter query builder failure": {
			filters: filters{
				PageFilter: filter.Page{
					Number: 10,
					Size:   20,
					Sort:   "id",
					SortBy: "asc",
				},
			},
			builder:     &filterQueryBuilderFailureStub{},
			expectedErr: querybuilders.ErrFilterQueryBuilder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			opts := []querybuilders.QueryBuilderOption{
				querybuilders.WithMetadataFilter(tc.filters.MetadataFilter),
				querybuilders.WithPageFilter(tc.filters.PageFilter),
				querybuilders.WithModelFilter(tc.filters.ModelFilter),
				querybuilders.WithFilterQueryBuilder(tc.builder),
			}
			builder := querybuilders.NewQueryBuilder(opts...)

			// then:
			got, err := builder.Build()
			require.ErrorIs(t, err, tc.expectedErr)

			if got != nil {
				require.Equal(t, tc.expectedParams, got.Values)
			}
		})
	}
}

type filterQueryBuilderFailureStub struct{}

func (f *filterQueryBuilderFailureStub) Build() (url.Values, error) {
	return nil, errors.New("filter query builder failure stub - query build op failure")
}
