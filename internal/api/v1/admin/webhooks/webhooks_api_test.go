package webhooks_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestWebhooksAPI_SubscribeWebhook(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 200": {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP POST /api/v1/admin/webhooks/subscriptions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/webhooks/subscriptions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			err := wallet.SubscribeWebhook(context.Background(), &commands.CreateWebhookSubscription{
				URL:         "http://webhook1.com",
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
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		"HTTP DELETE /api/v1/admin/webhooks/subscriptions str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/webhooks/subscriptions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.UnsubscribeWebhook(context.Background(), &commands.CancelWebhookSubscription{
				URL: "http://webhook1.com",
			})
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
