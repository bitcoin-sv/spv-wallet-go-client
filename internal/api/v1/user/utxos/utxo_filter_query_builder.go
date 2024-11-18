package utxos

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type utxoFilterQueryBuilder struct {
	utxoFilter         filter.UtxoFilter
	modelFilterBuilder querybuilders.ModelFilterBuilder
}

func (u *utxoFilterQueryBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := u.modelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("transactionId", u.utxoFilter.TransactionID)
	params.AddPair("outputIndex", u.utxoFilter.OutputIndex)
	params.AddPair("id", u.utxoFilter.ID)
	params.AddPair("satoshis", u.utxoFilter.Satoshis)
	params.AddPair("scriptPubKey", u.utxoFilter.ScriptPubKey)
	params.AddPair("type", u.utxoFilter.Type)
	params.AddPair("draftId", u.utxoFilter.DraftID)
	params.AddPair("reservedRange", u.utxoFilter.ReservedRange)
	params.AddPair("spendingTxId", u.utxoFilter.SpendingTxID)
	return params.Values, nil
}
