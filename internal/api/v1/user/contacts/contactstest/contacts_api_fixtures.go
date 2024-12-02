package contactstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedUserContactsPage(t *testing.T) *queries.UserContactsPage {
	return &queries.UserContactsPage{
		Content: []*response.Contact{
			{
				Model: response.Model{
					CreatedAt: ParseTime(t, "2024-10-18T12:07:44.739839Z"),
					UpdatedAt: ParseTime(t, "2024-10-18T15:08:44.739918Z"),
				},
				ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
				FullName: "John Doe",
				Paymail:  "john.doe.test5@john.doe.4chain.space",
				PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
				Status:   "unconfirmed",
			},
			{
				Model: response.Model{
					CreatedAt: ParseTime(t, "2024-10-18T12:07:44.739839Z"),
					UpdatedAt: ParseTime(t, "2024-10-18T15:08:44.739918Z"),
				},
				ID:       "e55a4d4e-4a4b-4720-8556-1c00dd6a5cf3",
				FullName: "Jane Doe",
				Paymail:  "jane.doe.test5@jane.doe.4chain.space",
				PubKey:   "f8898969-3f96-48d3-b122-bbb3e738dbf5",
				Status:   "unconfirmed",
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

func ExpectedContactWithWithPaymail(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-18T12:07:44.739839Z"),
			UpdatedAt: ParseTime(t, "2024-10-18T15:08:44.739918Z"),
		},
		ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
		FullName: "John Doe",
		Paymail:  "john.doe.test5@john.doe.4chain.space",
		PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
		Status:   "unconfirmed",
	}
}

func ExpectedUpsertContact(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-18T12:07:44.739839Z"),
			UpdatedAt: ParseTime(t, "2024-11-06T11:30:35.090124Z"),
			Metadata: map[string]interface{}{
				"example_key": "example_val",
			},
		},
		ID:       "68acf78f-5ece-4917-821d-8028ecf06c9a",
		FullName: "John Doe",
		Paymail:  "john.doe.test@john.doe.test.4chain.space",
		PubKey:   "0df36839-67bb-49e7-a9c7-e839aa564871",
		Status:   "unconfirmed",
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
