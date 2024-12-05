package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.UserContactsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/contacts response: 200": {
			expectedResponse: contactstest.ExpectedUserContactsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contacts_200.json")),
		},
		"HTTP GET /api/v1/admin/contacts response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/contacts response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/admin/contacts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Contacts(context.Background(), queries.ContactQueryWithPageFilter(filter.Page{Size: 1}))
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_ContactUpdate(t *testing.T) {
	id := "4d570959-dd85-4f53-bad1-18d0671761e9"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 200", id): {
			expectedResponse: contactstest.ExpectedUpdatedUserContact(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/put_contact_update_200.json")),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPut, url, tc.responder)

			// then:
			got, err := wallet.ContactUpdate(context.Background(), &commands.UpdateContact{
				ID:       "4d570959-dd85-4f53-bad1-18d0671761e9",
				FullName: "John Doe Williams",
				Metadata: map[string]any{"phoneNumber": "123456789"},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_DeleteContact(t *testing.T) {
	id := "4d570959-dd85-4f53-bad1-18d0671761e9"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s response: 200", id): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE/api/v1/admin/contacts/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.DeleteContact(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
