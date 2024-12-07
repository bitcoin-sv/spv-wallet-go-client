package utxos

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type UtxoFilterQueryBuilder struct {
	UtxoFilter         filter.UtxoFilter
	ModelFilterBuilder querybuilders.ModelFilterBuilder
}

func (u *UtxoFilterQueryBuilder) BuildExtendedURLValues() (*querybuilders.ExtendedURLValues, error) {
	modelFilterBuilder, err := u.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("transactionId", u.UtxoFilter.TransactionID)
	params.AddPair("outputIndex", u.UtxoFilter.OutputIndex)
	params.AddPair("id", u.UtxoFilter.ID)
	params.AddPair("satoshis", u.UtxoFilter.Satoshis)
	params.AddPair("scriptPubKey", u.UtxoFilter.ScriptPubKey)
	params.AddPair("type", u.UtxoFilter.Type)
	params.AddPair("draftId", u.UtxoFilter.DraftID)
	params.AddPair("reservedRange", u.UtxoFilter.ReservedRange)
	params.AddPair("spendingTxId", u.UtxoFilter.SpendingTxID)
	return params, nil
}

func (u *UtxoFilterQueryBuilder) Build() (url.Values, error) {
	params, err := u.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	return params.Values, nil
}
