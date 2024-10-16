package walletclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// MerkleRootsRepository is an interface responsible for storing synchronized MerkleRoots and retrieving the last evaluation key from the database.
type MerkleRootsRepository interface {
	// GetLastMerkleRoot should return the merkle root with the heighest height from your storage or undefined if empty
	GetLastMerkleRoot() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are ordered with ascending order by block height
	SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error
}

// SyncMerkleRoots syncs merkleroots known to spv-wallet with the client database
// If timeout is needed pass context.WithTimeout() as ctx param
func (wc *WalletClient) SyncMerkleRoots(ctx context.Context, repo MerkleRootsRepository) error {
	lastEvaluatedKey := repo.GetLastMerkleRoot()
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

			var merkleRootsResponse models.ExclusiveStartKeyPage[[]models.MerkleRoot]

			err := wc.doHTTPRequest(ctx, http.MethodGet, url, nil, wc.xPriv, true, &merkleRootsResponse)

			if err != nil {
				// In case if the context deadline exceeds its limit during http request, httpClient
				// cancels the request wrapping it as spverror, so we need to check if the message
				// is the same as context deadline exceeded error
				if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
					return ErrSyncMerkleRootsTimeout
				}
				return WrapError(err)
			}

			lastEvaluatedKey = merkleRootsResponse.Page.LastEvaluatedKey
			if lastEvaluatedKey != "" && previousLastEvaluatedKey == lastEvaluatedKey {
				return ErrStaleLastEvaluatedKey
			}

			err = repo.SaveMerkleRoots(merkleRootsResponse.Content)
			if err != nil {
				return err
			}

			previousLastEvaluatedKey = lastEvaluatedKey
			if previousLastEvaluatedKey == "" {
				return nil
			}

			lastEvaluatedKeyQuery = fmt.Sprintf("?lastEvaluatedKey=%s", previousLastEvaluatedKey)
		}
	}
}
