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
)

func main() {
	defer examples.HandlePanic()

	client := walletclient.NewWithAdminKey("http://localhost:3003/v1", "xprv9s21ZrQH143K2pmNeAHBzU4JHNDaFaPTbzKbBCw55ErhMDLsxDwKqcaDVV3PwmEmRZa9qUaU261iJaUx8eBiBF77zrPxTH8JGXC7LZQnsgA")
	wh := walletclient.NewWebhook(client, "http://localhost:5005/notification", "Authorization", "this-is-the-token")
	err := wh.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}

	http.Handle("/notification", wh.HTTPHandler())

	_ = walletclient.NewNotificationsDispatcher(context.Background(), wh.Channel, []walletclient.Handler{
		{Model: &walletclient.GeneralPurposeEvent{}, HandlerFunc: func(gpe *walletclient.GeneralPurposeEvent) {
			time.Sleep(50 * time.Millisecond) // simulate processing time
			fmt.Printf("Processing event: %s\n", gpe.Value)
		}},
	})

	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-wh.Channel:
	// 			time.Sleep(50 * time.Millisecond) // simulate processing time
	// 			fmt.Println("Processing event:", event)
	// 		case <-context.Background().Done():
	// 			return
	// 		}
	// 	}
	// }()

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
