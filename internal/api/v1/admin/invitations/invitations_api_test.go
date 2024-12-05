package invitations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestInvitationsAPI_AcceptInvitation(t *testing.T) {
	id := "34d0b1f9-6d00-4bdb-ba2e-146a3cbadd35"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 200", id): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/admin/invitations/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/invitations/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			err := wallet.AcceptInvitation(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestInvitationsAPI_RejectInvitation(t *testing.T) {
	id := "34d0b1f9-6d00-4bdb-ba2e-146a3cbadd35"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 200", id): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 400", id): {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s response: 500", id): {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/invitations/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/invitations/" + id
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.RejectInvitation(context.Background(), id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
