package utxos

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	userutxos "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type adminUtxoFilterQueryBuilder struct {
	utxoFilter filter.AdminUtxoFilter
}

func (a *adminUtxoFilterQueryBuilder) Build() (url.Values, error) {
	utxoFilterQueryBuilder := userutxos.UtxoFilterQueryBuilder{
		UtxoFilter:         a.utxoFilter.UtxoFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: a.utxoFilter.ModelFilter},
	}

	params, err := utxoFilterQueryBuilder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubId", a.utxoFilter.XpubID)
	return params.Values, nil
}
