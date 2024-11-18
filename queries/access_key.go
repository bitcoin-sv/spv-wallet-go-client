package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// AccessKeyPage is an alias for the access key response page model
// returned by the SPV Wallet API, which contains a paginated list of
// access keys along with pagination metadata.
type AccessKeyPage = response.PageModel[response.AccessKey]

// AccessKeyQuery aggregates query parameters for constructing the search access key endpoint URL.
// It holds filters for metadata, pagination, and specific access key attributes.
type AccessKeyQuery struct {
	Metadata        map[string]any         // Metadata filters for the search.
	PageFilter      filter.Page            // Pagination details (page number, size, sorting).
	AccessKeyFilter filter.AccessKeyFilter // Filters for access key properties (date ranges, deletion status).
}

// AccessKeyQueryOption defines a functional option for configuring an AccessKeyQuery instance.
type AccessKeyQueryOption func(*AccessKeyQuery)

// AccessKeyQueryWithMetadataFilter adds metadata filters to the search access key endpoint URL.
// The specified metadata attributes will be appended as query parameters.
func AccessKeyQueryWithMetadataFilter(m map[string]any) AccessKeyQueryOption {
	return func(akq *AccessKeyQuery) {
		akq.Metadata = m
	}
}

// AccessKeyQueryWithAccessKeyFilter adds filters for access key properties, such as date ranges
// (created, updated, revoked) and a flag indicating deletion status. These will be appended
// as query parameters to the search access key endpoint URL.
func AccessKeyQueryWithAccessKeyFilter(f filter.AccessKeyFilter) AccessKeyQueryOption {
	return func(akq *AccessKeyQuery) {
		akq.AccessKeyFilter = f
	}
}

// AccessKeyQueryWithPageFilter adds pagination details, such as page number, page size, and sorting
// options, to the search access key endpoint URL as query parameters.
func AccessKeyQueryWithPageFilter(f filter.Page) AccessKeyQueryOption {
	return func(akq *AccessKeyQuery) {
		akq.PageFilter = f
	}
}
