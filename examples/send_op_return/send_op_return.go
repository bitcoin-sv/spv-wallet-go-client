package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	ctx := context.Background()
	draftTransaction, err := usersAPI.DraftTransaction(ctx, &commands.DraftTransaction{
		Config: response.TransactionConfig{
			Outputs: []*response.TransactionOutput{
				{
					OpReturn: &response.OpReturn{StringParts: []string{"hello", "world"}},
				},
			},
		},
		Metadata: queryparams.Metadata{},
	})
	if err != nil {
		log.Fatalf("Failed to create draft transaction: %v", err)
	}
	exampleutil.PrettyPrint("Created DraftTransaction", draftTransaction)

	finalized, err := usersAPI.FinalizeTransaction(draftTransaction)
	if err != nil {
		log.Fatalf("Failed to finalize draft transaction: %v", err)
	}
	fmt.Printf("Finalized draft transaction hex: %s\n", finalized)

	transaction, err := usersAPI.RecordTransaction(ctx, &commands.RecordTransaction{
		Hex:         finalized,
		Metadata:    queryparams.Metadata{},
		ReferenceID: draftTransaction.ID,
	})
	if err != nil {
		log.Fatalf("Failed to record finalized transaction: %v", err)
	}
	exampleutil.PrettyPrint("Recorded transaction with OP_RETURN", transaction)

	transactionG, err := usersAPI.Transaction(context.Background(), transaction.ID)
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.PrettyPrint("Fetched transaction", transactionG)
}
