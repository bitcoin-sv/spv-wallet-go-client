package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	var db db
	exampleutil.PrettyPrint("Merkle roots in db before sync", db.roots)

	ctx := context.Background()
	err = usersAPI.SyncMerkleRoots(ctx, &db)
	if err != nil {
		log.Fatalf("Failed to sync merkle roots: %v", err)
	}

	exampleutil.PrettyPrint("Merkle roots in db after sync", db.roots)
}

type db struct {
	roots []models.MerkleRoot
}

func (d *db) GetLastMerkleRoot() string {
	if len(d.roots) == 0 {
		return ""
	}
	return d.roots[len(d.roots)-1].MerkleRoot
}

func (d *db) SaveMerkleRoots(roots []models.MerkleRoot) error {
	d.roots = append(d.roots, roots...)
	return nil
}
