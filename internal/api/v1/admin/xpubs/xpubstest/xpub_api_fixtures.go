package xpubstest

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedXPub(t *testing.T) *response.Xpub {
	return &response.Xpub{
		Model: response.Model{
			CreatedAt: testutils.ParseTime(t, "2024-11-22T07:51:37.708754Z"),
			UpdatedAt: testutils.ParseTime(t, "2024-11-22T08:51:37.708865+01:00"),
			Metadata:  map[string]any{"key": "value"},
		},
		ID:              "d7ff33b6-8c25-4955-bcea-a5557c18bb95",
		CurrentBalance:  0,
		NextInternalNum: 0,
		NextExternalNum: 0,
	}
}

func ExpectedXPubsPage(t *testing.T) *queries.XPubPage {
	return &queries.XPubPage{
		Content: []*response.Xpub{
			{
				Model: response.Model{
					CreatedAt: testutils.ParseTime(t, "2024-11-21T11:41:49.830635Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-11-21T11:41:49.830649Z"),
					Metadata:  map[string]any{"key": "val"},
				},
				ID:              "3c7a9d02-32e3-4d83-a391-af64f1933acb",
				CurrentBalance:  10,
				NextInternalNum: 20,
				NextExternalNum: 30,
			},
			{
				Model: response.Model{
					CreatedAt: testutils.ParseTime(t, "2024-11-21T11:26:43.091808Z"),
					UpdatedAt: testutils.ParseTime(t, "2024-11-21T11:26:43.091857Z"),
					Metadata:  map[string]any{"key": "val"},
				},
				ID:              "301f38e2-f1dc-43cb-9db2-f2835a648b8b",
				CurrentBalance:  40,
				NextInternalNum: 50,
				NextExternalNum: 60,
			},
		},
		Page: response.PageDescription{
			Size:          50,
			Number:        1,
			TotalElements: 40,
			TotalPages:    1,
		},
	}
}
