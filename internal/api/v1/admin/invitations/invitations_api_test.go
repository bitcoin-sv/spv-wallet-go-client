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
	invitationsURL = "/api/v1/admin/invitations"
	id             = "34d0b1f9-6d00-4bdb-ba2e-146a3cbadd35"
)

func TestInvitationsAPI_AcceptInvitation(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 200", id): {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, invitationsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			err := wallet.AcceptInvitation(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestInvitationsAPI_RejectInvitation(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 200", id): {
			responder: testutils.NewStringResponderStatusOK(http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, invitationsURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.RejectInvitation(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
