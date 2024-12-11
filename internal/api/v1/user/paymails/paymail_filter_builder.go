package paymails

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type PaymailFilterBuilder struct {
	PaymailFilter      filter.PaymailFilter
	ModelFilterBuilder querybuilders.ModelFilterBuilder
}

func (p *PaymailFilterBuilder) BuildExtendedURLValues() (*querybuilders.ExtendedURLValues, error) {
	modelFilterBuilder, err := p.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", p.PaymailFilter.ID)
	params.AddPair("alias", p.PaymailFilter.Alias)
	params.AddPair("domain", p.PaymailFilter.Domain)
	params.AddPair("publicName", p.PaymailFilter.PublicName)
	return params, nil
}

func (p *PaymailFilterBuilder) Build() (url.Values, error) {
	params, err := p.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}
	return params.Values, nil
}
