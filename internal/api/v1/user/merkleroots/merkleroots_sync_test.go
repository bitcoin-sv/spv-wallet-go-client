package merkleroots_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestSyncMerkleRoots(t *testing.T) {
	t.Run("Should properly sync database when empty", func(t *testing.T) {
		// setup
		server := merklerootstest.MockMerkleRootsAPIResponseNormal()
		defer server.Close()

		apiURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		// given
		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
		client := merkleroots.NewAPI(apiURL, resty.New())

		// when
		err = client.SyncMerkleRoots(context.Background(), repo)

		// then
		require.NoError(t, err)
		require.Len(t, repo.MerkleRoots, len(merklerootstest.MockedSPVWalletData))
		require.Equal(t, merklerootstest.LastMockedMerkleRoot(), repo.MerkleRoots[len(repo.MerkleRoots)-1])
	})

	t.Run("Should properly sync database when partially filled", func(t *testing.T) {
		// setup
		server := merklerootstest.MockMerkleRootsAPIResponseNormal()
		defer server.Close()

		apiURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		// given
		client := merkleroots.NewAPI(apiURL, resty.New())
		require.NoError(t, err)

		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})

		// when
		err = client.SyncMerkleRoots(context.Background(), repo)

		// then
		require.NoError(t, err)
		require.Len(t, repo.MerkleRoots, len(merklerootstest.MockedSPVWalletData))
		require.Equal(t, merklerootstest.LastMockedMerkleRoot(), repo.MerkleRoots[len(repo.MerkleRoots)-1])
	})

	t.Run("Should fail sync merkleroots due to the timeout", func(t *testing.T) {
		// setup
		server := merklerootstest.MockMerkleRootsAPIResponseDelayed()
		defer server.Close()

		apiURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		// given
		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
		defer cancel()

		client := merkleroots.NewAPI(apiURL, resty.New())
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(ctx, repo)

		// then
		require.ErrorIs(t, err, goclienterr.ErrSyncMerkleRootsTimeout)
	})

	t.Run("Should fail sync merkleroots due to last evaluated key being the same in the response", func(t *testing.T) {
		// setup
		server := merklerootstest.MockMerkleRootsAPIResponseStale()
		defer server.Close()

		apiURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		// given
		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
		client := merkleroots.NewAPI(apiURL, resty.New())
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(context.Background(), repo)

		// then
		require.ErrorIs(t, err, errors.ErrStaleLastEvaluatedKey)
	})
}
