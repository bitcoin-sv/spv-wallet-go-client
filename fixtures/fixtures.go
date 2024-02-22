package fixtures

import (
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/common"
)

const (
	RequestType     = "http"
	ServerURL       = "https://example.com/"
	XPubString      = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	XPrivString     = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	AccessKeyString = "7779d24ca6f8821f225042bf55e8f80aa41b08b879b72827f51e41e6523b9cd0"
	PaymailAddress  = "address@paymail.com"
)

func MarshallForTestHandler(object any) string {
	json, err := json.Marshal(object)
	if err != nil {
		// as this is just for tests, empty string will make the tests fail,
		// so it's acceptable as an "error" here, in case there's a problem with marshall
		return ""
	}
	return string(json)
}

var TestMetadata = &models.Metadata{"test-key": "test-value"}

var Xpub = &models.Xpub{
	Model:           common.Model{Metadata: *TestMetadata},
	ID:              "cba0be1e753a7609e1a2f792d2e80ea6fce241be86f0690ec437377477809ccc",
	CurrentBalance:  16680,
	NextInternalNum: 2,
	NextExternalNum: 1,
}

var AccessKey = &models.AccessKey{
	Model:  common.Model{Metadata: *TestMetadata},
	ID:     "access-key-id",
	XpubID: Xpub.ID,
	Key:    AccessKeyString,
}

var Destination = &models.Destination{
	Model:         common.Model{Metadata: *TestMetadata},
	ID:            "90d10acb85f37dd009238fe7ec61a1411725825c82099bd8432fcb47ad8326ce",
	XpubID:        Xpub.ID,
	LockingScript: "76a9140e0eb4911d79e9b7683f268964f595b66fa3604588ac",
	Type:          "pubkeyhash",
	Chain:         1,
	Num:           19,
	Address:       "18oETbMcqRB9S7NEGZgwsHKpoTpB3nKBMa",
	DraftID:       "3a0e1fdd9ac6046c0c82aa36b462e477a455880ceeb08d3aabb1bf031553d1df",
}

var Transaction = &models.Transaction{
	Model:                common.Model{Metadata: *TestMetadata},
	ID:                   "caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda",
	Hex:                  "0100000001cf4faa628ce1abdd2cfc641c948898bb7a3dbe043999236c3ea4436a0c79f5dc000000006a47304402206aeca14175e4477031970c1cda0af4d9d1206289212019b54f8e1c9272b5bac2022067c4d32086146ca77640f02a989f51b3c6738ebfa24683c4a923f647cf7f1c624121036295a81525ba33e22c6497c0b758e6a84b60d97c2d8905aa603dd364915c3a0effffffff023e030000000000001976a914f7fc6e0b05e91c3610efd0ce3f04f6502e2ed93d88ac99030000000000001976a914550e06a3aa71ba7414b53922c13f96a882bf027988ac00000000",
	XpubInIDs:            []string{Xpub.ID},
	XpubOutIDs:           []string{Xpub.ID},
	BlockHash:            "00000000000000000896d2b93efa4476c4bd47ed7a554aeac6b38044745a6257",
	BlockHeight:          825599,
	Fee:                  97,
	NumberOfInputs:       4,
	NumberOfOutputs:      2,
	DraftID:              "fe6fe12c25b81106b7332d58fe87dab7bc6e56c8c21ca45b4de05f673f3f653c",
	TotalValue:           6955,
	OutputValue:          1725,
	Outputs:              map[string]int64{"680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0": 62},
	Status:               "confirmed",
	TransactionDirection: "incoming",
}

var DraftTx = &models.DraftTransaction{
	Model:  common.Model{Metadata: *TestMetadata},
	ID:     "3a0e1fdd9ac6046c0c82aa36b462e477a455880ceeb08d3aabb1bf031553d1df",
	Hex:    "010000000123462f14e60556718916a8cff9dbf2258195a928777c0373200dba1cee105bdb0100000000ffffffff020c000000000000001976a914c4b15e7f65e3e6a062c1d21b7f1d7d2cd3b18e8188ac0b000000000000001976a91455873fd2baa7b51a624f6416b1d824939d99151a88ac00000000",
	XpubID: Xpub.ID,
	Configuration: models.TransactionConfig{
		ChangeDestinations:         []*models.Destination{Destination},
		ChangeStrategy:             "",
		ChangeMinimumSatoshis:      0,
		ChangeNumberOfDestinations: 0,
		ChangeSatoshis:             11,
		Fee:                        1,
		FeeUnit: &models.FeeUnit{
			Satoshis: 1,
			Bytes:    1000,
		},
		FromUtxos: []*models.UtxoPointer{{
			TransactionID: "caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda",
			OutputIndex:   1,
		}},
		IncludeUtxos: []*models.UtxoPointer{{
			TransactionID: "caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda",
			OutputIndex:   1,
		}},
		Inputs: []*models.TransactionInput{{
			Utxo: models.Utxo{
				UtxoPointer: models.UtxoPointer{
					TransactionID: "db5b10ee1cba0d2073037c7728a9958125f2dbf9cfa81689715605e6142f4623",
					OutputIndex:   1,
				},
				ID:           "041479f86c475603fd510431cf702bc8c9849a9c350390eb86b467d82a13cc24",
				XpubID:       "9fe44728bf16a2dde3748f72cc65ea661f3bf18653b320d31eafcab37cf7fb36",
				Satoshis:     24,
				ScriptPubKey: "76a914673d3a53dade2723c48b446578681e253b5c548b88ac",
				Type:         "pubkeyhash",
				DraftID:      "3a0e1fdd9ac6046c0c82aa36b462e477a455880ceeb08d3aabb1bf031553d1df",
				SpendingTxID: "",
			},
			Destination: *Destination,
		}},
		Outputs: []*models.TransactionOutput{
			{
				PaymailP4: &models.PaymailP4{
					Alias:           "dorzepowski",
					Domain:          "damiano.4chain.space",
					FromPaymail:     "test3@kuba.4chain.space",
					Note:            "paymail_note",
					PubKey:          "1DSsgJdB2AnWaFNgSbv4MZC2m71116JafG",
					ReceiveEndpoint: "https://damiano.serveo.net/v1/bsvalias/receive-transaction/{alias}@{domain.tld}",
					ReferenceID:     "9b48dde1821fa82cf797372a297363c8",
					ResolutionType:  "p2p",
				},
				Satoshis: 12,
				Scripts: []*models.ScriptOutput{{
					Address:    "1Jw1vRUq6pYqiMBAT6x3wBfebXCrXv6Qbr",
					Satoshis:   12,
					Script:     "76a914c4b15e7f65e3e6a062c1d21b7f1d7d2cd3b18e8188ac",
					ScriptType: "pubkeyhash",
				}},
				To:           "pubkeyhash",
				UseForChange: false,
			},
			{
				Satoshis: 11,
				Scripts: []*models.ScriptOutput{{
					Address:    "18oETbMcqRB9S7NEGZgwsHKpoTpB3nKBMa",
					Satoshis:   11,
					Script:     "76a91455873fd2baa7b51a624f6416b1d824939d99151a88ac",
					ScriptType: "pubkeyhash",
				}},
				To: "18oETbMcqRB9S7NEGZgwsHKpoTpB3nKBMa",
			},
		},
		SendAllTo: &models.TransactionOutput{
			OpReturn: &models.OpReturn{
				Hex:      "0100000001cf4faa628ce1abdd2cfc641c948898bb7a3dbe043999236c3ea4436a0c79f5dc000000006a47304402206aeca14175e4477031970c1cda0af4d9d1206289212019b54f8e1c9272b5bac2022067c4d32086146ca77640f02a989f51b3c6738ebfa24683c4a923f647cf7f1c624121036295a81525ba33e22c6497c0b758e6a84b60d97c2d8905aa603dd364915c3a0effffffff023e030000000000001976a914f7fc6e0b05e91c3610efd0ce3f04f6502e2ed93d88ac99030000000000001976a914550e06a3aa71ba7414b53922c13f96a882bf027988ac00000000",
				HexParts: []string{"0100000001cf4faa628ce1abdd2cfc641c948898bb7a3dbe043999236c3ea4436a0c79f5dc000000006a47304402206aeca14175e4477031970c1cda0af4d9d1206289212019b54f8e1c9272b5bac2022067c4d32086146ca77640f02a989f51b3c6738ebfa24683c4a923f647cf7f1c624121036295a81525ba33e22c6497c0b758e6a84b60d97c2d8905aa603dd364915c3a0effffffff023e030000000000001976a914f7fc6e0b05e91c3610efd0ce3f04f6502e2ed93d88ac99030000000000001976a914550e06a3aa71ba7414b53922c13f96a882bf027988ac00000000"},
				Map: &models.MapProtocol{
					App:  "app_protocol",
					Keys: map[string]interface{}{"test-key": "test-value"},
					Type: "app_protocol_type",
				},
				StringParts: []string{"string", "parts"},
			},
			PaymailP4: &models.PaymailP4{
				Alias:           "alias",
				Domain:          "domain.tld",
				FromPaymail:     "alias@paymail.com",
				Note:            "paymail_note",
				PubKey:          "1DSsgJdB2AnWaFNgSbv4MZC2m71116JafG",
				ReceiveEndpoint: "https://bsvalias.example.org/alias@domain.tld/payment-destination-response",
				ReferenceID:     "3d7c2ca83a46",
				ResolutionType:  "resolution_type",
			},
			Satoshis: 1220,
			Script:   "script",
			Scripts: []*models.ScriptOutput{{
				Address:    "12HL5RyEy3Rt6SCwxgpiFSTigem1Pzbq22",
				Satoshis:   1220,
				Script:     "script",
				ScriptType: "pubkeyhash",
			}},
			To:           "1DSsgJdB2AnWaFNgSbv4MZC2m71116JafG",
			UseForChange: false,
		},
		Sync: &models.SyncConfig{
			Broadcast:        true,
			BroadcastInstant: true,
			PaymailP2P:       true,
			SyncOnChain:      true,
		},
	},
	Status:    "draft",
	FinalTxID: "caae6e799210dfea7591e3d55455437eb7e1091bb01463ae1e7ddf9e29c75eda",
}
