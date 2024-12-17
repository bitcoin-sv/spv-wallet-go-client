package paymailstest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedCreatedPaymail(t *testing.T) *response.PaymailAddress {
	return &response.PaymailAddress{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-12-02T10:22:45.263654Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-12-02T11:22:45.263664+01:00"),
		},
		ID:         "069d0011-580e-4fc6-9f24-45471b732a8b",
		XpubID:     "22e6cba6-ef6e-432a-8612-63ac4b290ce9",
		Alias:      "john.doe.test",
		Domain:     "example.com",
		PublicName: "john.doe.test",
		Avatar:     "",
	}
}

func ExpectedPaymail(t *testing.T) *response.PaymailAddress {
	return &response.PaymailAddress{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-10-02T10:28:15.544234Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-10-02T10:34:54.836433Z"),
		},
		ID:         "98dbafe0-4e2b-4307-8fbf-c55209214bae",
		XpubID:     "0d71ac87-ef56-4b1a-8372-814481cface6",
		Alias:      "john.doe.test",
		Domain:     "john.doe.test.4chain.space",
		PublicName: "john.doe.test",
		Avatar:     "http://localhost:3003/static/paymail/avatar.jpg",
	}
}

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
