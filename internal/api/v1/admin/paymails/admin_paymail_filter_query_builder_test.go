package paymails

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TesAdminPaymailFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminPaymailFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"admin paymail filter: zero values": {
			expectedParams: make(url.Values),
		},
		"admin paymail filter: filter with only 'xPubId' field set": {
			filter: filter.AdminPaymailFilter{
				XpubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
			},
			expectedParams: url.Values{
				"xpubId": []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
			},
		},
		"admin paymail filter: all fields set": {
			filter: filter.AdminPaymailFilter{
				XpubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
				PaymailFilter: filter.PaymailFilter{
					ID:         ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
					PublicName: ptr("Alice"),
					Alias:      ptr("alias"),
				},
			},
			expectedParams: url.Values{
				"publicName": []string{"Alice"},
				"alias":      []string{"alias"},
				"id":         []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
				"xpubId":     []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := adminPaymailFilterBuilder{paymailFilter: tc.filter}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T { return &value }
