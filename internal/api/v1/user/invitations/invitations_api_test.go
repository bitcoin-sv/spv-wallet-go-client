package invitations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
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
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s/contacts str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/invitations/" + paymail + "/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// then:
			err := wallet.AcceptInvitation(context.Background(), paymail)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
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
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/invitations/%s str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    errors.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: http.StatusInternalServerError,
				Code:       "internal-server-error",
			},
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/invitations/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.RejectInvitation(context.Background(), paymail)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func ParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return t
}

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
