package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// UserContactsPage is an alias for the user contacts response page model returned by the SPV Wallet API.
// It provides a paginated list of user contacts along with pagination metadata.
type UserContactsPage = response.PageModel[response.Contact]

// ContactQuery aggregates query parameters for constructing the user contacts endpoint URL.
// It contains filters for metadata, pagination, and user contact-specific attributes.
type ContactQuery struct {
	Metadata      map[string]any       // Metadata filters for refining the search.
	PageFilter    filter.Page          // Pagination details, including page number, size, and sorting.
	ContactFilter filter.ContactFilter // Filters for contact attributes (paymail, public key, ID, status).
}

// ContactQueryOption defines a functional option for configuring a ContactQuery instance.
// These options allow flexible setup of filters and pagination for the query.
type ContactQueryOption func(*ContactQuery)

// ContactQueryWithMetadataFilter adds metadata filters to the user contacts search URL.
// The provided metadata attributes are appended as query parameters.
func ContactQueryWithMetadataFilter(m map[string]any) ContactQueryOption {
	return func(cq *ContactQuery) {
		cq.Metadata = m
	}
}

// ContactQueryWithPageFilter adds pagination settings, like page number, size, and sorting options,
// to the user contacts search URL as query parameters.
func ContactQueryWithPageFilter(f filter.Page) ContactQueryOption {
	return func(cq *ContactQuery) {
		cq.PageFilter = f
	}
}

// ContactQueryWithContactFilter adds filters for user contact attributes, such as paymail, public key,
// contact ID, and status. These filters are appended as query parameters to the user contacts search URL.
func ContactQueryWithContactFilter(cf filter.ContactFilter) ContactQueryOption {
	return func(cq *ContactQuery) {
		cq.ContactFilter = cf
	}
}
