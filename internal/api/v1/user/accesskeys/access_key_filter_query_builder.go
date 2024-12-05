package accesskeys

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type AccessKeyFilterQueryBuilder struct {
	AccessKeyFilter    filter.AccessKeyFilter
	ModelFilterBuilder querybuilders.ModelFilterBuilder
}

func (a *AccessKeyFilterQueryBuilder) BuildExtendedURLValues() (*querybuilders.ExtendedURLValues, error) {
	modelFilterBuilder, err := a.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("revokedRange", a.AccessKeyFilter.RevokedRange)
	return params, nil
}

func (a *AccessKeyFilterQueryBuilder) Build() (url.Values, error) {
	params, err := a.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	return params.Values, nil
}
