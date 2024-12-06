package paymails

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type paymailFilterBuilder struct {
	paymailFilter      filter.AdminPaymailFilter
	modelFilterBuilder querybuilders.ModelFilterBuilder
}

func (p *paymailFilterBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := p.modelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", p.paymailFilter.ID)
	params.AddPair("xpubId", p.paymailFilter.XpubID)
	params.AddPair("alias", p.paymailFilter.Alias)
	params.AddPair("domain", p.paymailFilter.Domain)
	params.AddPair("publicName", p.paymailFilter.PublicName)
	return params.Values, nil
}
