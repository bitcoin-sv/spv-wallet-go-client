package merkleroots

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
)

type merkleRootsFilterQueryBuilder struct {
	query queries.MerkleRootsQuery
}

func (m *merkleRootsFilterQueryBuilder) Build() (url.Values, error) {
	params := querybuilders.NewExtendedURLValues()
	params.AddPair("batchSize", m.query.BatchSize)
	params.AddPair("lastEvaluatedKey", m.query.LastEvaluatedKey)
	return params.Values, nil
}
