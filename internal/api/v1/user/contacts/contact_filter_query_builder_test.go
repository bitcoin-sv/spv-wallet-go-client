package contacts

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestContactFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.ContactFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"contact filter: zero values": {
			expectedParams: make(url.Values),
		},
		"contact filter: filter with only 'id' field set": {
			filter: filter.ContactFilter{
				ID: contactstest.Ptr("e3a1e174-cdf8-4683-b112-e198144eb9d2"),
			},
			expectedParams: url.Values{
				"id": []string{"e3a1e174-cdf8-4683-b112-e198144eb9d2"},
			},
		},
		"contact filter: filter with only 'full name' field set": {
			filter: filter.ContactFilter{
				FullName: contactstest.Ptr("John Doe"),
			},
			expectedParams: url.Values{
				"fullName": []string{"John Doe"},
			},
		},
		"contact filter: filter with only 'paymail' field set": {
			filter: filter.ContactFilter{
				Paymail: contactstest.Ptr("john.doe@test.com"),
			},
			expectedParams: url.Values{
				"paymail": []string{"john.doe@test.com"},
			},
		},
		"contact filter: filter with only 'status' field set": {
			filter: filter.ContactFilter{
				Status: contactstest.Ptr("confirmed"),
			},
			expectedParams: url.Values{
				"status": []string{"confirmed"},
			},
		},
		"contact filter: filter with all fields set": {
			filter: filter.ContactFilter{
				ID:       contactstest.Ptr("e3a1e174-cdf8-4683-b112-e198144eb9d2"),
				FullName: contactstest.Ptr("John Doe"),
				Paymail:  contactstest.Ptr("john.doe@test.com"),
				Status:   contactstest.Ptr("confirmed"),
			},
			expectedParams: url.Values{
				"paymail":  []string{"john.doe@test.com"},
				"status":   []string{"confirmed"},
				"id":       []string{"e3a1e174-cdf8-4683-b112-e198144eb9d2"},
				"fullName": []string{"John Doe"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := contactFilterQueryBuilder{
				contactFilter:      tc.filter,
				modelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
