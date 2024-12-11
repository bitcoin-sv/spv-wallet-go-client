package testutils

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/mock"
)

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
	return fmt.Errorf("SaveMerkleRoots error: %w", args.Error(0))
}

// SetupMerkleRootMockRepo sets up a mock MerkleRootsRepository with common behavior
func SetupMerkleRootMockRepo(mockRepo *MockMerkleRootsRepository, expectedCallCount int) {
	mockRepo.On("GetLastMerkleRoot").Return("") // Start with no data

	// Match any non-empty slice of MerkleRoot models
	mockRepo.On("SaveMerkleRoots", mock.MatchedBy(func(roots []models.MerkleRoot) bool {
		return len(roots) > 0
	})).Return(nil).Times(expectedCallCount)
}

// SetupStaleKeyMock sets up the mock repository for a stale LastEvaluatedKey scenario
func SetupStaleKeyMock(mockRepo *MockMerkleRootsRepository) {
	mockRepo.On("GetLastMerkleRoot").Return("stale-key") // Simulate a stale key
	mockRepo.On("SaveMerkleRoots", mock.Anything).Return(nil)
}

// SetupEmptyMerkleRootMock sets up a mock MerkleRootsRepository with no data
func SetupEmptyMerkleRootMock(mockRepo *MockMerkleRootsRepository) {
	mockRepo.On("GetLastMerkleRoot").Return("") // Simulate no data
}
