package transactionstest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedTransaction(t *testing.T) *response.Transaction {
	return &response.Transaction{
		Model: response.Model{
			CreatedAt: spvwallettest.ParseTime(t, "2024-10-02T10:34:57.931744Z"),
			UpdatedAt: spvwallettest.ParseTime(t, "2024-10-18T13:59:02.237607Z"),
			Metadata: map[string]any{
				"domain":     "john.doe.test.4chain.space",
				"ip_address": "127.0.0.1",
				"p2p_tx_metadata": map[string]any{
					"pubkey": "c110ad13-9ded-4df3-a7af-99215e80a609",
					"sender": "john.doe@handcash.io",
				},
				"paymail_request": "HandleReceivedP2pTransaction",
				"reference_id":    "802f85de-23fa-40e3-adc2-165f33b9c853",
				"user_agent":      "node-fetch",
			},
		},
		ID:                   "7efb8617-deb9-43cf-90df-d78782d40ab2",
		Hex:                  "7198c88c-b695-4a19-b4af-6d912f87bf29",
		XpubOutIDs:           []string{"df5b2ccc-04cc-4ae3-8e8b-d311160eba75"},
		BlockHash:            "d97337d6-d735-4b41-a118-70a8813e616d",
		BlockHeight:          864633,
		NumberOfInputs:       2,
		NumberOfOutputs:      3,
		TotalValue:           631,
		Status:               "MINED",
		TransactionDirection: "outgoing",
	}
}

func ExpectedTransactionsPage(t *testing.T) *response.PageModel[response.Transaction] {
	return &response.PageModel[response.Transaction]{
		Content: []*response.Transaction{
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-11-18T07:19:51.661646Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-11-18T08:24:08.217141Z"),
					Metadata: map[string]any{
						"receiver": "john.doe.test@john.test.4chain.space",
						"sender":   "jane.doe.test@jane.test.4chain.space",
					},
				},
				ID:              "8b17fb99-bb11-4cdd-b04c-2d35c3d5070f",
				Hex:             "ae3d28fe-2610-496d-95ea-fbfa50dd571d",
				XpubOutIDs:      []string{"f3d9355f-a152-4075-be05-52daa98cc87c"},
				XpubInIDs:       []string{"70dbd5a2-ab71-4a62-8240-ca5960d1327b"},
				BlockHash:       "8bcc3e0b-3a0f-44a2-a32e-7fa40fe50db5",
				BlockHeight:     871343,
				Fee:             1,
				NumberOfInputs:  2,
				NumberOfOutputs: 1,
				DraftID:         "c78a5ced-d3b0-4acf-bfff-7f4268760b0f",
				TotalValue:      1,
				Outputs: map[string]int64{
					"ba121a0d-af03-41cb-bfb2-592005f73e55": 1,
					"79c17f30-0580-48e1-8491-f147c869b73b": -2,
				},
			},
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-11-18T07:16:04.821925Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-11-19T13:09:39.501356Z"),
					Metadata: map[string]any{
						"domain":     "john.doe.test.4chain.space",
						"ip_address": "127.0.0.1",
						"p2p_tx_metadata": map[string]any{
							"note":   "example note",
							"pubkey": "bd568915-e532-4466-a0ab-70c7dead4b4b",
							"sender": "jane.doe@handcash.io",
						},
						"paymail_request": "HandleReceivedP2pTransaction",
						"reference_id":    "c7056aa68fb3586c74d0f8f7cb0ae52b",
						"user_agent":      "node-fetch",
					},
				},
				ID:              "2f130e50-3d0b-46b1-a1e6-ad866d60c2d2",
				Hex:             "8741462f-50ac-4fdc-bb22-b1db155899dc",
				XpubOutIDs:      []string{"1f1724bb-b167-4ac2-b97b-7a8e03a111c8"},
				BlockHash:       "fca9446e-5c56-44eb-930f-275ce0594f24",
				BlockHeight:     871343,
				Fee:             0,
				NumberOfInputs:  2,
				NumberOfOutputs: 3,
				TotalValue:      513,
				Outputs: map[string]int64{
					"e95ab632-2a96-466a-9676-5e86bc8a8d8d": 100,
				},
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
