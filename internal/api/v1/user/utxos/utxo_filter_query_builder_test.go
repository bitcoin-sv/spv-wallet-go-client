package utxos

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos/utxostest"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestUtoxFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.UtxoFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"utxo filter: zero values": {
			expectedParams: make(url.Values),
		},
		"utxo filter: filter with only 'transaction id' field set": {
			expectedParams: url.Values{
				"transactionId": []string{"124c2237-9b54-46c4-bf53-3cec86f7e316"},
			},
			filter: filter.UtxoFilter{
				TransactionID: utxostest.Ptr("124c2237-9b54-46c4-bf53-3cec86f7e316"),
			},
		},
		"utxo filter: filter with only 'output index' field set": {
			expectedParams: url.Values{
				"outputIndex": []string{"32"},
			},
			filter: filter.UtxoFilter{
				OutputIndex: utxostest.Ptr(uint32(32)),
			},
		},
		"utxo filter: filter with only 'id' field set": {
			expectedParams: url.Values{
				"id": []string{"abb6a871-dd95-4f7a-8090-ca34cff63801"},
			},
			filter: filter.UtxoFilter{
				ID: utxostest.Ptr("abb6a871-dd95-4f7a-8090-ca34cff63801"),
			},
		},
		"utxo filter: filter with only 'satoshis' field set": {
			expectedParams: url.Values{
				"satoshis": []string{"64"},
			},
			filter: filter.UtxoFilter{
				Satoshis: utxostest.Ptr(uint64(64)),
			},
		},
		"utxo filter: filter with only 'script pub key' field set": {
			expectedParams: url.Values{
				"scriptPubKey": []string{"3adec124-32eb-46f1-94f2-4949a86dbe8d"},
			},
			filter: filter.UtxoFilter{
				ScriptPubKey: utxostest.Ptr("3adec124-32eb-46f1-94f2-4949a86dbe8d"),
			},
		},
		"utxo filter: filter with only 'type' field set": {
			expectedParams: url.Values{
				"type": []string{"0f65e842-decf-4725-8ad9-877634280e4f"},
			},
			filter: filter.UtxoFilter{
				Type: utxostest.Ptr("0f65e842-decf-4725-8ad9-877634280e4f"),
			},
		},
		"utxo filter: filter with only 'draft id' field set": {
			expectedParams: url.Values{
				"draftId": []string{"2453797c-4089-4078-8723-5ecb68e70bd7"},
			},
			filter: filter.UtxoFilter{
				DraftID: utxostest.Ptr("2453797c-4089-4078-8723-5ecb68e70bd7"),
			},
		},
		"utxo filter: filter with only reserved range 'from' field set": {
			expectedParams: url.Values{
				"reservedRange[from]": []string{"2021-02-01T00:00:00Z"},
			},
			filter: filter.UtxoFilter{
				ReservedRange: &filter.TimeRange{
					From: utxostest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"utxo filter: filter with only reserved range 'to' field set": {
			expectedParams: url.Values{
				"reservedRange[to]": []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.UtxoFilter{
				ReservedRange: &filter.TimeRange{
					To: utxostest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"utxo filter: filter with only reserved range field set": {
			expectedParams: url.Values{
				"reservedRange[to]":   []string{"2021-02-02T00:00:00Z"},
				"reservedRange[from]": []string{"2021-02-01T00:00:00Z"},
			},
			filter: filter.UtxoFilter{
				ReservedRange: &filter.TimeRange{
					To:   utxostest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
					From: utxostest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"utxo filter: filter with only 'spending tx id' field set": {
			expectedParams: url.Values{
				"spendingTxId": []string{"7539366c-beb2-4405-8597-025bf2dc7cbd"},
			},
			filter: filter.UtxoFilter{
				SpendingTxID: utxostest.Ptr("7539366c-beb2-4405-8597-025bf2dc7cbd"),
			},
		},
		"utxo filter: all fields set": {
			expectedParams: url.Values{
				"scriptPubKey":        []string{"3adec124-32eb-46f1-94f2-4949a86dbe8d"},
				"draftId":             []string{"2453797c-4089-4078-8723-5ecb68e70bd7"},
				"reservedRange[to]":   []string{"2021-02-02T00:00:00Z"},
				"reservedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"transactionId":       []string{"124c2237-9b54-46c4-bf53-3cec86f7e316"},
				"spendingTxId":        []string{"7539366c-beb2-4405-8597-025bf2dc7cbd"},
				"type":                []string{"0f65e842-decf-4725-8ad9-877634280e4f"},
				"satoshis":            []string{"64"},
				"id":                  []string{"abb6a871-dd95-4f7a-8090-ca34cff63801"},
				"outputIndex":         []string{"32"},
			},
			filter: filter.UtxoFilter{
				SpendingTxID:  utxostest.Ptr("7539366c-beb2-4405-8597-025bf2dc7cbd"),
				DraftID:       utxostest.Ptr("2453797c-4089-4078-8723-5ecb68e70bd7"),
				Type:          utxostest.Ptr("0f65e842-decf-4725-8ad9-877634280e4f"),
				ScriptPubKey:  utxostest.Ptr("3adec124-32eb-46f1-94f2-4949a86dbe8d"),
				ID:            utxostest.Ptr("abb6a871-dd95-4f7a-8090-ca34cff63801"),
				OutputIndex:   utxostest.Ptr(uint32(32)),
				Satoshis:      utxostest.Ptr(uint64(64)),
				TransactionID: utxostest.Ptr("124c2237-9b54-46c4-bf53-3cec86f7e316"),
				ReservedRange: &filter.TimeRange{
					To:   utxostest.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
					From: utxostest.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := utxoFilterQueryBuilder{
				utxoFilter: tc.filter,
				modelFilterBuilder: querybuilders.ModelFilterBuilder{
					ModelFilter: tc.filter.ModelFilter,
				},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
