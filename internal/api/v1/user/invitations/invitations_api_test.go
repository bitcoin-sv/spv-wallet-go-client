package invitations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/invitations/invitationstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestInvitationsAPI_AcceptInvitation(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts response: 200", paymail): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts response: 400", paymail): {
			expectedErr: invitationstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, invitationstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts str response: 500", paymail): {
			expectedErr: invitationstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, invitationstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/invitations/" + paymail + "/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			err := wallet.AcceptInvitation(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestInvitationsAPI_RejectInvitation(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s response: 200", paymail): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s response: 400", paymail): {
			expectedErr: invitationstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, invitationstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s str response: 500", paymail): {
			expectedErr: invitationstest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, invitationstest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/invitations/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.RejectInvitation(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
