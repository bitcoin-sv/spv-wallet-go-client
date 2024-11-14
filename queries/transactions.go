package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// TransactionPage is an alias for the transactions response page model
// returned by the SPV Wallet API, which contains a paginated list of
// transactions along with pagination metadata.
type TransactionPage = response.PageModel[response.Transaction]

// TransactionsQuery aggregates query parameters for constructing a transactions endpoint URL.
// It holds filters for metadata, transaction-specific attributes, and pagination.
type TransactionsQuery struct {
	Metadata map[string]any           // Metadata filters for the transactions.
	Filter   filter.TransactionFilter // Transaction-specific filters (e.g., block height, status).
	Page     filter.Page              // Pagination details (page number, size, sorting).
}

// TransactionsQueryOption defines a functional option for configuring a TransactionsQuery instance.
type TransactionsQueryOption func(*TransactionsQuery)

// TransactionsQueryWithMetadataFilter adds metadata filters to the transactions endpoint URL.
// The specified metadata attributes will be appended as query parameters.
func TransactionsQueryWithMetadataFilter(m map[string]any) TransactionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.Metadata = m
	}
}

// TransactionsQueryWithFilter adds transaction-specific filters, such as block height, block hash,
// transaction status, etc., to the transactions endpoint URL as query parameters.
func TransactionsQueryWithFilter(tf filter.TransactionFilter) TransactionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.Filter = tf
	}
}

// TransactionsQueryWithPageFilter adds pagination details, like page number, page size, and sort order,
// to the transactions endpoint URL as query parameters.
func TransactionsQueryWithPageFilter(pf filter.Page) TransactionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.Page = pf
	}
}
