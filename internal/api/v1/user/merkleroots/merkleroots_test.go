package merkleroots_test

import (
	"context"
	"net/http"
	"testing"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestMerkleRootsAPI_MerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse []*models.MerkleRoot
		expectedErr      error
	}{
		"HTTP GET /api/v1/merkleroots response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: merklerootstest.ExpectedMerkleRoorts(),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("merklerootstest/merkleroots_200.json")),
		},
		"HTTP GET /api/v1/merkleroots response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, merklerootstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/merkleroots str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/merkleroots"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.MerkleRoots(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func Test_CreateQueryParams(t *testing.T) {
	tests := map[string]struct {
		opts           []queries.MerkleRootsQueryOption
		expectedParams map[string]string
	}{

		"query params: map entry [batchSize]=10": {
			opts: []queries.MerkleRootsQueryOption{queries.MerkleRootsQueryWithBatchSize(10)},
			expectedParams: map[string]string{
				"batchSize": "10",
			},
		},
		"query params: map entry [lastEvaluatedKey]=key": {
			opts: []queries.MerkleRootsQueryOption{queries.MerkleRootsQueryWithLastEvaluatedKey("key")},
			expectedParams: map[string]string{
				"lastEvaluatedKey": "key",
			},
		},
		"query params: map entries [lastEvaluatedKey]=key, [batchSize]=10": {
			opts: []queries.MerkleRootsQueryOption{
				queries.MerkleRootsQueryWithLastEvaluatedKey("key"),
				queries.MerkleRootsQueryWithBatchSize(10),
			},
			expectedParams: map[string]string{
				"lastEvaluatedKey": "key",
				"batchSize":        "10",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := merkleroots.CreateQueryParams(tc.opts...)
			require.Equal(t, got, tc.expectedParams)
		})
	}
}
