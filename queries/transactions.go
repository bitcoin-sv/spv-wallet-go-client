package queries

import (
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// TransactionsQuery aggregates query parameter filters describing the key-value pairs
// that should be appended to the transaction URL being constructed
type TransactionsQuery struct {
	Metadata    map[string]any
	Filter      filter.TransactionFilter
	QueryParams filter.QueryParams
}

// TransctionsQueryOption represents a functional option for creating a customized list
// of transaction query parameters.
type TransctionsQueryOption func(*TransactionsQuery)

// TransactionsQueryWithMetadataFilter applies specific metadata attributes as filters
// to the transactions endpoint URL being constructed.
func TransactionsQueryWithMetadataFilter(m map[string]any) TransctionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.Metadata = m
	}
}

// TransactionsQueryWithFilter applies general query parameters like BlockHeight, BlockHash,
// transaction status, etc. to the transactions endpoint URL being constructed.
func TransactionsQueryWithFilter(tf filter.TransactionFilter) TransctionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.Filter = tf
	}
}

// TransactionsQueryWithQueryParamsFilter applies general query parameters like pagination and sort order etc.
// to the transactions endpoint URL being constructed.
func TransactionsQueryWithQueryParamsFilter(q filter.QueryParams) TransctionsQueryOption {
	return func(tq *TransactionsQuery) {
		tq.QueryParams = q
	}
}
