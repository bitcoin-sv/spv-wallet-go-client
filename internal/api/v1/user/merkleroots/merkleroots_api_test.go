package merkleroots_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const merkleRootsURL = "/api/v1/merkleroots"

func TestMerkleRootsAPI_MerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.MerkleRootPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/merkleroots response: 200": {
			expectedResponse: merklerootstest.ExpectedMerkleRootsPage(),
			responder:        testutils.NewJSONFileResponderWithStatusOK("merklerootstest/get_merkleroots_200.json"),
		},
		"HTTP GET /api/v1/merkleroots response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/merkleroots response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/merkleroots str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, merkleRootsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.MerkleRoots(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

// TestSyncMerkleRoots tests the SyncMerkleRoots functionality
func TestMerkleRootsAPI_SyncMerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		setupMock   func(mockRepo *testutils.MockMerkleRootsRepository)
		expectedErr error
	}{
		"Successful Sync with Pagination": {
			responder: testutils.NewPaginatedJSONResponder(t,
				"merklerootstest/get_merkleroots_page1.json",
				"merklerootstest/get_merkleroots_page2.json",
			),
			setupMock: func(mockRepo *testutils.MockMerkleRootsRepository) {
				testutils.SetupMerkleRootMockRepo(mockRepo, 2)
			},
		},
		"Stale LastEvaluatedKey Error": {
			responder: testutils.NewJSONFileResponderWithStatusOK("merklerootstest/get_merkleroots_stale.json"),
			setupMock: func(mockRepo *testutils.MockMerkleRootsRepository) {
				testutils.SetupStaleKeyMock(mockRepo)
			},
			expectedErr: errors.ErrStaleLastEvaluatedKey,
		},
		"API Returns Error Response": {
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, testutils.NewInternalServerSPVError()),
			setupMock: func(mockRepo *testutils.MockMerkleRootsRepository) {
				testutils.SetupEmptyMerkleRootMock(mockRepo)
			},
			expectedErr: testutils.NewInternalServerSPVError(),
		},
	}

	url := testutils.FullAPIURL(t, merkleRootsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mockRepo := new(testutils.MockMerkleRootsRepository)
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			spvWalletClient, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)
			tc.setupMock(mockRepo)

			// when:
			err := spvWalletClient.SyncMerkleRoots(context.Background(), mockRepo)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestMerkleRootsAPI_SyncMerkleRoots_PartialResponsesStoredSuccessfully(t *testing.T) {
	// given:
	db := merklerootstest.CreateRepository([]models.MerkleRoot{})
	url := testutils.FullAPIURL(t, merkleRootsURL)
	spvWalletClient, transport := testutils.GivenSPVUserAPI(t)

	var expected []models.MerkleRoot
	expected = append(expected, merklerootstest.FirstMerkleRootsPage().Content...)
	expected = append(expected, merklerootstest.SecondMerkleRootsPage().Content...)
	expected = append(expected, merklerootstest.ThirdMerkleRootsPage().Content...)

	transport.RegisterResponder(http.MethodGet, url, merklerootstest.ResponderWithThreeMerkleRootPagesSuccess(t))

	// when:
	err := spvWalletClient.SyncMerkleRoots(context.Background(), db)

	// then:
	require.NoError(t, err)
	require.Equal(t, expected, db.MerkleRoots)
}
