package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// ContactsPage is an alias for the user contacts response page model returned by the SPV Wallet API.
// It provides a paginated list of user contacts along with pagination metadata.
type ContactsPage = response.PageModel[response.Contact]

// AccessKeyPage is an alias for the access key response page model
// returned by the SPV Wallet API, which contains a paginated list of
// access keys along with pagination metadata.
type AccessKeyPage = response.PageModel[response.AccessKey]

// PaymailsPage is an alias for the paymail addresses response page model returned by the SPV Wallet API.
// It provides a paginated list of paymails along with pagination metadata.
type PaymailsPage = response.PageModel[response.PaymailAddress]

// TransactionPage is an alias for the transactions response page model
// returned by the SPV Wallet API, which contains a paginated list of
// transactions along with pagination metadata.
type TransactionPage = response.PageModel[response.Transaction]

// UtxosPage is an alias for the UTXOs response page model returned by the SPV Wallet API.
// It contains a paginated list of UTXOs along with pagination metadata.
type UtxosPage = response.PageModel[response.Utxo]

// XPubPage represents a paginated response model containing XPubs,
// as provided by the SPV Wallet API.
type XPubPage = response.PageModel[response.Xpub]

// QueryOption defines a functional option for configuring a generic query instance.
// These options allow flexible setup of filters, metadata, and pagination for the query.
type QueryOption[F QueryFilters] func(*Query[F])

// QueryWithMetadataFilter adds metadata filters to the search parameters URL.
// The provided metadata attributes are appended as query parameters.
func QueryWithMetadataFilter[F QueryFilters](m map[string]any) QueryOption[F] {
	return func(q *Query[F]) {
		q.Metadata = m
	}
}

// QueryWithPageFilter adds pagination filters, such as page number, size, and sorting options,
// to the search URL as query parameters.
func QueryWithPageFilter[F QueryFilters](f filter.Page) QueryOption[F] {
	return func(q *Query[F]) {
		q.PageFilter = f
	}
}

// QueryWithFilter adds search parameters to the search URL corresponding to the specified filter type.
func QueryWithFilter[F QueryFilters](f F) QueryOption[F] {
	return func(q *Query[F]) {
		q.Filter = f
	}
}

// AdminQueryFilters aggregates the supported query filter types used for constructing query parameters
// for SPV Wallet API admin search endpoints.
type AdminQueryFilters interface {
	filter.AdminAccessKeyFilter | filter.AdminUtxoFilter | filter.AdminPaymailFilter | filter.XpubFilter | filter.AdminTransactionFilter | filter.AdminContactFilter
}

// NonAdminQueryFilters aggregates the supported query filter types used for constructing query parameters
// for SPV Wallet API non-admin search endpoints.
type NonAdminQueryFilters interface {
	filter.AccessKeyFilter | filter.ContactFilter | filter.PaymailFilter | filter.TransactionFilter | filter.UtxoFilter
}

// QueryFilters aggregates filter types for both admin and non-admin query types.
type QueryFilters interface {
	AdminQueryFilters | NonAdminQueryFilters
}

// Query represents a generic query structure that aggregates metadata, pagination, and specific filters.
type Query[F QueryFilters] struct {
	Metadata   map[string]any // Metadata filters for refining the search.
	PageFilter filter.Page    // Pagination details, including page number, size, and sorting.
	Filter     F              // Specific filter for refining the query.
}

// NewQuery creates a new Query instance, applying the provided functional options.
// It allows flexible configuration of metadata, filters, and pagination for the query.
func NewQuery[F QueryFilters](opts ...QueryOption[F]) *Query[F] {
	var q Query[F]
	for _, o := range opts {
		o(&q)
	}
	return &q
}

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
