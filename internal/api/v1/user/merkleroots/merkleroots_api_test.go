package merkleroots_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMerkleRootsAPI_MerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.MerkleRootPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/merkleroots response: 200": {
			expectedResponse: merklerootstest.ExpectedMerkleRootsPage(),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_200.json")),
		},
		"HTTP GET /api/v1/merkleroots response: 400": {
			expectedErr: spvwallettest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, spvwallettest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/merkleroots str response: 500": {
			expectedErr: spvwallettest.NewInternalServerSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/merkleroots"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := spvWalletClient.MerkleRoots(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

// Mock repository for testing
type MockMerkleRootsRepository struct {
	mock.Mock
}

// GetLastMerkleRoot retrieves the last Merkle root from storage.
func (m *MockMerkleRootsRepository) GetLastMerkleRoot() string {
	args := m.Called()
	return args.String(0)
}

// SaveMerkleRoots appends synced Merkle roots to the simulated storage.
func (m *MockMerkleRootsRepository) SaveMerkleRoots(roots []models.MerkleRoot) error {
	args := m.Called(roots)
	return args.Error(0)
}

// TestSyncMerkleRoots tests the SyncMerkleRoots functionality
func TestMerkleRootsAPI_SyncMerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder   httpmock.Responder
		setupMock   func(mockRepo *MockMerkleRootsRepository)
		expectedErr error
	}{
		"Successful Sync with Pagination": {
			responder: httpmock.ResponderFromMultipleResponses(
				[]*http.Response{
					httpmock.NewStringResponse(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_page1.json").String()),
					httpmock.NewStringResponse(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_page2.json").String()),
				},
			),
			setupMock: func(mockRepo *MockMerkleRootsRepository) {
				mockRepo.On("GetLastMerkleRoot").Return("") // Start with no data
				mockRepo.On("SaveMerkleRoots", mock.MatchedBy(func(roots []models.MerkleRoot) bool {
					return len(roots) > 0
				})).Return(nil).Twice() // Called twice for two pages
			},
		},
		"Stale LastEvaluatedKey Error": {
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_stale.json")),
			setupMock: func(mockRepo *MockMerkleRootsRepository) {
				mockRepo.On("GetLastMerkleRoot").Return("stale-key") // Simulate a stale key
				mockRepo.On("SaveMerkleRoots", mock.Anything).Return(nil)
			},
			expectedErr: errors.ErrStaleLastEvaluatedKey,
		},
		"API Returns Error Response": {
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, spvwallettest.NewInternalServerSPVError()),
			setupMock: func(mockRepo *MockMerkleRootsRepository) {
				mockRepo.On("GetLastMerkleRoot").Return("") // No data initially
			},
			expectedErr: spvwallettest.NewInternalServerSPVError(),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/merkleroots"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mockRepo := new(MockMerkleRootsRepository)
			spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)
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
	url := spvwallettest.TestAPIAddr + "/api/v1/merkleroots"
	spvWalletClient, transport := spvwallettest.GivenSPVUserAPI(t)

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
