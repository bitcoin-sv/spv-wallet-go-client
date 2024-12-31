package accesskeys_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/accesskeys/accesskeystest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	accessKeysURL = "/api/v1/users/current/keys"
	id            = "1fb70cc2-e9d9-41a3-842e-f71cc58d9787"
)

func TestAccessKeyAPI_GenerateAccessKey(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		"HTTP POST /api/v1/users/current/keys response: 200": {
			expectedResponse: accesskeystest.ExpectedCreatedAccessKey(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("accesskeystest/post_access_key_200.json"),
		},
		"HTTP POST /api/v1/users/current/keys response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP POST /api/v1/users/current/keys response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP POST /api/v1/users/current/keys str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, accessKeysURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
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
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 200", id): {
			expectedResponse: accesskeystest.ExpectedRertrivedAccessKey(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("accesskeystest/get_access_key_200.json"),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, accessKeysURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.AccessKey(context.Background(), id)

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
			expectedResponse: accesskeystest.ExpectedAccessKeyPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("accesskeystest/get_access_keys_200.json"),
		},
		"HTTP GET /api/v1/users/current/keys response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/users/current/keys response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/users/current/keys str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, accessKeysURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			opts := []queries.QueryOption[filter.AccessKeyFilter]{
				queries.QueryWithPageFilter[filter.AccessKeyFilter](filter.Page{
					Number: 1,
					Size:   1,
					Sort:   "asc",
					SortBy: "key",
				}),
				queries.QueryWithFilter(filter.AccessKeyFilter{
					ModelFilter: filter.ModelFilter{
						IncludeDeleted: testutils.Ptr(true),
					},
				}),
			}
			params := "page=1&size=1&sort=asc&sortBy=key&includeDeleted=true"
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponderWithQuery(http.MethodGet, url, params, tc.responder)

			// when:
			got, err := wallet.AccessKeys(context.Background(), opts...)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestAccessKeyAPI_RevokeAccessKey(t *testing.T) {

	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 200", id): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 400", id): {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, testutils.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 500", id): {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s str response: 500", id): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, accessKeysURL, id)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.RevokeAccessKey(context.Background(), id)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
