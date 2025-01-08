package paymailstest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedPaymailsPage(t *testing.T) *queries.PaymailsPage {
	return &queries.PaymailsPage{
		Content: []*response.PaymailAddress{
			{
				Model: response.Model{
					CreatedAt: testutils.ParseTime(t, "2024-11-18T06:50:07.144902Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-11-18T06:50:07.144932Z"),
				},
				ID:         "31b80181-4d8b-4766-9bc7-76a1d9c6b44d",
				XpubID:     "69245a3a-f9ed-4046-9acb-9d66c0b3750c",
				Alias:      "john.doe.test",
				Domain:     "john.doe.4chain.space",
				PublicName: "John Doe",
			},
			{
				Model: response.Model{
					CreatedAt: testutils.ParseTime(t, "2024-11-08T15:10:44.688653Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-11-18T07:19:51.561691Z"),
				},
				ID:         "ec91273e-9fb7-4f10-9ecb-d1848d238814",
				XpubID:     "68026cb6-a549-45e8-97b1-11426bb16769",
				Alias:      "jane.doe.test",
				Domain:     "jane.doe.4chain.space",
				PublicName: "Jane Doe",
				Avatar:     "http://localhost:3003/static/paymail/avatar.jpg",
			},
		},
		Page: response.PageDescription{
			Size:          10,
			Number:        1,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}
