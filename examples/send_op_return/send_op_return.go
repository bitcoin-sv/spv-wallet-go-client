package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	metadata := map[string]any{}

	opReturn := response.OpReturn{StringParts: []string{"hello", "world"}}
	draftTransactionCmd := commands.DraftTransaction{
		Config: response.TransactionConfig{
			Outputs: []*response.TransactionOutput{
				{
					OpReturn: &opReturn,
				},
			},
		},
		Metadata: metadata,
	}

	draftTransaction, err := usersAPI.DraftTransaction(ctx, &draftTransactionCmd)
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("DraftTransaction response: ", draftTransaction)

	finalized, err := usersAPI.FinalizeTransaction(draftTransaction)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finalized transaction hex : ", finalized)

	recordTransactionCmd := commands.RecordTransaction{
		Hex:         finalized,
		Metadata:    metadata,
		ReferenceID: draftTransaction.ID,
	}
	transaction, err := usersAPI.RecordTransaction(ctx, &recordTransactionCmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction with OP_RETURN: ", transaction)

	transactionG, err := usersAPI.Transaction(context.Background(), transaction.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction: ", transactionG)
}
