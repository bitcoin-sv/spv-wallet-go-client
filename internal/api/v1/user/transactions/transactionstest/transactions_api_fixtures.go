package transactionstest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedDraftTransactionWithWrongInputs(t *testing.T) *response.DraftTransaction {
	draftWithWrongInputs := ExpectedDraftTransactionWithHex(t)
	draftWithWrongInputs.Configuration.Inputs[0].TransactionID = "wrong-input-transaction-id"
	return draftWithWrongInputs
}

func ExpectedDraftTransactionWithWrongLockingScript(t *testing.T) *response.DraftTransaction {
	draftWithWrongLockingScript := ExpectedDraftTransactionWithHex(t)
	draftWithWrongLockingScript.Configuration.Inputs[0].Destination.LockingScript = "wrong-locking-script"
	return draftWithWrongLockingScript
}

func ExpectedDraftTransactionWithWrongHex(t *testing.T) *response.DraftTransaction {
	draftWithWrongHex := ExpectedDraftTransactionWithHex(t)
	draftWithWrongHex.Hex = "wrong-hex"
	return draftWithWrongHex
}

func ExpectedDraftTransactionWithHex(t *testing.T) *response.DraftTransaction {
	return &response.DraftTransaction{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-12-02T12:04:33.855018Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-12-02T13:04:33.855036+01:00"),
		},
		ID:        "de3b8ef7041b2a528bc47ecdb3b87b06b61407fe24789bc02f9d49bfc234b4d5",
		Hex:       "01000000014c037d55e72d2ee6a95ff67bd758c4cee9c7545bb4d72ba77584152fcfa070120100000000ffffffff0200000000000000000e006a0568656c6c6f05776f726c6408000000000000001976a914702cef80a7039a1aebb70dc05ce1e439646fa33788ac00000000",
		XpubID:    "55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973",
		ExpiresAt: testutils.ParseTime(t, "2024-12-02T12:04:53.840989Z"),
		Configuration: response.TransactionConfig{
			ChangeSatoshis: 8,
			ChangeDestinations: []*response.Destination{
				{
					Model: response.Model{
						CreatedAt: testutils.ParseTime(t, "2024-12-02T12:04:33.853019Z"),
						UpdatedAt: testutils.ParseTime(t, "2024-12-02T13:04:33.853035+01:00"),
					},
					ID:            "872a51f9eed774e7e5051cec19db192783521b5a9e0d4d814d46bdce338a32dc",
					XpubID:        "55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973",
					LockingScript: "76a914702cef80a7039a1aebb70dc05ce1e439646fa33788ac",
					Type:          "pubkeyhash",
					Chain:         1,
					Num:           18,
					Address:       "1BE8WfQkDDYE3zEgxBdRNuAxsnHkDcuPdT",
					DraftID:       "de3b8ef7041b2a528bc47ecdb3b87b06b61407fe24789bc02f9d49bfc234b4d5",
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
							CreatedAt: testutils.ParseTime(t, "2024-11-29T23:13:54.0229Z"),
							UpdatedAt: testutils.ParseTime(t, "2024-12-02T12:04:33.847931Z"),
						},
						UtxoPointer: response.UtxoPointer{
							TransactionID: "1270a0cf2f158475a72bd7b45b54c7e9cec458d77bf65fa9e62e2de7557d034c",
						},
						ID:           "062bc7c22bf1f08e1716169e73f647e8ae2175ee3e8479de1bbc052eaa514d1d",
						XpubID:       "55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973",
						Satoshis:     9,
						ScriptPubKey: "76a9146637345046fd4d78a9ce187370db0ab7c15dd10488ac",
						Type:         "pubkeyhash",
						DraftID:      "de3b8ef7041b2a528bc47ecdb3b87b06b61407fe24789bc02f9d49bfc234b4d5",
						ReservedAt:   testutils.ParseTime(t, "2024-12-02T12:04:33.846479Z"),
					},
					Destination: response.Destination{
						Model: response.Model{
							CreatedAt: testutils.ParseTime(t, "2024-11-29T23:13:54.000014Z"),
							UpdatedAt: testutils.ParseTime(t, "2024-11-30T00:13:54.000029+01:00"),
						},
						ID:            "886d2ac60ac7fa630ad68954d6eb865314c484b4418e2469d48e4170dec7771f",
						XpubID:        "55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973",
						LockingScript: "76a9146637345046fd4d78a9ce187370db0ab7c15dd10488ac",
						Type:          "pubkeyhash",
						Chain:         1,
						Num:           16,
						Address:       "1AKU4EU46p38GWhaEcvuLL2UK23Fv14cwn",
						DraftID:       "254686b74d37878852b41503cf33604d5f6ba692705df08a855dc4d926b47251",
					},
				},
			},
			Outputs: []*response.TransactionOutput{
				{
					PaymailP4: nil,
					Satoshis:  0,
					Scripts: []*response.ScriptOutput{
						{
							Script:     "006a0568656c6c6f05776f726c64",
							ScriptType: "nulldata",
						},
					},
					UseForChange: false,
				},
				{
					Satoshis: 8,
					Scripts: []*response.ScriptOutput{
						{
							Address:    "1BE8WfQkDDYE3zEgxBdRNuAxsnHkDcuPdT",
							Satoshis:   8,
							Script:     "76a914702cef80a7039a1aebb70dc05ce1e439646fa33788ac",
							ScriptType: "pubkeyhash",
						},
					},
					To:           "1BE8WfQkDDYE3zEgxBdRNuAxsnHkDcuPdT",
					UseForChange: false,
				},
			},
		},
	}
}

func ExpectedSendToRecipientsTransaction(t *testing.T) *response.Transaction {
	return &response.Transaction{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-12-03T16:10:48.551774Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-12-03T16:10:49.080876Z"),
		},
		ID:                   "a4f86cdfefc3339bd3bd7861ad642feab05798f8a31cd67f81aec3c8c87083e0",
		Hex:                  "01000000017c8b38da58d766c75d94ca65f919723651090ed03e89a6cbc31ab7b87923d6ea010000006b483045022100b2e654169dda17a68c74b24d21e7b2e0dfef7fccdad9ab1e1ec87d2cab910e1d02206e247df32c4fe845af61001c8ec1dda718d737021df27c1d3f2e47d9fe76dd4241210265332864a94ed4c82bf3dacafbb828479b0a7fd0a73e62f60f6224dbf1504261ffffffff0200000000000000000e006a0568656c6c6f05776f726c6407000000000000001976a91464d00b8045c9e432b469f762b7e5beac2ef5a20c88ac00000000",
		XpubInIDs:            []string{"55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973"},
		XpubOutIDs:           []string{"55e5aeae101bf7dc49db2abfccfab9fb5f56a6b594fdcc87e5f5a94bfe94b973"},
		BlockHash:            "47758f612c6bf5b454bcd642fe8031f6",
		BlockHeight:          1024,
		Fee:                  1,
		NumberOfInputs:       1,
		NumberOfOutputs:      2,
		DraftID:              "4b0571c11a8a96d5af85bfbc32b98d4de6f3cc788afa2a6cd028ef9b69acc779",
		TotalValue:           0,
		OutputValue:          -1,
		Status:               "BROADCASTED",
		TransactionDirection: "outgoing",
	}
}

func ExpectedDraftTransaction(t *testing.T) *response.DraftTransaction {
	return &response.DraftTransaction{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
			Metadata: map[string]interface{}{
				"receiver": "john.doe.test4@john.doe.test.4chain.space",
				"sender":   "john.doe.test4@john.doe.test.4chain.space",
			},
		},
		ID:        "36be741b-31c7-4aed-8840-5e5b2eafeb41",
		Hex:       "c959fdb6-f438-4ef9-aef9-92a1852885ef",
		XpubID:    "3f0a90d3-4f8b-45f6-81e4-9858fa47ecc0",
		ExpiresAt: testutils.ParseTime(t, "2024-11-05T07:30:27.372912Z"),
		Configuration: response.TransactionConfig{
			ChangeSatoshis: 98,
			ChangeDestinations: []*response.Destination{
				{
					Model: response.Model{
						CreatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
						UpdatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
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
							CreatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
							UpdatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
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
						ReservedAt:   testutils.ParseTime(t, "2024-11-05T07:30:14.207287Z"),
					},
					Destination: response.Destination{
						Model: response.Model{
							CreatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
							UpdatedAt: testutils.ParseTime(t, "2024-11-05T07:30:14.219077Z"),
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
						PaymailExternalDerivationNum: testutils.Ptr(uint32(3)),
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
			CreatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
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
			CreatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
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
			CreatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
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

func ExpectedTransactionsPage(t *testing.T) *response.PageModel[response.Transaction] {
	return &response.PageModel[response.Transaction]{
		Content: []*response.Transaction{
			{
				Model: response.Model{
					CreatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
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
					CreatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-10-07T14:03:26.736816Z"),
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
		},
		Page: response.PageDescription{
			Size:          2,
			Number:        2,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}
