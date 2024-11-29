package merklerootstest

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func ExpectedMerkleRootsPage() *queries.MerkleRootPage {
	return &queries.MerkleRootPage{
		Content: []models.MerkleRoot{
			{
				MerkleRoot:  "d02ab7b5-ac3e-4612-9377-9bffe05ac689",
				BlockHeight: 1,
			},
			{
				MerkleRoot:  "132a2a38-b23f-404b-940f-f811de886114",
				BlockHeight: 2,
			},
			{
				MerkleRoot:  "d229c224-6c21-4c68-ba25-261119e9b8dc",
				BlockHeight: 3,
			},
		},
		Page: models.ExclusiveStartKeyPageInfo{
			OrderByField:     Ptr("blockHeight"),
			SortDirection:    Ptr("asc"),
			TotalElements:    10,
			Size:             20,
			LastEvaluatedKey: "6bad63f5-8f2e-4756-aca9-cc9cb4a001c6",
		},
	}
}

func Ptr[T any](value T) *T {
	return &value
}

func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
