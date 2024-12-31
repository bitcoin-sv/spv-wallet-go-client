package merkleroots

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

const (
	route = "api/v1/merkleroots"
	api   = "User Merkle roots API"
)

// MerkleRootsRepository is an interface responsible for storing synchronized MerkleRoots and retrieving the last evaluation key from the database.
type MerkleRootsRepository interface {
	// GetLastMerkleRoot should return the Merkle root with the highest height from your memory, or undefined if empty.
	GetLastMerkleRoot() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are sorted in ascending order by block height.
	SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error
}

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) MerkleRoots(ctx context.Context, merkleRootOpts ...queries.MerkleRootsQueryOption) (*queries.MerkleRootPage, error) {
	var query queries.MerkleRootsQuery
	for _, o := range merkleRootOpts {
		o(&query)
	}

	params := queryparams.NewURLValues()
	params.AddPair("batchSize", query.BatchSize)
	params.AddPair("lastEvaluatedKey", query.LastEvaluatedKey)

	var result queries.MerkleRootPage
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{
		url:        url.JoinPath(route),
		httpClient: httpClient,
	}
}

func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    api,
		Err:    err,
	}
}

func (a *API) SyncMerkleRoots(ctx context.Context, repo MerkleRootsRepository) error {
	lastEvaluatedKey := repo.GetLastMerkleRoot()
	previousLastEvaluatedKey := lastEvaluatedKey

	for {
		select {
		case <-ctx.Done():
			return goclienterr.ErrSyncMerkleRootsTimeout
		default:
			// Query the MerkleRoots API
			result, err := a.MerkleRoots(ctx, queries.MerkleRootsQueryWithLastEvaluatedKey(lastEvaluatedKey))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return goclienterr.ErrSyncMerkleRootsTimeout
				}
				return fmt.Errorf("failed to fetch merkle roots from API: %w", err)
			}

			// Handle empty results
			if len(result.Content) == 0 {
				return nil
			}

			// Update the last evaluated key
			lastEvaluatedKey = result.Page.LastEvaluatedKey
			if lastEvaluatedKey != "" && previousLastEvaluatedKey == lastEvaluatedKey {
				return goclienterr.ErrStaleLastEvaluatedKey
			}

			// Save fetched Merkle roots
			err = repo.SaveMerkleRoots(result.Content)
			if err != nil {
				return fmt.Errorf("failed to save merkle roots: %w", err)
			}

			if lastEvaluatedKey == "" {
				return nil
			}

			previousLastEvaluatedKey = lastEvaluatedKey
		}
	}
}
