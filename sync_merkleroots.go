package walletclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/models"
)

// SyncMerkleRoots syncs merkleroots known to spv-wallet with the client database
func (wc *WalletClient) SyncMerkleRoots(ctx context.Context, repo models.MerkleRootsRepository, timeoutMs time.Duration) error {
	var cancel context.CancelFunc
	if timeoutMs > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeoutMs)
		defer cancel()
	}

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
