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
)

// simulate database
type db struct {
	MerkleRoots []walletclient.MerkleRoot
}

func (db *db) SaveMerkleRoots(syncedMerkleRoots []walletclient.MerkleRoot) error {
	db.MerkleRoots = append(db.MerkleRoots, syncedMerkleRoots...)
	time.Sleep(1 * time.Second)
	return nil
}

func (db *db) GetLastEvaluatedKey() string {
	if len(db.MerkleRoots) == 0 {
		return ""
	}
	return db.MerkleRoots[len(db.MerkleRoots)-1].MerkleRoot
}

// simulate repository
var repository = &db{
	MerkleRoots: []walletclient.MerkleRoot{
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
		{
			MerkleRoot:  "612209eca3ff078e55d1ffe4602e78ac9a53458c7f9196b38c232d8af9ed635d",
			BlockHeight: 864550,
		},
	},
}

func main() {
	defer examples.HandlePanic()

	server := "http://localhost:3003/api/v1"

	// client := walletclient.NewWithAccessKey(server, examples.ExampleAccessKey)
	client := walletclient.NewWithXPriv(server, examples.ExampleXPriv)
	ctx := context.Background()

	fmt.Printf("\n\n Initial State: \n %+v\n\n", repository.MerkleRoots)

	err := client.SyncMerkleRoots(ctx, repository)
	if err != nil {
		fmt.Println("Error: ", err)
		examples.GetFullErrorMessage(err)
		os.Exit(1)
	}

	fmt.Printf("\n\n After Sync State: \n %+v\n\n", repository.MerkleRoots)
}
