package walletclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

// MerkleRootsRepository is an interface responsible for saving synced merkleroots and getting last evaluat key from database
type MerkleRootsRepository interface {
	// GetLastEvaluatedKey should return the merkle root with the heighest height from your storage or undefined if empty
	GetLastEvaluatedKey() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are ordered with ascending order by block height
	SaveMerkleRoots(syncedMerkleRoots []MerkleRoot) error
}

// SyncMerkleRoots syncs merkleroots known to spv-wallet with the client database
func (wc *WalletClient) SyncMerkleRoots(ctx context.Context, repo MerkleRootsRepository, timeoutMs time.Duration) error {
	var cancel context.CancelFunc
	if timeoutMs > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeoutMs)
		defer cancel()
	}

	lastEvaluatedKey := repo.GetLastEvaluatedKey()
	requestPath := "merkleroots"
	lastEvaluatedKeyQuery := ""
	previousLastEvaluatedKey := lastEvaluatedKey

	if lastEvaluatedKey != "" {
		lastEvaluatedKeyQuery = fmt.Sprintf("?lastEvaluatedKey=%s", lastEvaluatedKey)
	}

	for {
		select {
		case <-ctx.Done():
			return ErrSyncMerkleRootsTimeout
		default:
			url := fmt.Sprintf("/%s%s", requestPath, lastEvaluatedKeyQuery)

			var merkleRootsResponse exclusiveStartKeyPage[[]MerkleRoot]

			err := wc.doHTTPRequest(ctx, http.MethodGet, url, nil, wc.xPriv, true, &merkleRootsResponse)
			if err != nil {
				return err
			}

			if previousLastEvaluatedKey == merkleRootsResponse.Page.LastEvaluatedKey {
				return ErrStaleLastEvaluatedKey
			}

			err = repo.SaveMerkleRoots(merkleRootsResponse.Content)
			if err != nil {
				return err
			}

			if merkleRootsResponse.Page.LastEvaluatedKey == "" {
				break
			}

			lastEvaluatedKeyQuery = fmt.Sprintf("?lastEvaluatedKey=%s", merkleRootsResponse.Page.LastEvaluatedKey)
			previousLastEvaluatedKey = merkleRootsResponse.Page.LastEvaluatedKey
		}
	}
	return nil
}
