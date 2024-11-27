package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// XPubPage represents a paginated response model containing XPubs,
// as provided by the SPV Wallet API.
type XPubPage = response.PageModel[response.Xpub]

// XPubQuery defines the query parameters used to construct the XPubs search endpoint URL.
// It includes filters for metadata, pagination, and attributes specific to XPubs and UTXOs.
type XPubQuery struct {
	Metadata   map[string]any    // Filters based on key-value pairs of metadata.
	PageFilter filter.Page       // Pagination settings, including page number, size, and sort order.
	XpubFilter filter.XpubFilter // Filters for XPub properties, such as ID, balance, and date ranges.
}

// XPubQueryOption specifies a functional option for configuring an XPubQuery instance.
type XPubQueryOption func(*XPubQuery)

// XPubQueryWithMetadataFilter applies metadata filters to the XPubQuery instance.
// The specified key-value pairs will be added as query parameters to the search URL.
func XPubQueryWithMetadataFilter(m map[string]any) XPubQueryOption {
	return func(xq *XPubQuery) {
		xq.Metadata = m
	}
}

// XPubQueryWithPageFilter applies pagination settings to the XPubQuery instance.
// These include details like page number, page size, and sort order, which will be
// appended as query parameters to the search URL.
func XPubQueryWithPageFilter(f filter.Page) XPubQueryOption {
	return func(xq *XPubQuery) {
		xq.PageFilter = f
	}
}

// XPubQueryWithXPubFilter applies XPub-specific filters to the XPubQuery instance.
// This includes filters for attributes like ID, balance, or date ranges, which will
// be added as query parameters to the search URL.
func XPubQueryWithXPubFilter(f filter.XpubFilter) XPubQueryOption {
	return func(xq *XPubQuery) {
		xq.XpubFilter = f
	}
}
