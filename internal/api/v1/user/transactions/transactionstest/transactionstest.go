package transactionstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedDraftTransaction(t *testing.T) *response.DraftTransaction {
	return &response.DraftTransaction{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
			UpdatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
			Metadata: map[string]interface{}{
				"receiver": "john.doe.test4@john.doe.test.4chain.space",
				"sender":   "john.doe.test4@john.doe.test.4chain.space",
			},
		},
		ID:        "36be741b-31c7-4aed-8840-5e5b2eafeb41",
		Hex:       "c959fdb6-f438-4ef9-aef9-92a1852885ef",
		XpubID:    "3f0a90d3-4f8b-45f6-81e4-9858fa47ecc0",
		ExpiresAt: ParseTime(t, "2024-11-05T07:30:27.372912Z"),
		Configuration: response.TransactionConfig{
			ChangeSatoshis: 98,
			ChangeDestinations: []*response.Destination{
				{
					Model: response.Model{
						CreatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
						UpdatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
					},
					ID:            "c86dd8f4-316f-4d71-be00-7bd1a38079e4",
					XpubID:        "d6884260-1624-415b-8625-652a59345ead",
					LockingScript: "189593db-0048-4fb7-80da-b69bce8fbf78",
					Type:          "pubkeyhash",
					Chain:         1,
					Num:           5,
					Address:       "3f96ea59-ac83-476e-a0ea-f0d668086081",
					DraftID:       "fc60742e-92b5-4a98-90a7-422d89879494",
				},
			},
			FeeUnit: &response.FeeUnit{
				Satoshis: 1,
				Bytes:    1000,
			},
			Inputs: []*response.TransactionInput{
				{
					Utxo: response.Utxo{
						Model: response.Model{
							CreatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
							UpdatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
						},
						UtxoPointer: response.UtxoPointer{
							TransactionID: "3e0c5f6d-0dfc-462d-8a63-31b7a20d0c6b",
						},
						ID:           "203277ff-006a-4e48-bbe9-2f1b6fb9ddfd",
						XpubID:       "4676a7d6-45f8-46b3-850b-68a9bb7642bc",
						Satoshis:     100,
						ScriptPubKey: "9d7eede4-00cd-47fd-ab3d-b0ae6d2ca6a6",
						Type:         "pubkeyhash",
						DraftID:      "f1ebe294-d921-4fb7-8b22-ed33e090e7ea",
						ReservedAt:   ParseTime(t, "2024-11-05T07:30:14.207287Z"),
					},
					Destination: response.Destination{
						Model: response.Model{
							CreatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
							UpdatedAt: ParseTime(t, "2024-11-05T07:30:14.219077Z"),
							Metadata: map[string]interface{}{
								"domain":          "john.doe.test.4chain.space",
								"ip_address":      "127.0.0.1",
								"paymail_request": "CreateP2PDestinationResponse",
								"reference_id":    "1a461311db24115cd5e0525f8c9b5613",
								"satoshis":        float64(100),
								"user_agent":      "node-fetch",
							},
						},
						ID:                           "bc22a0b9-d91c-4d0b-a7e4-8ea2d37e42db",
						XpubID:                       "325b1440-3af4-4a65-bf90-d88ed978948b",
						LockingScript:                "e459d941-d820-4663-a5d8-6a12457825e9",
						Type:                         "pubkeyhash",
						Chain:                        0,
						Num:                          0,
						PaymailExternalDerivationNum: Ptr(uint32(3)),
						Address:                      "6e4f50b1-356b-4453-a83e-2f412f328c25",
						DraftID:                      "",
					},
				},
			},
			Outputs: []*response.TransactionOutput{
				{
					PaymailP4: &response.PaymailP4{
						Alias:           "john.doe.test4",
						Domain:          "john.doe.test.4chain.space",
						FromPaymail:     "from@domain.com",
						ReceiveEndpoint: "https://john.doe.test.4chain.space:443/v1/bsvalias/beef/{alias}@{domain.tld}",
						ReferenceID:     "bdac6a12ec7f31feb5ae426e28c9ddfa",
						ResolutionType:  "p2p",
					},
					Satoshis: 1,
					Scripts: []*response.ScriptOutput{
						{
							Address:    "18p1xtQQeaVVpsxrSiRUhUKMyR5jPEvAhY",
							Satoshis:   1,
							Script:     "45a858f8-c645-48c3-bff0-f776d8d8452d",
							ScriptType: "pubkeyhash",
						},
					},
					To:           "john.doe.test4@john.doe.test.4chain.space",
					UseForChange: false,
				},
				{
					Satoshis: 98,
					Scripts: []*response.ScriptOutput{
						{
							Address:    "19a5857d-3eb9-43f8-b240-c29c05909fdc",
							Satoshis:   98,
							Script:     "cca457ab-2277-457b-bf53-17face515f5c",
							ScriptType: "pubkeyhash",
						},
					},
					To:           "b1e97d9c-e1e5-4120-b0f1-0363693b1959",
					UseForChange: false,
				},
			},
		},
	}
}

func ExpectedRecordTransaction(t *testing.T) *response.Transaction {
	return &response.Transaction{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			Metadata: map[string]interface{}{
				"key":  "value",
				"key2": "value2",
			},
		},
		ID:              "fdad0324-1185-4a54-8eae-f0c8858fa3ce",
		Hex:             "fda8f356-615e-4b4c-a3c8-53a47531a446",
		XpubInIDs:       []string{"e2be970c-a867-4e65-b141-7f2aafd44a42"},
		XpubOutIDs:      []string{"475e5e90-a117-46b6-b9e5-6983f2721b19"},
		BlockHash:       "47758f612c6bf5b454bcd642fe8031f6",
		BlockHeight:     1024,
		Fee:             1,
		NumberOfInputs:  3,
		NumberOfOutputs: 2,
		DraftID:         "d3fb66d6-6e3b-4a1f-aa80-dda848079663",
		TotalValue:      51,
		OutputValue:     50,
		Outputs: map[string]int64{
			"92640954841510a9d95f7737a43075f22ebf7255976549de4c52e8f3faf57470": -51,
			"9d07977d2fc14402426288a6010b4cdf7d91b61461acfb75af050b209d2d07ba": 50,
		},
		Status:               "MINED",
		TransactionDirection: "outgoing",
	}
}

func ExpectedTransactionWithUpdatedMetadata(t *testing.T) *response.Transaction {
	return &response.Transaction{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			Metadata: map[string]any{
				"domain":          "john.doe.test.4chain.space",
				"example_key1":    "example_key10_val",
				"example_key2":    "example_key20_val",
				"ip_address":      "127.0.0.01",
				"user_agent":      "node-fetch",
				"paymail_request": "HandleReceivedP2pTransaction",
				"reference_id":    "1c2dcc61-f48f-44f2-aba2-9a759a514d49",
				"p2p_tx_metadata": map[string]any{
					"pubkey": "3fa7af5b-4568-4873-86da-0aa442ca91dd",
					"sender": "john.doe@handcash.io",
				},
			},
		},
		ID:                   "2c250e21-c33a-41e3-a4e3-77c68b03244e",
		Hex:                  "283b1c6deb6d6263b3cec7a4701d46d3",
		XpubOutIDs:           []string{"4c9a0a0d-ea4f-4f03-b740-84438b3d210d"},
		BlockHash:            "47758f612c6bf5b454bcd642fe8031f6",
		BlockHeight:          512,
		Fee:                  1,
		NumberOfInputs:       2,
		NumberOfOutputs:      3,
		TotalValue:           311,
		OutputValue:          100,
		Status:               "MINED",
		TransactionDirection: "incoming",
	}
}

func ExpectedTransaction(t *testing.T) *response.Transaction {
	return &response.Transaction{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			Metadata: map[string]any{
				"domain":          "john.doe.test.4chain.space",
				"example_key1":    "example_key10_val",
				"ip_address":      "127.0.0.01",
				"user_agent":      "node-fetch",
				"paymail_request": "HandleReceivedP2pTransaction",
				"reference_id":    "1c2dcc61-f48f-44f2-aba2-9a759a514d49",
				"p2p_tx_metadata": map[string]any{
					"pubkey": "3fa7af5b-4568-4873-86da-0aa442ca91dd",
					"sender": "john.doe@handcash.io",
				},
			},
		},
		ID:                   "2c250e21-c33a-41e3-a4e3-77c68b03244e",
		Hex:                  "283b1c6deb6d6263b3cec7a4701d46d3",
		XpubOutIDs:           []string{"4c9a0a0d-ea4f-4f03-b740-84438b3d210d"},
		BlockHash:            "47758f612c6bf5b454bcd642fe8031f6",
		BlockHeight:          512,
		Fee:                  1,
		NumberOfInputs:       2,
		NumberOfOutputs:      3,
		TotalValue:           311,
		OutputValue:          100,
		Status:               "MINED",
		TransactionDirection: "incoming",
	}
}

func ExpectedTransactions(t *testing.T) []*response.Transaction {
	return []*response.Transaction{
		{
			Model: response.Model{
				CreatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
				UpdatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
				Metadata: map[string]any{
					"domain":          "john.doe.test.4chain.space",
					"example_key1":    "example_key10_val",
					"ip_address":      "127.0.0.01",
					"user_agent":      "node-fetch",
					"paymail_request": "HandleReceivedP2pTransaction",
					"reference_id":    "1c2dcc61-f48f-44f2-aba2-9a759a514d49",
					"p2p_tx_metadata": map[string]any{
						"pubkey": "3efe9fcb-859c-47f1-b85f-0fa8b1eee065",
						"sender": "john.doe@handcash.io",
					},
				},
			},
			ID:                   "2c250e21-c33a-41e3-a4e3-77c68b03244e",
			Hex:                  "283b1c6deb6d6263b3cec7a4701d46d3",
			XpubOutIDs:           []string{"4c9a0a0d-ea4f-4f03-b740-84438b3d210d"},
			BlockHash:            "47758f612c6bf5b454bcd642fe8031f6",
			BlockHeight:          512,
			Fee:                  1,
			NumberOfInputs:       2,
			NumberOfOutputs:      3,
			TotalValue:           311,
			OutputValue:          100,
			Status:               "MINED",
			TransactionDirection: "incoming",
		},
		{
			Model: response.Model{
				CreatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
				UpdatedAt: ParseTime(t, "2024-10-07T14:03:26.736816Z"),
				Metadata: map[string]any{
					"domain":          "jane.doe.test.4chain.space",
					"example_key101":  "example_key101_val",
					"ip_address":      "127.0.0.01",
					"user_agent":      "node-fetch",
					"paymail_request": "HandleReceivedP2pTransaction",
					"reference_id":    "2c6dcc71-f42f-54f2-ada1-1c658a515d50",
					"p2p_tx_metadata": map[string]any{
						"pubkey": "4fa8af6b-3217-2373-76da-0aa552ca88aa",
						"sender": "jane.doe@handcash.io",
					},
				},
			},
			ID:                   "1c110e11-c23a-51e5-a7e7-99c12b01233e",
			Hex:                  "283b1c7deb7d7773b3cec7a8801d47d2",
			XpubOutIDs:           []string{"2c8a1a1d-ea5f-5f04-b890-92418b2d411d"},
			BlockHash:            "56659f622c6bf5b554bcd742fe8132f9",
			BlockHeight:          1024,
			Fee:                  1,
			NumberOfInputs:       2,
			NumberOfOutputs:      3,
			TotalValue:           500,
			OutputValue:          200,
			Status:               "MINED",
			TransactionDirection: "incoming",
		},
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

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
