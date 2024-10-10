package walletclient

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
	"github.com/bitcoin-sv/spv-wallet-go-client/models"
	"github.com/stretchr/testify/require"
)

func TestSyncMerkleRoots(t *testing.T) {

	t.Run("Should properly sync database when empty", func(t *testing.T) {
		// setup
		server := fixtures.MockMerkleRootsAPIResponseNormal()
		defer server.Close()

		// given
		repo := fixtures.CreateRepository([]models.MerkleRoot{})
		client, err := NewWithXPriv(server.URL, fixtures.XPrivString)
		require.NotNil(t, client.xPriv)
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(context.Background(), repo, 0)

		// then
		require.NoError(t, err)
		require.Equal(t, len(fixtures.MockedSPVWalletData), len(repo.MerkleRoots))
		require.Equal(t, fixtures.MockedSPVWalletData[len(fixtures.MockedSPVWalletData)-1].MerkleRoot, repo.MerkleRoots[len(repo.MerkleRoots)-1].MerkleRoot)
		require.Equal(t, fixtures.MockedSPVWalletData[len(fixtures.MockedSPVWalletData)-1].BlockHeight, repo.MerkleRoots[len(repo.MerkleRoots)-1].BlockHeight)
	})

	t.Run("Should properly sync database when partially filled", func(t *testing.T) {
		// setup
		server := fixtures.MockMerkleRootsAPIResponseNormal()
		defer server.Close()

		// given
		repo := fixtures.CreateRepository([]models.MerkleRoot{
			{
				MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
				BlockHeight: 0,
			},
			{
				MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
				BlockHeight: 1,
			},
			{
				MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
				BlockHeight: 2,
			},
		})
		client, err := NewWithXPriv(server.URL, fixtures.XPrivString)
		require.NotNil(t, client.xPriv)
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(context.Background(), repo, 0)

		// then
		require.NoError(t, err)
		require.Equal(t, len(fixtures.MockedSPVWalletData), len(repo.MerkleRoots))
		require.Equal(t, fixtures.MockedSPVWalletData[len(fixtures.MockedSPVWalletData)-1].MerkleRoot, repo.MerkleRoots[len(repo.MerkleRoots)-1].MerkleRoot)
		require.Equal(t, fixtures.MockedSPVWalletData[len(fixtures.MockedSPVWalletData)-1].BlockHeight, repo.MerkleRoots[len(repo.MerkleRoots)-1].BlockHeight)
	})

	t.Run("Should fail sync merkleroots due to the time out", func(t *testing.T) {
		// setup
		server := fixtures.MockMerkleRootsAPIResponseDelayed()
		defer server.Close()

		// given
		repo := fixtures.CreateRepository([]models.MerkleRoot{})
		client, err := NewWithXPriv(server.URL, fixtures.XPrivString)
		require.NotNil(t, client.xPriv)
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(context.Background(), repo, 10)

		// then
		require.ErrorIs(t, err, ErrSyncMerkleRootsTimeout)
	})

	t.Run("Should fail sync merkleroots due to last evaluated key being the same in the response", func(t *testing.T) {
		// setup
		server := fixtures.MockMerkleRootsAPIResponseStale()
		defer server.Close()

		// given
		repo := fixtures.CreateRepository([]models.MerkleRoot{})
		client, err := NewWithXPriv(server.URL, fixtures.XPrivString)
		require.NotNil(t, client.xPriv)
		require.NoError(t, err)

		// when
		err = client.SyncMerkleRoots(context.Background(), repo, 0)

		// then
		require.ErrorIs(t, err, ErrStaleLastEvaluatedKey)
	})
}
