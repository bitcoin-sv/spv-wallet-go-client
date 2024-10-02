package walletclient

import (
	"context"
	"fmt"
	"net/http"
)

// exclusiveStartKeyPage represents a paginated response for database records using Exclusive Start Key paging
type exclusiveStartKeyPage[T any] struct {
	// List of records for the response
	Content T
	// Pagination details
	Page exclusiveStartKeyPageInfo
}

// exclusiveStartKeyPageInfo represents the pagination information for limiting and sorting database query results
type exclusiveStartKeyPageInfo struct {
	// Field by which to order the results
	OrderByField *string `json:"orderByField,omitempty"` // Optional ordering field
	// Direction in which to order the results (ASC or DESC)
	SortDirection *string `json:"sortDirection,omitempty"` // Optional sort direction
	// Total count of elements
	TotalElements int `json:"totalElements"`
	// Size of the page or returned data
	Size int `json:"size"`
	// Last evaluated key returned from the database
	LastEvaluatedKey string `json:"lastEvaluatedKey"`
}

// MerkleRoot holds the content of the synced Merkle root response
type MerkleRoot struct {
	MerkleRoot  string `json:"merkleRoot"`
	BlockHeight int    `json:"blockHeight"`
}

// Repository is an interface responsible for saving synced merkleroots and getting last evaluat key from database
type Repository interface {
	GetLastEvaluatedKey() string
	SaveMerkleRoots(syncedMerkleRoots []MerkleRoot) error
}

func (wc *WalletClient) SyncMerkleRoots(ctx context.Context, repo Repository) error {
	lastEvaluatedKey := repo.GetLastEvaluatedKey()
	requestPath := "merkleroots"
	lastEvaluatedKeyQuery := ""

	if lastEvaluatedKey != "" {
		lastEvaluatedKeyQuery = fmt.Sprintf("?lastEvaluatedKey=%s", lastEvaluatedKey)
	}

	for {
		url := fmt.Sprintf("/%s%s", requestPath, lastEvaluatedKeyQuery)

		var merkleRootsResponse exclusiveStartKeyPage[[]MerkleRoot]

		err := wc.doHTTPRequest(ctx, http.MethodGet, url, nil, wc.xPriv, true, &merkleRootsResponse)
		if err != nil {
			return err
		}

		err = repo.SaveMerkleRoots(merkleRootsResponse.Content)
		if err != nil {
			return err
		}

		if merkleRootsResponse.Page.LastEvaluatedKey == "" {
			break
		}

		lastEvaluatedKey = fmt.Sprintf("?lastEvaluatedKey=%s", *&merkleRootsResponse.Page.LastEvaluatedKey)
	}
	return nil
}
