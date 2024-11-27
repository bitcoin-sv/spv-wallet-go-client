package querybuilders

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type PageFilterBuilder struct {
	Page filter.Page
}

func (p *PageFilterBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	params.AddPair("page", p.Page.Number)
	params.AddPair("size", p.Page.Size)
	params.AddPair("sort", p.Page.Sort)
	params.AddPair("sortBy", p.Page.SortBy)

	return params.Values, nil
}
