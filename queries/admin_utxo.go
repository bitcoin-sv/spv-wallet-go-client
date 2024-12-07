package queries

import "github.com/bitcoin-sv/spv-wallet/models/filter"

// AdminUtxoQuery aggregates query parameters for constructing the UTXOs search endpoint URL.
// It includes filters for metadata, pagination, and specific UTXO attributes.
type AdminUtxoQuery struct {
	Metadata   map[string]any         // Filters based on metadata attributes.
	PageFilter filter.Page            // Pagination details, such as page number, size, and sort order.
	UtxoFilter filter.AdminUtxoFilter // Filters for UTXO properties (e.g., date ranges, transaction ID, satoshis, status) and xPub for database queries.
}

// AdminUtxoQueryOption defines a functional option for configuring an AdminUtxoQuery instance.
type AdminUtxoQueryOption func(*AdminUtxoQuery)

// AdminUtxoQueryWithFilter applies filters for UTXO properties to the search endpoint URL.
// These include reserved date ranges (e.g., created, updated), transaction ID, output index,
// satoshis, etc., added as query parameters.
func AdminUtxoQueryWithFilter(f filter.AdminUtxoFilter) AdminUtxoQueryOption {
	return func(q *AdminUtxoQuery) {
		q.UtxoFilter = f
	}
}

// AdminUtxoQueryWithMetadataFilter adds metadata filters to the UTXOs search endpoint URL.
// The specified metadata attributes will be appended as query parameters.
func AdminUtxoQueryWithMetadataFilter(m map[string]any) AdminUtxoQueryOption {
	return func(q *AdminUtxoQuery) {
		q.Metadata = m
	}
}

// AdminUtxoQueryWithPageFilter adds pagination details, such as page number, page size, and sorting
// options, to the UTXOs search endpoint URL as query parameters.
func AdminUtxoQueryWithPageFilter(f filter.Page) AdminUtxoQueryOption {
	return func(q *AdminUtxoQuery) {
		q.PageFilter = f
	}
}
