package accesskeys

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/accesskeys"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type adminAccessKeyFilterQueryBuilder struct {
	adminAccessKeyFilter filter.AdminAccessKeyFilter
}

func (a *adminAccessKeyFilterQueryBuilder) Build() (url.Values, error) {
	accessKeyFilterQueryBuilder := accesskeys.AccessKeyFilterQueryBuilder{
		AccessKeyFilter:    a.adminAccessKeyFilter.AccessKeyFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: a.adminAccessKeyFilter.ModelFilter},
	}
	params, err := accessKeyFilterQueryBuilder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubId", a.adminAccessKeyFilter.XpubID)
	return params.Values, nil
}
