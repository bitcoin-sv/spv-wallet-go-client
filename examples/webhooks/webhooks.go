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

	client := walletclient.NewWithAdminKey("http://localhost:3003/v1", examples.ExampleAdminKey)
	wh := walletclient.NewWebhook(client, "http://localhost:5005/notification", "", "")
	err := wh.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}

	http.Handle("/notification", wh.HTTPHandler())

	go func() {
		for {
			select {
			case event := <-wh.Channel:
				time.Sleep(100 * time.Millisecond) // simulate processing time
				fmt.Println(event)
			case <-context.Background().Done():
				return
			}
		}
	}()

	http.ListenAndServe(":5005", nil)
}
