package querybuilders

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type ModelFilterBuilder struct {
	ModelFilter filter.ModelFilter
}

func (m *ModelFilterBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	params.AddPair("includeDeleted", m.ModelFilter.IncludeDeleted)
	params.AddPair("createdRange", m.ModelFilter.CreatedRange)
	params.AddPair("updatedRange", m.ModelFilter.UpdatedRange)

	return params.Values, nil
}
