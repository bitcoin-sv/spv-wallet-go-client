/*
Package main - send_op_return example
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/notifications"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	defer examples.HandlePanic()

	examples.CheckIfAdminKeyExists()

	client := walletclient.NewWithAdminKey("http://localhost:3003/v1", examples.ExampleAdminKey)
	wh := notifications.NewWebhook(
		context.Background(),
		client,
		"http://localhost:5005/notification",
		notifications.WithToken("Authorization", "this-is-the-token"),
		notifications.WithProcessors(3),
	)
	err := wh.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}

	http.Handle("/notification", wh.HTTPHandler())

	if err = notifications.RegisterHandler(wh, func(gpe *models.StringEvent) {
		time.Sleep(50 * time.Millisecond) // simulate processing time
		fmt.Printf("Processing event-string: %s\n", gpe.Value)
	}); err != nil {
		panic(err)
	}

	if err = notifications.RegisterHandler(wh, func(gpe *models.TransactionEvent) {
		time.Sleep(50 * time.Millisecond) // simulate processing time
		fmt.Printf("Processing event-transaction: XPubID: %s, TxID: %s, Status: %s\n", gpe.XPubID, gpe.TransactionID, gpe.Status)
	}); err != nil {
		panic(err)
	}

	go func() {
		_ = http.ListenAndServe(":5005", nil)
	}()

	<-time.After(30 * time.Second)

	fmt.Printf("Unsubscribing...\n")
	err = wh.Unsubscribe(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shutting down...\n")
}
