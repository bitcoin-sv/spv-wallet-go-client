package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse []*response.Contact
		expectedErr      error
	}{
		"HTTP GET /api/v1/contacts response: 200": {
			statusCode: http.StatusOK,
			expectedResponse: []*response.Contact{
				{
					Model: response.Model{
						CreatedAt: ParseTime("2024-10-18T12:07:44.739839Z"),
						UpdatedAt: ParseTime("2024-10-18T15:08:44.739918Z"),
					},
					ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
					FullName: "John Doe",
					Paymail:  "john.doe.test5@john.doe.4chain.space",
					PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
					Status:   "unconfirmed",
				},
				{
					Model: response.Model{
						CreatedAt: ParseTime("2024-10-18T12:07:44.739839Z"),
						UpdatedAt: ParseTime("2024-10-18T15:08:44.739918Z"),
					},
					ID:       "e55a4d4e-4a4b-4720-8556-1c00dd6a5cf3",
					FullName: "Jane Doe",
					Paymail:  "jane.doe.test5@jane.doe.4chain.space",
					PubKey:   "f8898969-3f96-48d3-b122-bbb3e738dbf5",
					Status:   "unconfirmed",
				},
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/contacts_200.json")),
		},
		"HTTP GET /api/v1/contacts response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/contacts str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.Contacts(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_ContactWithPaymail(t *testing.T) {
	paymail := "john.doe.test5@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s response: 200", paymail): {
			statusCode: http.StatusOK,
			expectedResponse: &response.Contact{
				Model: response.Model{
					CreatedAt: ParseTime("2024-10-18T12:07:44.739839Z"),
					UpdatedAt: ParseTime("2024-10-18T15:08:44.739918Z"),
				},
				ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
				FullName: "John Doe",
				Paymail:  "john.doe.test5@john.doe.4chain.space",
				PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
				Status:   "unconfirmed",
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/contact_paymail_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.ContactWithPaymail(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_UpsertContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s response: 200", paymail): {
			statusCode: http.StatusOK,
			expectedResponse: &response.Contact{
				Model: response.Model{
					CreatedAt: ParseTime("2024-10-18T12:07:44.739839Z"),
					UpdatedAt: ParseTime("2024-11-06T11:30:35.090124Z"),
					Metadata: map[string]interface{}{
						"example_key": "example_val",
					},
				},
				ID:       "68acf78f-5ece-4917-821d-8028ecf06c9a",
				FullName: "John Doe",
				Paymail:  "john.doe.test@john.doe.test.4chain.space",
				PubKey:   "0df36839-67bb-49e7-a9c7-e839aa564871",
				Status:   "unconfirmed",
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/contact_upsert_200.json")),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPut, URL, tc.responder)

			// then:
			got, err := wallet.UpsertContact(context.Background(), client.UpsertContactArgs{
				FullName: "John Doe",
				Metadata: map[string]any{"example_key": "example_val"},
				Paymail:  paymail,
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_RemoveContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s response: 200", paymail): {
			statusCode: http.StatusOK,
			responder:  httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodDelete, URL, tc.responder)

			// then:
			err := wallet.RemoveContact(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_ConfirmContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 200", paymail): {
			statusCode: http.StatusOK,
			responder:  httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation str response: 500", paymail): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts/" + paymail + "/confirmation"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			err := wallet.ConfirmContact(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_UnconfirmContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation response: 200", paymail): {
			statusCode: http.StatusOK,
			responder:  httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation str response: 500", paymail): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/contacts/" + paymail + "/confirmation"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodDelete, URL, tc.responder)

			// then:
			err := wallet.UnconfirmContact(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func ParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return t
}

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
