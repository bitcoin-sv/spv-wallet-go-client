package querybuilders

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type QueryParamsFilterBuilder struct {
	QueryParamsFilter filter.QueryParams
}

func (q *QueryParamsFilterBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	params.AddPair("page", q.QueryParamsFilter.Page)
	params.AddPair("size", q.QueryParamsFilter.PageSize)
	params.AddPair("sortBy", q.QueryParamsFilter.OrderByField)
	params.AddPair("sort", q.QueryParamsFilter.SortDirection)
	return params.Values, nil
}
