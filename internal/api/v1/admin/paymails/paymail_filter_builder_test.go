package paymails

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestPaymailFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminPaymailFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"admin paymail filter: zero values": {
			expectedParams: make(url.Values),
		},
		"admin paymail filter: filter with only 'id' field set": {
			filter: filter.AdminPaymailFilter{
				ID: ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
			},
			expectedParams: url.Values{
				"id": []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
			},
		},
		"admin paymail filter: filter with only 'xPubId' field set": {
			filter: filter.AdminPaymailFilter{
				XpubID: ptr("7d373830-1d74-4c4b-a435-04ce09398027"),
			},
			expectedParams: url.Values{
				"xpubId": []string{"7d373830-1d74-4c4b-a435-04ce09398027"},
			},
		},
		"admin paymail filter: filter with only 'alias' field set": {
			filter: filter.AdminPaymailFilter{
				Alias: ptr("alias"),
			},
			expectedParams: url.Values{
				"alias": []string{"alias"},
			},
		},
		"admin paymail filter: filter with only 'public name' field set": {
			filter: filter.AdminPaymailFilter{
				PublicName: ptr("Alice"),
			},
			expectedParams: url.Values{
				"publicName": []string{"Alice"},
			},
		},
		"admin paymail filter: all fields set": {
			filter: filter.AdminPaymailFilter{
				ID:         ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
				XpubID:     ptr("7d373830-1d74-4c4b-a435-04ce09398027"),
				PublicName: ptr("Alice"),
				Alias:      ptr("alias"),
			},
			expectedParams: url.Values{
				"publicName": []string{"Alice"},
				"xpubId":     []string{"7d373830-1d74-4c4b-a435-04ce09398027"},
				"alias":      []string{"alias"},
				"id":         []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := paymailFilterBuilder{
				paymailFilter:      tc.filter,
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
