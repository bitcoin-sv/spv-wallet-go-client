package users_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/users/userstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyAPI_GenerateAccessKey(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		"HTTP POST /api/v1/users/current/keys response: 200": {
			expectedResponse: userstest.ExpectedCreatedAccessKey(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/post_access_key_200.json")),
		},
		"HTTP POST /api/v1/users/current/keys response: 400": {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/users/current/keys str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current/keys"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			got, err := wallet.GenerateAccessKey(context.Background(), &commands.GenerateAccessKey{
				Metadata: map[string]any{"example_key": "example_value"},
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestAccessKeyAPI_AccessKey(t *testing.T) {
	ID := "1fb70cc2-e9d9-41a3-842e-f71cc58d9787"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 200", ID): {
			expectedResponse: userstest.ExpectedRertrivedAccessKey(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/get_access_key_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 400", ID): {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s str response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current/keys/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.AccessKey(context.Background(), ID)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestAccessKeyAPI_AccessKeys(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.AccessKeyPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/users/current/keys response: 200": {
			expectedResponse: userstest.ExpectedAccessKeyPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/get_access_keys_200.json")),
		},
		"HTTP GET /api/v1/users/current/keys response: 400": {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/users/current/keys str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current/keys"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.AccessKeys(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestAccessKeyAPI_RevokeAccessKey(t *testing.T) {
	ID := "081743f7-040e-45a3-8c36-dde39001e20d"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 200", ID): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 400", ID): {
			expectedErr: userstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s str response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/users/current/keys/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.RevokeAccessKey(context.Background(), ID)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
