package paymails_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/paymails/paymailstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestPaymailsAPI_DeletePaymail(t *testing.T) {
	id := "xpub22e6cba6-ef6e-432a-8612-63ac4b290ce9"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PaymailAddress
		expectedErr      error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 200", id): {
			expectedResponse: paymailstest.ExpectedCreatedPaymail(t),
			responder:        httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/paymails/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.DeletePaymail(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestPaymailsAPI_CreatePaymail(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PaymailAddress
		expectedErr      error
	}{
		"HTTP POST /api/v1/admin/paymails response: 200": {
			expectedResponse: paymailstest.ExpectedCreatedPaymail(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("paymailstest/post_paymail_200.json")),
		},
		"HTTP POST /api/v1/admin/paymails response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/admin/paymails response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP POST /api/v1/admin/paymails str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/paymails"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			got, err := wallet.CreatePaymail(context.Background(), &commands.CreatePaymail{
				Key: "xpub22e6cba6-ef6e-432a-8612-63ac4b290ce9",
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestPaymailsAPI_Paymails(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.PaymailAddressPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/paymails response: 200": {
			expectedResponse: paymailstest.ExpectedPaymailsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("paymailstest/get_paymails_200.json")),
		},
		"HTTP GET /api/v1/admin/paymails response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/paymails response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP GET /api/v1/admin/paymails str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/paymails"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Paymails(context.Background(), queries.PaymailQueryWithPageFilter[filter.AdminPaymailFilter](filter.Page{Size: 1}))
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestPaymailsAPI_Paymail(t *testing.T) {
	id := "98dbafe0-4e2b-4307-8fbf-c55209214bae"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PaymailAddress
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 200", id): {
			expectedResponse: paymailstest.ExpectedPaymail(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("paymailstest/get_paymail_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/paymails/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Paymail(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
