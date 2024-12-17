package transactions

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type adminTransactionFilterQueryBuilder struct {
	transactionFilter filter.AdminTransactionFilter
}

func (a *adminTransactionFilterQueryBuilder) Build() (url.Values, error) {
	builder := transactions.TransactionFilterBuilder{
		TransactionFilter:  a.transactionFilter.TransactionFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: a.transactionFilter.ModelFilter},
	}
	params, err := builder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubid", a.transactionFilter.XPubID) // xpubid should be replaced by xpubId in filter model.
	return params.Values, nil
}
