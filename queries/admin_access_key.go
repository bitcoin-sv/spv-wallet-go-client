package queries

import "github.com/bitcoin-sv/spv-wallet/models/filter"

// AdminAccessKeyQuery aggregates query parameters for constructing the access key search endpoint URL.
// It includes filters for metadata, pagination, and specific access key attributes, such as xPub.
type AdminAccessKeyQuery struct {
	AdminAccessKeyFilter filter.AdminAccessKeyFilter // Filters for access key properties (date ranges, deletion status) and xPub for database queries.
	Metadata             map[string]any              // Filters based on metadata attributes.
	PageFilter           filter.Page                 // Pagination details, including page number, size, and sorting options.
}

// AdminAccessKeyQueryOption defines a functional option for configuring an AdminAccessKeyQuery instance.
type AdminAccessKeyQueryOption func(*AdminAccessKeyQuery)

// AdminAccessKeyQueryWithAdminAccessKeyFilter adds filters for access key properties, such as date ranges
// (created, updated, revoked) and a flag for deletion status. These filters are appended as query
// parameters to the access key search endpoint URL.
func AdminAccessKeyQueryWithAdminAccessKeyFilter(f filter.AdminAccessKeyFilter) AdminAccessKeyQueryOption {
	return func(aakq *AdminAccessKeyQuery) {
		aakq.AdminAccessKeyFilter = f
	}
}

// AdminAccessKeyQueryWithMetadataFilter adds metadata filters to the search access key endpoint URL.
// The specified metadata attributes will be appended as query parameters.
func AdminAccessKeyQueryWithMetadataFilter(m map[string]any) AdminAccessKeyQueryOption {
	return func(aakq *AdminAccessKeyQuery) {
		aakq.Metadata = m
	}
}

// AdminAccessKeyQueryWithPageFilter adds pagination details, such as page number, page size, and sorting
// options, to the search access key endpoint URL as query parameters.
func AdminAccessKeyQueryWithPageFilter(f filter.Page) AdminAccessKeyQueryOption {
	return func(aakq *AdminAccessKeyQuery) {
		aakq.PageFilter = f
	}
}
