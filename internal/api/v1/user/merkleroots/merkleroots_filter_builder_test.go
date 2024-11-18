package merkleroots

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/stretchr/testify/require"
)

func TestMerklerootsFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		query          queries.MerkleRootsQuery
		expectedParams url.Values
		expectedErr    error
	}{
		"merkle roots query: zero value": {
			expectedParams: make(url.Values),
		},
		"merkle roots query: query with 'batch size' set only": {
			query: queries.MerkleRootsQuery{
				BatchSize: 10,
			},
			expectedParams: url.Values{
				"batchSize": []string{"10"},
			},
		},
		"merkle roots query: query with 'last evaluated key' set only": {
			query: queries.MerkleRootsQuery{
				LastEvaluatedKey: "key",
			},
			expectedParams: url.Values{
				"lastEvaluatedKey": []string{"key"},
			},
		},
		"merkle roots query: all fields set": {
			query: queries.MerkleRootsQuery{
				BatchSize:        10,
				LastEvaluatedKey: "key",
			},
			expectedParams: url.Values{
				"batchSize":        []string{"10"},
				"lastEvaluatedKey": []string{"key"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := merkleRootsFilterQueryBuilder{
				query: tc.query,
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, got, tc.expectedParams)
		})
	}
}
