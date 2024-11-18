package users

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type accessKeyFilterQueryBuilder struct {
	accessKeyFilter    filter.AccessKeyFilter
	modelFilterBuilder querybuilders.ModelFilterBuilder
}

func (a *accessKeyFilterQueryBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := a.modelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("revokedRange", a.accessKeyFilter.RevokedRange)
	return params.Values, nil
}
