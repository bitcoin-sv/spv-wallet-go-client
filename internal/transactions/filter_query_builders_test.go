package transactions_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/transactions"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestMetadataFilterQueryBuilder(t *testing.T) {
	tests := map[string]struct {
		metadata       transactions.Metadata
		expectedParams url.Values
		expectedErr    error
		depth          int
	}{
		"metadata: empty map": {
			depth:          transactions.DefaultMaxDepth,
			expectedParams: make(url.Values),
		},
		"metadata: map entry [key]=value1": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key]": []string{"value1"},
			},
			metadata: transactions.Metadata{
				"key": "value1",
			},
		},
		"metadata: map entries [key1]=value1, [key2]=1024": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]": []string{"value1"},
				"metadata[key2]": []string{"1024"},
			},
			metadata: transactions.Metadata{
				"key1": "value1",
				"key2": 1024,
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[]{value2,value3,value4}": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":   []string{"value1"},
				"metadata[key2][]": []string{"value2", "value3", "value4"},
			},
			metadata: transactions.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[]{value2, value3, value4}, [key3]=value5, [key4]=[]{value6,value7,value8}": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":   []string{"value1"},
				"metadata[key2][]": []string{"value2", "value3", "value4"},
				"metadata[key3]":   []string{"value5"},
				"metadata[key4][]": []string{"value6", "value7", "value8"},
			},
			metadata: transactions.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
				"key3": "value5",
				"key4": []string{"value6", "value7", "value8"},
			},
		},
		"metadata: map entries [key1]=value1, [key2]=[]{value1,value2,value3,value4}, [key3][key3_nested]=value5, [key4][key4_nested]=[]{6,7}": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
			metadata: transactions.Metadata{
				"key1": "value1",
				"key2": []string{"value2", "value3", "value4"},
				"key3": transactions.Metadata{
					"key3_nested": "value5",
				},
				"key4": transactions.Metadata{
					"key4_nested": []int{6, 7},
				},
			},
		},
		"metadata: 11 map entries, complex nesting, max depth set to 100": {
			depth: transactions.DefaultMaxDepth,
			expectedParams: url.Values{
				"metadata[key1][key2][key3][key1]":                                             []string{"abc"},
				"metadata[key1][key2][key3][key2][key1]":                                       []string{"9"},
				"metadata[key1][key2][key3][key3][key1][key2][key1][]":                         []string{"1", "2", "3", "4"},
				"metadata[key1][key2][key3][key3][key1][key2][key2]":                           []string{"10"},
				"metadata[key1][key2][key3][key3][key1][key2][key3]":                           []string{"abc"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key1]":         []string{"2"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key2]":         []string{"cde"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key3][key1][]": []string{"5", "6", "7", "8"},
				"metadata[key1][key2][key3][key3][key1][key2][key4][key1][key1][key3][key2][]": []string{"a", "b", "c"},
			},
			metadata: transactions.Metadata{
				"key1": transactions.Metadata{
					"key2": transactions.Metadata{
						"key3": transactions.Metadata{
							"key1": "abc",
							"key2": transactions.Metadata{
								"key1": 9,
							},
							"key3": transactions.Metadata{
								"key1": transactions.Metadata{
									"key2": transactions.Metadata{
										"key1": []int{1, 2, 3, 4},
										"key2": 10,
										"key3": "abc",
										"key4": transactions.Metadata{
											"key1": transactions.Metadata{
												"key1": transactions.Metadata{
													"key1": 2,
													"key2": "cde",
													"key3": transactions.Metadata{
														"key1": []int{5, 6, 7, 8},
														"key2": []string{"a", "b", "c"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"metadata: map entries depth exceeded - map entries: 4, max depth: 3": {
			metadata: transactions.Metadata{
				"key1": transactions.Metadata{
					"key2": transactions.Metadata{
						"key3": transactions.Metadata{
							"key4": "value1",
						},
					},
				},
			},
			depth:       3,
			expectedErr: transactions.ErrMetadataFilterMaxDepthExceeded,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mp := transactions.MetadataFilterQueryBuilder{
				MaxDepth: tc.depth,
				Metadata: tc.metadata,
			}
			got, err := mp.Build()
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func TestModelFilterQueryBuilder(t *testing.T) {
	tests := map[string]struct {
		filter         filter.ModelFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"model filter: filter with only 'include deleted field set": {
			expectedParams: url.Values{
				"includeDeleted": []string{"true"},
			},
			filter: filter.ModelFilter{
				IncludeDeleted: ptr(true),
			},
		},
		"model filter: filter with only created range 'from' field set": {
			expectedParams: url.Values{
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					From: ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter wtth only created range 'to' field set": {
			expectedParams: url.Values{
				"createdRange[to]": []string{"2021-01-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					To: ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only created range both fields set": {
			expectedParams: url.Values{
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					From: ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'from' field set": {
			expectedParams: url.Values{
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'to' field set": {
			expectedParams: url.Values{
				"updatedRange[to]": []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					To: ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range both fields set": {
			expectedParams: url.Values{
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: all fields set": {
			expectedParams: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
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
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := transactions.ModelFilterQueryBuilder{ModelFilter: tc.filter}
			got, err := m.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}

func TestQueryParamsFilterQueryBuilder(t *testing.T) {
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
			qp := transactions.QueryParamsFilterQueryBuilder{
				QueryParamsFilter: tc.filter,
			}
			got, err := qp.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func TestTransactionFilterQueryBuilder(t *testing.T) {
	tests := map[string]struct {
		filter         filter.TransactionFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"transaction filter: zero values": {
			filter: filter.TransactionFilter{
				Id:              ptr(""),
				Hex:             ptr(""),
				BlockHash:       ptr(""),
				BlockHeight:     ptr(uint64(0)),
				Fee:             ptr(uint64(0)),
				NumberOfInputs:  ptr(uint32(0)),
				NumberOfOutputs: ptr(uint32(0)),
				DraftID:         ptr(""),
				TotalValue:      ptr(uint64(0)),
				Status:          ptr(""),
				ModelFilter: filter.ModelFilter{
					CreatedRange: &filter.TimeRange{},
					UpdatedRange: &filter.TimeRange{},
				},
			},
			expectedParams: make(url.Values),
		},
		"transaction filter: filter with only 'id' field set": {
			filter: filter.TransactionFilter{
				Id: ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
			},
			expectedParams: url.Values{
				"id": []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
			},
		},
		"transaction filter: filter with only 'hex' field set": {
			filter: filter.TransactionFilter{
				Hex: ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
			},
			expectedParams: url.Values{
				"hex": []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
			},
		},
		"transaction filter: filter with only 'block hash' field set": {
			filter: filter.TransactionFilter{
				BlockHash: ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
			},
			expectedParams: url.Values{
				"blockHash": []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
			},
		},
		"transaction filter: filter with only 'block height' field set": {
			filter: filter.TransactionFilter{
				BlockHeight: ptr(uint64(839376)),
			},
			expectedParams: url.Values{
				"blockHeight": []string{"839376"},
			},
		},
		"transaction filter: filter with only 'fee' field set": {
			filter: filter.TransactionFilter{
				Fee: ptr(uint64(1)),
			},
			expectedParams: url.Values{
				"fee": []string{"1"},
			},
		},
		"transaction filter: filter with only 'number of inputs' field set": {
			filter: filter.TransactionFilter{
				NumberOfInputs: ptr(uint32(10)),
			},
			expectedParams: url.Values{
				"numberOfInputs": []string{"10"},
			},
		},
		"transaction filter: filter with only 'number of outputs' field set": {
			filter: filter.TransactionFilter{
				NumberOfOutputs: ptr(uint32(20)),
			},
			expectedParams: url.Values{
				"numberOfOutputs": []string{"20"},
			},
		},
		"transaction filter: filter with only 'draft id' field set": {
			filter: filter.TransactionFilter{
				DraftID: ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
			},
			expectedParams: url.Values{
				"draftId": []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
			},
		},
		"transaction filter: filter with only 'total value' field set": {
			filter: filter.TransactionFilter{
				TotalValue: ptr(uint64(100000000)),
			},
			expectedParams: url.Values{
				"totalValue": []string{"100000000"},
			},
		},
		"transaction filter: filter with only 'status' field set": {
			filter: filter.TransactionFilter{
				Status: ptr("RECEIVED"),
			},
			expectedParams: url.Values{
				"status": []string{"RECEIVED"},
			},
		},
		"transaction filter: filter with only 'model filter' fields set": {
			filter: filter.TransactionFilter{
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
			expectedParams: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"transaction filter: all fields set": {
			filter: filter.TransactionFilter{
				Id:              ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
				Hex:             ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
				BlockHash:       ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
				BlockHeight:     ptr(uint64(839376)),
				Fee:             ptr(uint64(1)),
				NumberOfInputs:  ptr(uint32(10)),
				NumberOfOutputs: ptr(uint32(20)),
				DraftID:         ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
				TotalValue:      ptr(uint64(100000000)),
				Status:          ptr("RECEIVED"),
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
			expectedParams: url.Values{
				"id":                 []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"hex":                []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
				"blockHash":          []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
				"blockHeight":        []string{"839376"},
				"fee":                []string{"1"},
				"numberOfInputs":     []string{"10"},
				"numberOfOutputs":    []string{"20"},
				"draftId":            []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"totalValue":         []string{"100000000"},
				"status":             []string{"RECEIVED"},
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tfp := transactions.TransactionFilterQueryBuilder{
				ModelFilterQueryBuilder: transactions.ModelFilterQueryBuilder{ModelFilter: tc.filter.ModelFilter},
				TransactionFilter:       tc.filter,
			}
			got, err := tfp.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
