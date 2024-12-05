package accesskeys

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestAdminAccessKeyFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminAccessKeyFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"access key filter: zero values": {
			expectedParams: make(url.Values),
		},
		"access key filter: filter with only 'revoked range' field set": {
			filter: filter.AdminAccessKeyFilter{
				XpubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
			},
			expectedParams: url.Values{
				"xPubId": []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
			},
		},
		"access key filter: all fields set": {
			filter: filter.AdminAccessKeyFilter{
				XpubID: ptr("9b496655-616a-48cd-a3f8-89608473a5f1"),
				AccessKeyFilter: filter.AccessKeyFilter{
					RevokedRange: &filter.TimeRange{
						From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
						To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
					},
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
				"xPubId":             []string{"9b496655-616a-48cd-a3f8-89608473a5f1"},
				"revokedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"revokedRange[to]":   []string{"2021-02-02T00:00:00Z"},
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
			queryBuilder := adminAccessKeyFilterQueryBuilder{
				adminAccessKeyFilter: tc.filter,
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
