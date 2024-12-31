package queryparams_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestQueryParser_Parse_AdminUtxosQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.AdminUtxoFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"admin utxos query: with only metadata": {
			query: &queries.Query[filter.AdminUtxoFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"admin utxos query: with only page filter": {
			query: &queries.Query[filter.AdminUtxoFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"admin utxos query: with only model filter": {
			query: &queries.Query[filter.AdminUtxoFilter]{
				Filter: filter.AdminUtxoFilter{
					UtxoFilter: filter.UtxoFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"admin utxos query: all fields set": {
			query: &queries.Query[filter.AdminUtxoFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.AdminUtxoFilter{
					UtxoFilter: filter.UtxoFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC))},
						},
						SpendingTxID:  testutils.Ptr("7539366c-beb2-4405-8597-025bf2dc7cbd"),
						DraftID:       testutils.Ptr("2453797c-4089-4078-8723-5ecb68e70bd7"),
						Type:          testutils.Ptr("0f65e842-decf-4725-8ad9-877634280e4f"),
						ScriptPubKey:  testutils.Ptr("3adec124-32eb-46f1-94f2-4949a86dbe8d"),
						ID:            testutils.Ptr("abb6a871-dd95-4f7a-8090-ca34cff63801"),
						OutputIndex:   testutils.Ptr(uint32(32)),
						Satoshis:      testutils.Ptr(uint64(64)),
						TransactionID: testutils.Ptr("124c2237-9b54-46c4-bf53-3cec86f7e316"),
						ReservedRange: &filter.TimeRange{
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
					XpubID: testutils.Ptr("c4dab549-902e-42fe-97d2-dc556b716f9a"),
				},
			},
			expectedValues: url.Values{
				"xpubId":                        []string{"c4dab549-902e-42fe-97d2-dc556b716f9a"},
				"scriptPubKey":                  []string{"3adec124-32eb-46f1-94f2-4949a86dbe8d"},
				"draftId":                       []string{"2453797c-4089-4078-8723-5ecb68e70bd7"},
				"reservedRange[to]":             []string{"2021-02-02T00:00:00Z"},
				"reservedRange[from]":           []string{"2021-02-01T00:00:00Z"},
				"transactionId":                 []string{"124c2237-9b54-46c4-bf53-3cec86f7e316"},
				"spendingTxId":                  []string{"7539366c-beb2-4405-8597-025bf2dc7cbd"},
				"type":                          []string{"0f65e842-decf-4725-8ad9-877634280e4f"},
				"satoshis":                      []string{"64"},
				"id":                            []string{"abb6a871-dd95-4f7a-8090-ca34cff63801"},
				"outputIndex":                   []string{"32"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_AdminTransactionsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.AdminTransactionFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"admin transactions query: with only metadata": {
			query: &queries.Query[filter.AdminTransactionFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"admin transactions query: with only page filter": {
			query: &queries.Query[filter.AdminTransactionFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"admin transactions query: with only model filter": {
			query: &queries.Query[filter.AdminTransactionFilter]{
				Filter: filter.AdminTransactionFilter{
					TransactionFilter: filter.TransactionFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"admin transactions query: all fields set": {
			query: &queries.Query[filter.AdminTransactionFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.AdminTransactionFilter{
					TransactionFilter: filter.TransactionFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
						Id:              testutils.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
						Hex:             testutils.Ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
						BlockHash:       testutils.Ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
						BlockHeight:     testutils.Ptr(uint64(839376)),
						Fee:             testutils.Ptr(uint64(1)),
						NumberOfInputs:  testutils.Ptr(uint32(10)),
						NumberOfOutputs: testutils.Ptr(uint32(20)),
						DraftID:         testutils.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
						TotalValue:      testutils.Ptr(uint64(100000000)),
						Status:          testutils.Ptr("RECEIVED"),
					},
					XPubID: testutils.Ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
				},
			},
			expectedValues: url.Values{
				"xpubId":                        []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
				"id":                            []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"hex":                           []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
				"blockHash":                     []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
				"blockHeight":                   []string{"839376"},
				"fee":                           []string{"1"},
				"numberOfInputs":                []string{"10"},
				"numberOfOutputs":               []string{"20"},
				"draftId":                       []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"totalValue":                    []string{"100000000"},
				"status":                        []string{"RECEIVED"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_AdminPaymailsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.AdminPaymailFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"admin paymails query: with only metadata": {
			query: &queries.Query[filter.AdminPaymailFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"admin paymails query: with only page filter": {
			query: &queries.Query[filter.AdminPaymailFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"admin paymails query: with only model filter": {
			query: &queries.Query[filter.AdminPaymailFilter]{
				Filter: filter.AdminPaymailFilter{
					PaymailFilter: filter.PaymailFilter{

						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"admin paymails query: all fields set": {
			query: &queries.Query[filter.AdminPaymailFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.AdminPaymailFilter{
					PaymailFilter: filter.PaymailFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
						ID:         testutils.Ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
						PublicName: testutils.Ptr("Alice"),
						Alias:      testutils.Ptr("alias"),
					},
					XpubID: testutils.Ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
				},
			},
			expectedValues: url.Values{
				"publicName":                    []string{"Alice"},
				"alias":                         []string{"alias"},
				"id":                            []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
				"xpubId":                        []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_AdminContactsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.AdminContactFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"admin contacts query: with only metadata": {
			query: &queries.Query[filter.AdminContactFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"admin contacts query: with only page filter": {
			query: &queries.Query[filter.AdminContactFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"admin contacts query: with only model filter": {
			query: &queries.Query[filter.AdminContactFilter]{
				Filter: filter.AdminContactFilter{
					ContactFilter: filter.ContactFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"admin contacts query: all fields set": {
			query: &queries.Query[filter.AdminContactFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.AdminContactFilter{
					ContactFilter: filter.ContactFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
						ID:       testutils.Ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
						FullName: testutils.Ptr("John Doe"),
						Paymail:  testutils.Ptr("test@example.com"),
						PubKey:   testutils.Ptr("pubKey"),
						Status:   testutils.Ptr("confirmed"),
					},
					XPubID: testutils.Ptr("xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"),
				},
			},
			expectedValues: url.Values{
				"id":                            []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
				"fullName":                      []string{"John Doe"},
				"paymail":                       []string{"test@example.com"},
				"pubKey":                        []string{"pubKey"},
				"status":                        []string{"confirmed"},
				"xpubId":                        []string{"xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_AdminAccessKeysQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.AdminAccessKeyFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"admin access keys query: with only metadata": {
			query: &queries.Query[filter.AdminAccessKeyFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"admin access keys query: with only page filter": {
			query: &queries.Query[filter.AdminAccessKeyFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"admin access keys query: with only model filter": {
			query: &queries.Query[filter.AdminAccessKeyFilter]{
				Filter: filter.AdminAccessKeyFilter{
					AccessKeyFilter: filter.AccessKeyFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"admin access keys query: all fields set": {
			query: &queries.Query[filter.AdminAccessKeyFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.AdminAccessKeyFilter{
					XpubID: testutils.Ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
					AccessKeyFilter: filter.AccessKeyFilter{
						RevokedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
							CreatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
							},
							UpdatedRange: &filter.TimeRange{
								From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
								To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
							},
						},
					},
				},
			},
			expectedValues: url.Values{
				"xpubId":                        []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
				"revokedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"revokedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_XpubsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.XpubFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"xpubs query: with only metadata": {
			query: &queries.Query[filter.XpubFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"xpubs query: with only page filter": {
			query: &queries.Query[filter.XpubFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"xpubs query: with only model filter": {
			query: &queries.Query[filter.XpubFilter]{
				Filter: filter.XpubFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"xpubs query: all fields set": {
			query: &queries.Query[filter.XpubFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.XpubFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
					ID:             testutils.Ptr("5505cbc3-b38f-40d4-885f-c53efd84828f"),
					CurrentBalance: testutils.Ptr(uint64(24)),
				},
			},
			expectedValues: url.Values{
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"id":                            []string{"5505cbc3-b38f-40d4-885f-c53efd84828f"},
				"currentBalance":                []string{"24"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_UtxosQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.UtxoFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"utxos query: with only metadata": {
			query: &queries.Query[filter.UtxoFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"utxos query: with only page filter": {
			query: &queries.Query[filter.UtxoFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"utxos query: with only model filter": {
			query: &queries.Query[filter.UtxoFilter]{
				Filter: filter.UtxoFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"utxos query: all fields set": {
			query: &queries.Query[filter.UtxoFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.UtxoFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
					SpendingTxID:  testutils.Ptr("7539366c-beb2-4405-8597-025bf2dc7cbd"),
					DraftID:       testutils.Ptr("2453797c-4089-4078-8723-5ecb68e70bd7"),
					Type:          testutils.Ptr("0f65e842-decf-4725-8ad9-877634280e4f"),
					ScriptPubKey:  testutils.Ptr("3adec124-32eb-46f1-94f2-4949a86dbe8d"),
					ID:            testutils.Ptr("abb6a871-dd95-4f7a-8090-ca34cff63801"),
					OutputIndex:   testutils.Ptr(uint32(32)),
					Satoshis:      testutils.Ptr(uint64(64)),
					TransactionID: testutils.Ptr("124c2237-9b54-46c4-bf53-3cec86f7e316"),
					ReservedRange: &filter.TimeRange{
						To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
			expectedValues: url.Values{
				"scriptPubKey":                  []string{"3adec124-32eb-46f1-94f2-4949a86dbe8d"},
				"draftId":                       []string{"2453797c-4089-4078-8723-5ecb68e70bd7"},
				"reservedRange[to]":             []string{"2021-02-02T00:00:00Z"},
				"reservedRange[from]":           []string{"2021-02-01T00:00:00Z"},
				"transactionId":                 []string{"124c2237-9b54-46c4-bf53-3cec86f7e316"},
				"spendingTxId":                  []string{"7539366c-beb2-4405-8597-025bf2dc7cbd"},
				"type":                          []string{"0f65e842-decf-4725-8ad9-877634280e4f"},
				"satoshis":                      []string{"64"},
				"id":                            []string{"abb6a871-dd95-4f7a-8090-ca34cff63801"},
				"outputIndex":                   []string{"32"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_TransactionsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.TransactionFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"transactions query: with only metadata": {
			query: &queries.Query[filter.TransactionFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"transactions query: with only page filter": {
			query: &queries.Query[filter.TransactionFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"transactions query: with only model filter": {
			query: &queries.Query[filter.TransactionFilter]{
				Filter: filter.TransactionFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"transactions query: all fields set": {
			query: &queries.Query[filter.TransactionFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.TransactionFilter{
					Id:              testutils.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
					Hex:             testutils.Ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
					BlockHash:       testutils.Ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
					BlockHeight:     testutils.Ptr(uint64(839376)),
					Fee:             testutils.Ptr(uint64(1)),
					NumberOfInputs:  testutils.Ptr(uint32(10)),
					NumberOfOutputs: testutils.Ptr(uint32(20)),
					DraftID:         testutils.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
					TotalValue:      testutils.Ptr(uint64(100000000)),
					Status:          testutils.Ptr("RECEIVED"),
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"id":                            []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"hex":                           []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
				"blockHash":                     []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
				"blockHeight":                   []string{"839376"},
				"fee":                           []string{"1"},
				"numberOfInputs":                []string{"10"},
				"numberOfOutputs":               []string{"20"},
				"draftId":                       []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"totalValue":                    []string{"100000000"},
				"status":                        []string{"RECEIVED"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_PaymailsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.PaymailFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"paymails query: with only metadata": {
			query: &queries.Query[filter.PaymailFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"paymails query: with only page filter": {
			query: &queries.Query[filter.PaymailFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"paymails query: with only model filter": {
			query: &queries.Query[filter.PaymailFilter]{
				Filter: filter.PaymailFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"paymails query: all fields set": {
			query: &queries.Query[filter.PaymailFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.PaymailFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
					ID:         testutils.Ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
					PublicName: testutils.Ptr("Alice"),
					Alias:      testutils.Ptr("alias"),
				},
			},
			expectedValues: url.Values{
				"publicName":                    []string{"Alice"},
				"alias":                         []string{"alias"},
				"id":                            []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}

func TestQueryParser_Parse_ContactsQuery(t *testing.T) {
	tests := map[string]struct {
		query          *queries.Query[filter.ContactFilter]
		expectedValues url.Values
		expectedErr    error
	}{
		"contacts query: with only metadata": {
			query: &queries.Query[filter.ContactFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
			},
			expectedValues: url.Values{
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
			},
		},
		"contacts query: with only page filter": {
			query: &queries.Query[filter.ContactFilter]{
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
			},
			expectedValues: url.Values{
				"page":   []string{"1"},
				"size":   []string{"2"},
				"sort":   []string{"asc"},
				"sortBy": []string{"key"},
			},
		},
		"contacts query: with only model filter": {
			query: &queries.Query[filter.ContactFilter]{
				Filter: filter.ContactFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			expectedValues: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"contacts query: all fields set": {
			query: &queries.Query[filter.ContactFilter]{
				Metadata: queryparams.Metadata{
					"key1": "value1",
					"key2": []string{"value2", "value3", "value4"},
					"key3": queryparams.Metadata{
						"key3_nested": "value5",
					},
					"key4": queryparams.Metadata{
						"key4_nested": []int{6, 7},
					},
				},
				PageFilter: filter.Page{
					Number: 1,
					Size:   2,
					Sort:   "asc",
					SortBy: "key",
				},
				Filter: filter.ContactFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
						CreatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
						},
						UpdatedRange: &filter.TimeRange{
							From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
							To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						},
					},
					ID:       testutils.Ptr("e3a1e174-cdf8-4683-b112-e198144eb9d2"),
					FullName: testutils.Ptr("John Doe"),
					Paymail:  testutils.Ptr("john.doe@test.com"),
					Status:   testutils.Ptr("confirmed"),
				},
			},
			expectedValues: url.Values{
				"paymail":                       []string{"john.doe@test.com"},
				"status":                        []string{"confirmed"},
				"id":                            []string{"e3a1e174-cdf8-4683-b112-e198144eb9d2"},
				"fullName":                      []string{"John Doe"},
				"page":                          []string{"1"},
				"size":                          []string{"2"},
				"sort":                          []string{"asc"},
				"sortBy":                        []string{"key"},
				"metadata[key1]":                []string{"value1"},
				"metadata[key2][]":              []string{"value2", "value3", "value4"},
				"metadata[key3][key3_nested]":   []string{"value5"},
				"metadata[key4][key4_nested][]": []string{"6", "7"},
				"includeDeleted":                []string{"true"},
				"createdRange[from]":            []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":              []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]":            []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":              []string{"2021-02-02T00:00:00Z"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser, err := queryparams.NewQueryParser(tc.query)
			require.NoError(t, err)

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedValues, got.Values)
		})
	}
}
