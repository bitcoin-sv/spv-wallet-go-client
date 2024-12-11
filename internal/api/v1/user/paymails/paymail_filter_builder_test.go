package paymails_test

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/paymails"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestPaymailFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.PaymailFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"paymail filter: zero values": {
			expectedParams: make(url.Values),
		},
		"admin paymail filter: filter with only 'id' field set": {
			filter: filter.PaymailFilter{
				ID: testutils.Ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
			},
			expectedParams: url.Values{
				"id": []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
			},
		},
		"admin paymail filter: filter with only 'alias' field set": {
			filter: filter.PaymailFilter{
				Alias: testutils.Ptr("alias"),
			},
			expectedParams: url.Values{
				"alias": []string{"alias"},
			},
		},
		"admin paymail filter: filter with only 'public name' field set": {
			filter: filter.PaymailFilter{
				PublicName: testutils.Ptr("Alice"),
			},
			expectedParams: url.Values{
				"publicName": []string{"Alice"},
			},
		},
		"admin paymail filter: all fields set": {
			filter: filter.PaymailFilter{
				ID:         testutils.Ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
				PublicName: testutils.Ptr("Alice"),
				Alias:      testutils.Ptr("alias"),
			},
			expectedParams: url.Values{
				"publicName": []string{"Alice"},
				"alias":      []string{"alias"},
				"id":         []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := paymails.PaymailFilterBuilder{
				PaymailFilter:      tc.filter,
				ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
