package contacts

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type adminContactFilterBuilder struct {
	contactFilter filter.AdminContactFilter
}

func (p *adminContactFilterBuilder) Build() (url.Values, error) {
	builder := contacts.ContactFilterQueryBuilder{
		ContactFilter:      p.contactFilter.ContactFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: p.contactFilter.ModelFilter},
	}
	params, err := builder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubId", p.contactFilter.XPubID)
	return params.Values, nil
}
