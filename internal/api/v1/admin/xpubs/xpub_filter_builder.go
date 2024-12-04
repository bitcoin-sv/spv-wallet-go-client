package xpubs

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type xpubFilterBuilder struct {
	xpubFilter         filter.XpubFilter
	modelFilterBuilder querybuilders.ModelFilterBuilder
}

func (x *xpubFilterBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := x.modelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", x.xpubFilter.ID)
	params.AddPair("currentBalance", x.xpubFilter.CurrentBalance)
	return params.Values, nil
}
