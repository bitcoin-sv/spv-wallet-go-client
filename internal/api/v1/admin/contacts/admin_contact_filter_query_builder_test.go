package contacts

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestAdminContactFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminContactFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"admin contact filter: zero values": {
			filter:         filter.AdminContactFilter{},
			expectedParams: url.Values{},
		},
		"admin contact filter: only 'xPubId' field set": {
			filter: filter.AdminContactFilter{
				XPubID: ptr("xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"),
			},
			expectedParams: url.Values{
				"xpubId": []string{"xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"},
			},
		},
		"admin contact filter: all fields set": {
			filter: filter.AdminContactFilter{
				ContactFilter: filter.ContactFilter{
					ID:       ptr("b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"),
					FullName: ptr("John Doe"),
					Paymail:  ptr("test@example.com"),
					PubKey:   ptr("pubKey"),
					Status:   ptr("confirmed"),
				},
				XPubID: ptr("xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"),
			},
			expectedParams: url.Values{
				"id":       []string{"b950f5de-3d3a-40b6-bdf8-c9d60e9e0a0a"},
				"fullName": []string{"John Doe"},
				"paymail":  []string{"test@example.com"},
				"pubKey":   []string{"pubKey"},
				"status":   []string{"confirmed"},
				"xpubId":   []string{"xpub6CUGRUonZSQ4TWtTMmzXdrXDtyPWKi"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := adminContactFilterBuilder{
				contactFilter: tc.filter,
			}
			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T { return &value }
