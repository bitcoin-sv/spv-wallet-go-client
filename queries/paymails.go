package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// PaymailAddressPage is an alias for the paymail addresses response page model returned by the SPV Wallet API.
// It provides a paginated list of paymails along with pagination metadata.
type PaymailAddressPage = response.PageModel[response.PaymailAddress]

// PaymailQuery aggregates query parameters for constructing the paymail addresses endpoint URL.
// It contains filters for metadata, pagination, and user paymails-specific attributes.
type PaymailQuery struct {
	Metadata      map[string]any            // Metadata filters for refining the search.
	PageFilter    filter.Page               // Pagination details, including page number, size, and sorting.
	PaymailFilter filter.AdminPaymailFilter // Filters for paymail attributes (ID, xPubID, alias, domain, public name).
}

// PaymailQueryOption defines a functional option for configuring a ContactQuery instance.
// These options allow flexible setup of filters and pagination for the query.
type PaymailQueryOption func(*PaymailQuery)

// PaymailQueryWithMetadataFilter adds metadata filters to the paymail addresses search URL.
// The provided metadata attributes are appended as query parameters.
func PaymailQueryWithMetadataFilter(m map[string]any) PaymailQueryOption {
	return func(pq *PaymailQuery) {
		pq.Metadata = m
	}
}

// PaymailQueryWithPageFilter adds pagination filters, like page number, size, and sorting options,
// to the paymail addresses search URL as query parameters.
func PaymailQueryWithPageFilter(f filter.Page) PaymailQueryOption {
	return func(pq *PaymailQuery) {
		pq.PageFilter = f
	}
}

// PaymailQueryWithPaymailFilter adds filters for paymail address attributes like ID, xPubID, alias, domain, public name.
// These filters are appended as query parameters to the paymail addresses search URL.
func PaymailQueryWithPaymailFilter(pf filter.AdminPaymailFilter) PaymailQueryOption {
	return func(pq *PaymailQuery) {
		pq.PaymailFilter = pf
	}
}
