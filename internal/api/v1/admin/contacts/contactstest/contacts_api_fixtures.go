package contactstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedUpdatedUserContact(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: parseTime(t, "2024-11-28T13:34:52.11722Z"),
			UpdatedAt: parseTime(t, "2024-11-29T08:23:19.66093Z"),
			Metadata:  map[string]any{"phoneNumber": "123456789"},
		},
		ID:       "4d570959-dd85-4f53-bad1-18d0671761e9",
		FullName: "John Doe Williams",
		Paymail:  "john.doe.test@john.doe.test.4chain.space",
		PubKey:   "96843af4-fc9c-4778-945d-2131ac5b1a8a",
		Status:   "awaiting",
	}
}

func ExpectedUserContactsPage(t *testing.T) *queries.UserContactsPage {
	return &queries.UserContactsPage{
		Content: []*response.Contact{
			{
				Model: response.Model{
					CreatedAt: parseTime(t, "2024-11-28T14:58:13.262238Z"),
					UpdatedAt: parseTime(t, "2024-11-28T16:18:43.842434Z"),
				},
				ID:       "7a5625ac-8256-454a-84a3-7f03f50cd7dc",
				FullName: "John Doe",
				Paymail:  "john.doe.test@john.doe.4chain.space",
				PubKey:   "bbbb7a4e-a3f4-4ca4-800a-fdd8029eda37",
				Status:   "confirmed",
			},
			{
				Model: response.Model{
					CreatedAt: parseTime(t, "2024-11-28T14:58:13.029966Z"),
					UpdatedAt: parseTime(t, "2024-11-28T14:58:13.03002Z"),
					Metadata: map[string]any{
						"phoneNumber": "123456789",
					},
				},
				ID:       "d05d2388-3c16-426d-98f1-ced9d9c5f4e1",
				FullName: "Jane Doe",
				Paymail:  "jane.doe.jane@john.doe.4chain.space",
				PubKey:   "ee191d63-1619-4fd3-ae3d-2202cfab751d",
				Status:   "unconfirmed",
			},
		},
		Page: response.PageDescription{
			Size:          50,
			Number:        1,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}

func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}

func NewInternalServerSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Code:       models.UnknownErrorCode,
	}
}

func parseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}
