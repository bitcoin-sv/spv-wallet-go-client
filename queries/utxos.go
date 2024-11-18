package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// UtxosPage is an alias for the UTXOs response page model returned by the SPV Wallet API.
// It contains a paginated list of UTXOs along with pagination metadata.
type UtxosPage = response.PageModel[response.Utxo]

// UtxoQuery aggregates query parameters for constructing the UTXOs search endpoint URL.
// It includes filters for metadata, pagination, and specific UTXO attributes.
type UtxoQuery struct {
	Metadata   map[string]any    // Filters based on metadata attributes.
	PageFilter filter.Page       // Pagination details, such as page number, size, and sort order.
	UtxoFilter filter.UtxoFilter // Filters for UTXO properties (e.g., date ranges, transaction ID, satoshis, status).
}

// UtxoQueryOption defines a functional option for configuring a UtxoQuery instance.
type UtxoQueryOption func(*UtxoQuery)

// UtxoQueryWithMetadataFilter applies metadata filters to the UTXOs search endpoint URL.
// The provided metadata attributes are added as query parameters.
func UtxoQueryWithMetadataFilter(m map[string]any) UtxoQueryOption {
	return func(uq *UtxoQuery) {
		uq.Metadata = m
	}
}

// UtxoQueryWithPageFilter sets pagination details for the UTXOs search endpoint URL.
// This includes page number, page size, and sort order, added as query parameters.
func UtxoQueryWithPageFilter(pf filter.Page) UtxoQueryOption {
	return func(uq *UtxoQuery) {
		uq.PageFilter = pf
	}
}

// UtxoQueryWithUtxoFilter applies filters for UTXO properties to the search endpoint URL.
// These include reserved date ranges (e.g., created, updated), transaction ID, output index,
// satoshis, etc., added as query parameters.
func UtxoQueryWithUtxoFilter(uf filter.UtxoFilter) UtxoQueryOption {
	return func(uq *UtxoQuery) {
		uq.UtxoFilter = uf
	}
}
