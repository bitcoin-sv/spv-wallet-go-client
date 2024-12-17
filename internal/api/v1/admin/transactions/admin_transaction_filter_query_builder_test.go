package transactions

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestAdminTransactionFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminTransactionFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"admin transaction filter: zero values": {
			expectedParams: make(url.Values),
		},
		"admin transaction: filter with only 'xPubId' field set": {
			filter: filter.AdminTransactionFilter{
				XPubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
			},
			expectedParams: url.Values{
				"xpubid": []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
			},
		},
		"admin transaction filter: all fields set": {
			filter: filter.AdminTransactionFilter{
				XPubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
				TransactionFilter: filter.TransactionFilter{
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
			},
			expectedParams: url.Values{
				"xpubid":             []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
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
			// when:
			queryBuilder := adminTransactionFilterQueryBuilder{transactionFilter: tc.filter}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T { return &value }
