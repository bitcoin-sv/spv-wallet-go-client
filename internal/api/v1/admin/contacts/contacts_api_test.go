package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	contactsURL = "/api/v1/admin/contacts"
	id          = "4d570959-dd85-4f53-bad1-18d0671761e9"
)

func TestContactsAPI_CreateContact(t *testing.T) {
	paymail := "john.doe@test.4chain.space"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 200", paymail): {
			expectedResponse: contactstest.ExpectedCreatedContact(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("contactstest/post_contact_200.json"),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 400", paymail): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 404", paymail): {
			expectedErr: testutils.NewResourceNotFoundSPVError(),
			responder:   testutils.NewResourceNotFoundSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 409", paymail): {
			expectedErr: testutils.NewConflictRequestSPVError(),
			responder:   testutils.NewConflictRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 500", paymail): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/contacts/%s response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, contactsURL+"/"+paymail, tc.responder)

			// when:
			got, err := wallet.CreateContact(context.Background(), &commands.CreateContact{
				CreatorPaymail: "admin@test.4chain.space",
				Paymail:        paymail,
				FullName:       "John Doe",
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.ContactsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/contacts response: 200": {
			expectedResponse: contactstest.ExpectedContactsPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("contactstest/get_contacts_200.json"),
		},
		"HTTP GET /api/v1/admin/contacts response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/contacts response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/contacts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, contactsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			opts := []queries.QueryOption[filter.AdminContactFilter]{
				queries.QueryWithPageFilter[filter.AdminContactFilter](filter.Page{
					Number: 1,
					Size:   1,
					Sort:   "asc",
					SortBy: "key",
				}),
				queries.QueryWithFilter(filter.AdminContactFilter{
					ContactFilter: filter.ContactFilter{
						ModelFilter: filter.ModelFilter{
							IncludeDeleted: testutils.Ptr(true),
						},
					},
				}),
			}
			params := "page=1&size=1&sort=asc&sortBy=key&includeDeleted=true"
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponderWithQuery(http.MethodGet, url, params, tc.responder)

			// when:
			got, err := wallet.Contacts(context.Background(), opts...)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_ConfirmContacts(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		"HTTP POST /api/v1/admin/contacts/confirmations response: 200": {
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		"HTTP POST /api/v1/admin/contacts/confirmations response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/contacts/confirmations response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/contacts/confirmations str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, contactsURL+"/confirmations")
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			err := wallet.ConfirmContacts(context.Background(), &commands.ConfirmContacts{
				PaymailA: "alice@paymail.com",
				PaymailB: "bob@paymail.com",
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_ContactUpdate(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 200", id): {
			expectedResponse: contactstest.ExpectedUpdatedUserContact(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("contactstest/put_contact_update_200.json"),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, contactsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPut, url, tc.responder)

			// when:
			got, err := wallet.ContactUpdate(context.Background(), &commands.UpdateContact{
				ID:       id,
				FullName: "John Doe Williams",
				Metadata: map[string]any{"phoneNumber": "123456789"},
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_DeleteContact(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s response: 200", id): {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE/api/v1/admin/contacts/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, contactsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.DeleteContact(context.Background(), id)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
