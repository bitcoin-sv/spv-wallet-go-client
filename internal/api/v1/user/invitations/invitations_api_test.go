package invitations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	invitationsURL = "/api/v1/invitations"
	paymail        = "john.doe.test@john.doe.test.4chain.space"
	contactsURI    = "/contacts"
)

func TestInvitationsAPI_AcceptInvitation(t *testing.T) {

	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts response: 200", paymail): {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts response: 400", paymail): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts response: 500", paymail): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, invitationsURL, paymail, contactsURI)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			err := wallet.AcceptInvitation(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestInvitationsAPI_RejectInvitation(t *testing.T) {

	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s response: 200", paymail): {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s response: 400", paymail): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s response: 500", paymail): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s str response: 500", paymail): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, invitationsURL, paymail)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.RejectInvitation(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
