package xpubs

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestXPubFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.XpubFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"xpub filter: zero values": {
			expectedParams: make(url.Values),
		},
		"xpub filter: filter with only 'id' field set": {
			filter: filter.XpubFilter{
				ID: ptr("5505cbc3-b38f-40d4-885f-c53efd84828f"),
			},
			expectedParams: url.Values{
				"id": []string{"5505cbc3-b38f-40d4-885f-c53efd84828f"},
			},
		},
		"xpub filter: filter with only 'current balance' field set": {
			filter: filter.XpubFilter{
				CurrentBalance: ptr(uint64(24)),
			},
			expectedParams: url.Values{
				"currentBalance": []string{"24"},
			},
		},
		"xpub filter: all fields set": {
			filter: filter.XpubFilter{
				ID:             ptr("5505cbc3-b38f-40d4-885f-c53efd84828f"),
				CurrentBalance: ptr(uint64(24)),
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
				"id":                 []string{"5505cbc3-b38f-40d4-885f-c53efd84828f"},
				"currentBalance":     []string{"24"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := xpubFilterBuilder{
				xpubFilter:         tc.filter,
				modelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}
