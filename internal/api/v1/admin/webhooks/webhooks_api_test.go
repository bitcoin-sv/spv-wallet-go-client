package webhooks_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/notifications"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	webhooksURL = "/api/v1/admin/webhooks/subscriptions"
	url         = "http://webhook1.com"
)

func TestWebhooksAPI_SubscribeWebhook(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 200": {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, webhooksURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			err := wallet.SubscribeWebhook(context.Background(), &commands.CreateWebhookSubscription{
				URL:         url,
				TokenHeader: "Header",
				TokenValue:  "76dd388f-62de-4957-afae-967c3a424bc7",
			})
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestWebhooksAPI_UnsubscribeWebhook(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions response: 200": {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, webhooksURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.UnsubscribeWebhook(context.Background(), &commands.CancelWebhookSubscription{
				URL: url,
			})
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestWebhooksAPI_GetAllWebhooks(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedErr      error
		expectedResponse []*notifications.Webhook
	}{
		"HTTP GET /api/v1/admin/webhooks/subscriptions response: 200": {
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, []map[string]interface{}{
				{
					"url":         "http://webhook1.com",
					"tokenHeader": "Authorization",
					"tokenValue":  "abcd1234",
				},
				{
					"url":         "http://webhook2.com",
					"tokenHeader": "Auth",
					"tokenValue":  "xyz5678",
				},
			}),
			expectedResponse: []*notifications.Webhook{
				{
					URL:         url,
					TokenHeader: "Authorization",
					TokenValue:  "abcd1234",
				},
				{
					URL:         "http://webhook2.com",
					TokenHeader: "Auth",
					TokenValue:  "xyz5678",
				},
			},
		},
		"HTTP GET /api/v1/admin/webhooks/subscriptions response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/webhooks/subscriptions response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/admin/webhooks/subscriptions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, webhooksURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			webhooks, err := wallet.GetAllWebhooks(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.ElementsMatch(t, tc.expectedResponse, webhooks)
		})
	}
}
