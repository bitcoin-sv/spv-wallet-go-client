package queries

// MerkleRootsQuery aggregates query parameter filters for constructing a URL
// to retrieve merkle roots. The parameters, such as BatchSize and LastEvaluatedKey,
// determine the behavior of the API request.
type MerkleRootsQuery struct {
	BatchSize        int
	LastEvaluatedKey string
}

// MerkleRootsQueryOption represents a functional option that can be used to customize
// the construction of a MerkleRootsQuery. These options provide a flexible way to configure
// the query by applying filters such as batch size and last evaluated key.
type MerkleRootsQueryOption func(*MerkleRootsQuery)

// MerkleRootsQueryWithBatchSize creates a MerkleRootsQueryOption to set the batch size for the query.
// This function applies a filter to the URL being constructed, specifying how many merkle roots
// should be fetched in a single API request. If the given batch size `n` is greater than 0, it
// sets the BatchSize field in the MerkleRootsQuery.
func MerkleRootsQueryWithBatchSize(n int) MerkleRootsQueryOption {
	return func(q *MerkleRootsQuery) {
		q.BatchSize = n
	}
}

// MerkleRootsQueryWithLastEvaluatedKey creates a MerkleRootsQueryOption to set the last evaluated key for the query.
// This function applies a filter to the URL being constructed, using the last processed merkle root
// in the client's database as a key for pagination. If the provided string `key` is non-empty, it sets
// the LastEvaluatedKey field in the MerkleRootsQuery.
func MerkleRootsQueryWithLastEvaluatedKey(key string) MerkleRootsQueryOption {
	return func(q *MerkleRootsQuery) {
		q.LastEvaluatedKey = key
	}
}
