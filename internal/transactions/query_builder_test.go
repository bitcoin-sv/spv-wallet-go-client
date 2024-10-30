package transactions_test

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/transactions"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder(t *testing.T) {
	type filters struct {
		TransactionFilter filter.TransactionFilter
		QueryParamsFilter filter.QueryParams
		MetadataFilter    transactions.Metadata
	}
	tests := map[string]struct {
		filters        filters
		expectedParams url.Values
		expectedErr    error
		builder        transactions.FilterQueryBuilder
	}{
		"query bilder: empty filters": {
			filters:        filters{},
			expectedParams: make(url.Values),
		},
		"query builder: HTTP GET transactions query with query params filter-only": {
			filters: filters{
				QueryParamsFilter: filter.QueryParams{
					Page:          10,
					PageSize:      20,
					OrderByField:  "id",
					SortDirection: "asc",
				},
			},
			expectedParams: url.Values{
				"page":   []string{"10"},
				"size":   []string{"20"},
				"sortBy": []string{"id"},
				"sort":   []string{"asc"},
			},
		},
		"query builder: HTTP GET transactions query with transaction filter-only": {
			filters: filters{
				TransactionFilter: filter.TransactionFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: ptr(true),
						CreatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedParams: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"query builder: HTTP GET transactions query with metadata filter-only": {
			expectedParams: url.Values{
				"metadata[key1]": []string{"value1"},
				"metadata[key2]": []string{"1024"},
			},
			filters: filters{
				MetadataFilter: transactions.Metadata{
					"key1": "value1",
					"key2": 1024,
				},
			},
		},
		"query builder: HTTP GET transactions query all filters set": {
			filters: filters{
				QueryParamsFilter: filter.QueryParams{
					Page:          10,
					PageSize:      20,
					OrderByField:  "id",
					SortDirection: "asc",
				},
				TransactionFilter: filter.TransactionFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: ptr(true),
						CreatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
				MetadataFilter: transactions.Metadata{
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
		"query builder: filter query builder failure": {
			filters: filters{
				QueryParamsFilter: filter.QueryParams{
					Page:          10,
					PageSize:      20,
					OrderByField:  "id",
					SortDirection: "asc",
				},
			},
			builder:     &filterQueryBuilderFailureStub{},
			expectedErr: transactions.ErrFilterQueryBuilder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []transactions.QueryBuilderOption{
				transactions.WithMetadataFilter(tc.filters.MetadataFilter),
				transactions.WithQueryParamsFilter(tc.filters.QueryParamsFilter),
				transactions.WithTransactionFilter(tc.filters.TransactionFilter),
				transactions.WithFilterQueryBuilder(tc.builder),
			}
			qb := transactions.NewQueryBuilder(opts...)
			got, err := qb.Build()
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

type filterQueryBuilderFailureStub struct{}

func (f *filterQueryBuilderFailureStub) Build() (url.Values, error) {
	return nil, errors.New("filter query builder failure stub - query build op failure")
}
