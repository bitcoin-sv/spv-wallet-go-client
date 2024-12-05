package accesskeys_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/accesskeys"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AccessKeyFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"access key filter: zero values": {
			expectedParams: make(url.Values),
		},
		"access key filter: filter with only 'revoked range' field set": {
			filter: filter.AccessKeyFilter{
				RevokedRange: &filter.TimeRange{
					From: spvwallettest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   spvwallettest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectedParams: url.Values{
				"revokedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"revokedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"access key filter: all fields set": {
			filter: filter.AccessKeyFilter{
				RevokedRange: &filter.TimeRange{
					From: spvwallettest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   spvwallettest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
				ModelFilter: filter.ModelFilter{
					IncludeDeleted: spvwallettest.Ptr(true),
					CreatedRange: &filter.TimeRange{
						From: spvwallettest.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
						To:   spvwallettest.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					UpdatedRange: &filter.TimeRange{
						From: spvwallettest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
						To:   spvwallettest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
			expectedParams: url.Values{
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
			queryBuilder := accesskeys.AccessKeyFilterQueryBuilder{
				AccessKeyFilter:    tc.filter,
				ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
