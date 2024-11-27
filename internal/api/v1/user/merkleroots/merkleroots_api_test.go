package merkleroots_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMerkleRootsAPI_MerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *queries.MerkleRootPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/merkleroots response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: merklerootstest.ExpectedMerkleRootsPage(),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_200.json")),
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
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/merkleroots"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := spvWalletClient.MerkleRoots(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
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
			responder: httpmock.NewStringResponder(http.StatusInternalServerError, "Internal Server Error"),
			setupMock: func(mockRepo *MockMerkleRootsRepository) {
				mockRepo.On("GetLastMerkleRoot").Return("") // No data initially
			},
			expectedErr: errors.ErrUnrecognizedAPIResponse,
		},
	}

	url := clienttest.TestAPIAddr + "/api/v1/merkleroots"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			mockRepo := new(MockMerkleRootsRepository)
			tc.setupMock(mockRepo)

			// Act
			err := spvWalletClient.SyncMerkleRoots(context.Background(), mockRepo)

			// Assert
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestMerkleRootsAPI_SyncMerkleRoots_PartialResponsesStoredSuccessfully tests the SyncMerkleRoots functionality
func TestMerkleRootsAPI_SyncMerkleRoots_PartialResponsesStoredSuccessfully(t *testing.T) {
	// given:
	db := merklerootstest.CreateRepository([]models.MerkleRoot{})
	url := clienttest.TestAPIAddr + "/api/v1/merkleroots"
	spvWalletClient, transport := clienttest.GivenSPVWalletClient(t)

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
