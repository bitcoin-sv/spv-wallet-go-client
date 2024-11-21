package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *queries.UserContactsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/contacts response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: contactstest.ExpectedUserContactsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contacts_200.json")),
		},
		"HTTP GET /api/v1/contacts response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/contacts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

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
			statusCode:       http.StatusOK,
			expectedResponse: contactstest.ExpectedContactWithWithPaymail(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contact_paymail_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

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
			statusCode:       http.StatusOK,
			expectedResponse: contactstest.ExpectedUpsertContact(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/put_contact_upsert_200.json")),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPut, url, tc.responder)

			// then:
			got, err := wallet.UpsertContact(context.Background(), commands.UpsertContact{
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
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.RemoveContact(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_ConfirmContact(t *testing.T) {
	contact := &models.Contact{
		Paymail: "alice@example.com",
		PubKey:  clienttest.MockPKI(t, clienttest.UserXPub),
	}

	tests := map[string]struct {
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 200", contact.Paymail): {
			statusCode: http.StatusOK,
			responder:  httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 400", contact.Paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusBadRequest,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation str response: 500", contact.Paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts/" + contact.Paymail + "/confirmation"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wrappedTransport := clienttest.NewTransportWrapper()
			aliceClient, _ := clienttest.GivenSPVWalletClientWithTransport(t, wrappedTransport)
			wrappedTransport.RegisterResponder(http.MethodPost, url, tc.responder)

			passcode, err := aliceClient.GenerateTotpForContact(contact, 3600, 6)
			require.NoError(t, err)

			// then:
			err = aliceClient.ConfirmContact(context.Background(), contact, passcode, contact.Paymail, 3600, 6)
			require.ErrorIs(t, err, tc.expectedErr)

			// Assert status code:
			resp, respErr := wrappedTransport.GetResponse()
			require.NoError(t, respErr)
			require.NotNil(t, resp, "response should not be nil")
			require.Equal(t, tc.statusCode, resp.StatusCode, "unexpected HTTP status code")

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
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/contacts/" + paymail + "/confirmation"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.UnconfirmContact(context.Background(), paymail)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
