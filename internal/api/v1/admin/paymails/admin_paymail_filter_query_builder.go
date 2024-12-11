package paymails

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/paymails"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type adminPaymailFilterBuilder struct {
	paymailFilter filter.AdminPaymailFilter
}

func (p *adminPaymailFilterBuilder) Build() (url.Values, error) {
	paymailFilterBuilder := paymails.PaymailFilterBuilder{
		PaymailFilter:      p.paymailFilter.PaymailFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: p.paymailFilter.ModelFilter},
	}
	params, err := paymailFilterBuilder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubId", p.paymailFilter.XpubID)
	return params.Values, nil
}
