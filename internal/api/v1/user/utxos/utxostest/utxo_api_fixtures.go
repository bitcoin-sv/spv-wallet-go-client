package utxostest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}

func ParseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}

func Ptr[T any](value T) *T {
	return &value
}

func ExpectedUtxosPage(t *testing.T) *queries.UtxosPage {
	return &queries.UtxosPage{
		Content: []*response.Utxo{
			{
				ID:           "db9bdd87-432d-44e6-b08f-9c0abd0d90ef",
				XpubID:       "0f8ff805-a282-48d6-be70-8b607deba5f1",
				Satoshis:     100,
				ScriptPubKey: "88ca49f2-816e-4a0b-b5c5-e5c574e2d292",
				Type:         "pubkeyhash",
				ReservedAt:   ParseTime(t, "0001-01-01T00:00:00Z"),
				UtxoPointer: response.UtxoPointer{
					TransactionID: "f365f697-3db9-44fd-bd0d-ba8e94ca63f2",
					OutputIndex:   0,
				},
				Model: response.Model{
					CreatedAt: ParseTime(t, "2024-11-12T11:31:07.728974Z"),
					UpdatedAt: ParseTime(t, "2024-11-12T11:31:07.732139Z"),
				},
				Transaction: &response.Transaction{
					Model: response.Model{
						CreatedAt: ParseTime(t, "2024-11-12T11:31:07.72894Z"),
						UpdatedAt: ParseTime(t, "2024-11-12T12:33:35.266758Z"),

						Metadata: map[string]any{
							"domain":     "john.doe.test.space",
							"ip_address": "127.0.0.1",
							"p2p_tx_metadata": map[string]any{
								"pubkey": "d90c6998-010a-466f-83d7-25c39188a1c5",
								"sender": "john.doe@test.com",
							},
							"paymail_request": "HandleReceivedP2pTransaction",
							"reference_id":    "81fbfb26-e648-463e-99ce-ade498774c8f",
							"user_agent":      "node-fetch",
						},
					},
					ID:                   "ec943c46-bfa8-4764-820a-a604c8b6c890",
					Hex:                  "d825088b-1f04-406a-b046-059bc0736b11",
					XpubOutIDs:           []string{"6e980e21-a8f8-4699-9d11-98aef96bdf98"},
					BlockHash:            "a7755931-eceb-473e-ab5b-6a6459948166",
					BlockHeight:          1024,
					NumberOfInputs:       2,
					NumberOfOutputs:      3,
					TotalValue:           1305,
					Status:               "MINED",
					TransactionDirection: "outgoing",
				},
			},
			{
				ID:           "7ed4a935-6b62-4e83-9d97-7e9a7f9eab30",
				XpubID:       "68019cf3-616c-4d14-b9bf-cd9486b63f4f",
				Satoshis:     18,
				ScriptPubKey: "5e63148d-f506-43fb-88c3-2d98491625da",
				Type:         "pubkeyhash",
				ReservedAt:   ParseTime(t, "0001-01-01T00:00:00Z"),
				UtxoPointer: response.UtxoPointer{
					TransactionID: "54ed5bcb-a964-47af-892b-1054065c28a8",
					OutputIndex:   1,
				},
				Model: response.Model{
					CreatedAt: ParseTime(t, "2024-11-08T13:40:55.592Z"),
					UpdatedAt: ParseTime(t, "2024-11-08T13:40:55.593441Z"),
				},
				Transaction: &response.Transaction{
					Model: response.Model{
						CreatedAt: ParseTime(t, "2024-11-08T13:40:55.591986Z"),
						UpdatedAt: ParseTime(t, "2024-11-08T14:43:56.256571Z"),
					},
					ID:                   "29b89717-f139-45ae-9848-f2d7415ea596",
					Hex:                  "6a1c1ddb-f3c1-4491-98b4-9ce3eb016e60",
					XpubInIDs:            []string{"32dfa8c9-82e3-4f49-8d33-ff7130e1cfae"},
					XpubOutIDs:           []string{"b0559e5f-b4b5-416f-b1f8-116f19a89f30"},
					BlockHash:            "f90fb747-4cec-4e00-912a-582d46090d61",
					BlockHeight:          2048,
					Fee:                  1,
					NumberOfInputs:       2,
					NumberOfOutputs:      2,
					DraftID:              "057a743c-4c97-444b-b6ac-8b4a757aee8c",
					TotalValue:           0,
					Status:               "MINED",
					TransactionDirection: "outgoing",
				},
			},
		},
		Page: response.PageDescription{
			Size:          2,
			Number:        1,
			TotalElements: 9,
			TotalPages:    5,
		},
	}
}
