package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MerkleRootPage is an alias for the Merkle roots response page model
// returned by the SPV Wallet API. It provides a paginated list of Merkle roots
// along with pagination metadata.
type MerkleRootPage = models.MerkleRootsBHSResponse

// MerkleRootsQuery aggregates query parameters for constructing a URL to retrieve Merkle roots.
// These parameters, such as BatchSize and LastEvaluatedKey, control how the API request is processed.
type MerkleRootsQuery struct {
	BatchSize        int    // The number of Merkle roots to fetch in a single API request.
	LastEvaluatedKey string // A key used for pagination, indicating where to continue the query.
}

// MerkleRootsQueryOption defines a functional option for customizing a MerkleRootsQuery.
// These options allow for flexible configuration by applying filters like batch size or
// the last evaluated key for pagination.
type MerkleRootsQueryOption func(*MerkleRootsQuery)

// MerkleRootsQueryWithBatchSize returns a MerkleRootsQueryOption to set the batch size for the query.
// This option specifies how many Merkle roots should be retrieved in a single API request.
func MerkleRootsQueryWithBatchSize(n int) MerkleRootsQueryOption {
	return func(q *MerkleRootsQuery) {
		q.BatchSize = n
	}
}

// MerkleRootsQueryWithLastEvaluatedKey returns a MerkleRootsQueryOption to set the last evaluated key for pagination.
// This option uses the last processed Merkle root in the client's database to continue the query.
func MerkleRootsQueryWithLastEvaluatedKey(key string) MerkleRootsQueryOption {
	return func(q *MerkleRootsQuery) {
		q.LastEvaluatedKey = key
	}
}
