package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// PaymailAddressPage is an alias for the paymail addresses response page model returned by the SPV Wallet API.
// It provides a paginated list of paymails along with pagination metadata.
type PaymailAddressPage = response.PageModel[response.PaymailAddress]

// PaymailQueryFilter aggregates available paymail filters used in SPV Wallet API search request calls.
type PaymailQueryFilter interface {
	filter.AdminPaymailFilter | filter.PaymailFilter
}

// PaymailQuery aggregates query parameters for constructing the paymail addresses endpoint URL.
// It contains filters for metadata, pagination, and paymail-specific attributes.
type PaymailQuery[T PaymailQueryFilter] struct {
	Metadata      map[string]any // Metadata filters for refining the search.
	PageFilter    filter.Page    // Pagination details, including page number, size, and sorting.
	PaymailFilter T              // Filters for paymail attributes (e.g., ID, xPubID, alias, domain, public name, etc.).
}

// PaymailQueryOption defines a functional option for configuring a PaymailQuery instance.
// These options allow flexible setup of filters and pagination for the query.
type PaymailQueryOption[T PaymailQueryFilter] func(*PaymailQuery[T])

// PaymailQueryWithMetadataFilter adds metadata filters to the paymail addresses search URL.
// The provided metadata attributes are appended as query parameters.
func PaymailQueryWithMetadataFilter[T PaymailQueryFilter](m map[string]any) PaymailQueryOption[T] {
	return func(pq *PaymailQuery[T]) {
		pq.Metadata = m
	}
}

// PaymailQueryWithPageFilter adds pagination filters, such as page number, size, and sorting options,
// to the paymail addresses search URL as query parameters.
func PaymailQueryWithPageFilter[T PaymailQueryFilter](f filter.Page) PaymailQueryOption[T] {
	return func(pq *PaymailQuery[T]) {
		pq.PageFilter = f
	}
}

// PaymailQueryWithPaymailFilter adds filters based on the provided paymail filter type.
// These filters can include attributes like ID, xPubID, alias, domain, and public name.
// The filters are appended as query parameters to the paymail addresses search URL.
func PaymailQueryWithPaymailFilter[T PaymailQueryFilter](f T) PaymailQueryOption[T] {
	return func(pq *PaymailQuery[T]) {
		pq.PaymailFilter = f
	}
}
