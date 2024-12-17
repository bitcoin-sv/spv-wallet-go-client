package paymails_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/paymails/paymailstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	paymailsURL = "/api/v1/admin/paymails"
	xpubID      = "xpub22e6cba6-ef6e-432a-8612-63ac4b290ce9"
	id          = "98dbafe0-4e2b-4307-8fbf-c55209214bae"
)

func TestPaymailsAPI_DeletePaymail(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PaymailAddress
		expectedErr      error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 200", xpubID): {
			expectedResponse: paymailstest.ExpectedCreatedPaymail(t),
			responder:        testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 400", xpubID): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s response: 500", xpubID): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/paymails/%s str response: 500", xpubID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, paymailsURL, xpubID)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.DeletePaymail(context.Background(), xpubID)
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
			responder:        testutils.NewJSONFileResponderWithStatusOK("paymailstest/post_paymail_200.json"),
		},
		"HTTP POST /api/v1/admin/paymails response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/paymails response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/paymails str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, paymailsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			got, err := wallet.CreatePaymail(context.Background(), &commands.CreatePaymail{
				Key: xpubID,
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestPaymailsAPI_Paymails(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.PaymailsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/paymails response: 200": {
			expectedResponse: paymailstest.ExpectedPaymailsPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("paymailstest/get_paymails_200.json"),
		},
		"HTTP GET /api/v1/admin/paymails response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/paymails response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/paymails str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, paymailsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Paymails(context.Background(), queries.QueryWithPageFilter[filter.AdminPaymailFilter](filter.Page{Size: 1}))
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestPaymailsAPI_Paymail(t *testing.T) {

	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.PaymailAddress
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 200", id): {
			expectedResponse: paymailstest.ExpectedPaymail(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("paymailstest/get_paymail_200.json"),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/admin/paymails/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, paymailsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Paymail(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
