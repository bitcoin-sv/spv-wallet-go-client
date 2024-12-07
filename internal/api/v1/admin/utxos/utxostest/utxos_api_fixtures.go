package utxostest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedUtxosPage(t *testing.T) *queries.UtxosPage {
	return &queries.UtxosPage{
		Content: []*response.Utxo{
			{
				Model:        response.Model{CreatedAt: spvwallettest.ParseTime(t, "2024-11-18T07:19:51.661656Z"), UpdatedAt: spvwallettest.ParseTime(t, "2024-11-18T07:19:51.663878Z")},
				UtxoPointer:  response.UtxoPointer{TransactionID: "ba371529-9746-4912-b9e4-4b3dc0539a40"},
				ID:           "21e69deb-8ea9-451b-9e1c-e89086ae439e",
				XpubID:       "2054a737-18d1-4e6c-9d3e-370a68ffe7f0",
				Satoshis:     1,
				ScriptPubKey: "2c2de44c-acf8-4507-9cdb-cf9cf5273253",
				Type:         "pubkeyhash",
				DraftID:      "",
				ReservedAt:   spvwallettest.ParseTime(t, "0001-01-01T00:00:00Z"),
				SpendingTxID: "",
				Transaction: &response.Transaction{
					Model: response.Model{
						CreatedAt: spvwallettest.ParseTime(t, "2024-11-18T07:19:51.661646Z"),
						UpdatedAt: spvwallettest.ParseTime(t, "2024-11-18T08:24:08.217141Z"),

						Metadata: map[string]any{
							"receiver": "john.doe.test@test.4chain.space",
							"sender":   "john.doe.test@test.4chain.space",
						},
					},
					ID:                   "d4fb1106-f023-43ce-9924-1d6f94bd5fbc",
					Hex:                  "cc75acd2-66b5-4970-9965-6cd621cd40cd",
					XpubInIDs:            []string{"9ecc1e7a-e122-41bd-9949-79ae8811de05"},
					XpubOutIDs:           []string{"227bcc1e-95cb-4113-9d56-d583857fdf86"},
					BlockHash:            "f3063295-7caf-4b11-8361-193b01a302c1",
					BlockHeight:          871343,
					Fee:                  1,
					NumberOfInputs:       2,
					NumberOfOutputs:      1,
					DraftID:              "6ea5d19a-f57f-42fb-94a4-287a627e76ce",
					TotalValue:           1,
					Status:               "MINED",
					TransactionDirection: "outgoing",
				},
			},
		},
		Page: response.PageDescription{
			Size:          1,
			Number:        1,
			TotalElements: 1,
			TotalPages:    1,
		},
	}
}
