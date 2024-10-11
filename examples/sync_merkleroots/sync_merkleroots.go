/*
Package main - sync_merkleroots example
*/
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/models"
)

// simulate a storage of merkle roots that exists on a client side that is using SyncMerkleRoots method
type db struct {
	MerkleRoots []models.MerkleRoot
}

func (db *db) SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error {
	fmt.Print("\nSaveMerkleRoots called\n")
	db.MerkleRoots = append(db.MerkleRoots, syncedMerkleRoots...)
	return nil
}

func (db *db) GetLastMerkleRoot() string {
	if len(db.MerkleRoots) == 0 {
		return ""
	}
	return db.MerkleRoots[len(db.MerkleRoots)-1].MerkleRoot
}

// initalize the storage that exists on a client side
var repository = &db{
	MerkleRoots: []models.MerkleRoot{
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
	},
}

func getLastFiveOrFewer(merkleroots []models.MerkleRoot) []models.MerkleRoot {
	startIndex := len(merkleroots) - 5
	if startIndex < 0 {
		startIndex = 0
	}

	return merkleroots[startIndex:]
}

func main() {
	defer examples.HandlePanic()

	server := "http://localhost:3003/api/v1"

	client, err := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	if err != nil {
		fmt.Println("Error: ", err)
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	fmt.Printf("\n\n Initial State Length: \n %d\n\n", len(repository.MerkleRoots))
	fmt.Printf("\n\nInitial State Last 5 MerkleRoots (or fewer):\n%+v\n", getLastFiveOrFewer(repository.MerkleRoots))

	err = client.SyncMerkleRoots(ctx, repository)
	if err != nil {
		fmt.Println("Error: ", err)
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	fmt.Printf("\n\n After Sync State Length: \n %d\n\n", len(repository.MerkleRoots))
	fmt.Printf("\n\n After Sync State Last 5 MerkleRoots (or fewer):\n%+v\n", getLastFiveOrFewer(repository.MerkleRoots))
}
