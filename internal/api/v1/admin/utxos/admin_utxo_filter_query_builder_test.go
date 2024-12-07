package utxos

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestAdminUtoxFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.AdminUtxoFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"admin utxo filter: zero values": {
			expectedParams: make(url.Values),
		},
		"admin utxo filter: filter with only 'xpub id' field set": {
			filter: filter.AdminUtxoFilter{
				XpubID: ptr("c4dab549-902e-42fe-97d2-dc556b716f9a"),
			},
			expectedParams: url.Values{
				"xpubId": []string{"c4dab549-902e-42fe-97d2-dc556b716f9a"},
			},
		},
		"admin utxo filter: all fields set": {
			filter: filter.AdminUtxoFilter{
				XpubID: ptr("c4dab549-902e-42fe-97d2-dc556b716f9a"),
				UtxoFilter: filter.UtxoFilter{
					SpendingTxID:  ptr("7539366c-beb2-4405-8597-025bf2dc7cbd"),
					DraftID:       ptr("2453797c-4089-4078-8723-5ecb68e70bd7"),
					Type:          ptr("0f65e842-decf-4725-8ad9-877634280e4f"),
					ScriptPubKey:  ptr("3adec124-32eb-46f1-94f2-4949a86dbe8d"),
					ID:            ptr("abb6a871-dd95-4f7a-8090-ca34cff63801"),
					OutputIndex:   ptr(uint32(32)),
					Satoshis:      ptr(uint64(64)),
					TransactionID: ptr("124c2237-9b54-46c4-bf53-3cec86f7e316"),
					ReservedRange: &filter.TimeRange{
						To:   ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
						From: ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
			expectedParams: url.Values{
				"xpubId":              []string{"c4dab549-902e-42fe-97d2-dc556b716f9a"},
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
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := adminUtxoFilterQueryBuilder{utxoFilter: tc.filter}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedParams, got)
		})
	}
}

func ptr[T any](value T) *T { return &value }
