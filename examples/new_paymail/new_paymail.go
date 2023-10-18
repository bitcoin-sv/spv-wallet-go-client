package main

import (
	"context"
	"github.com/BuxOrg/go-buxclient/xkeys"
	"log"

	"github.com/BuxOrg/go-buxclient"
)

func main() {

	// Generate keys
	keys, resErr := xkeys.Generate()
	if resErr != nil {
		log.Fatalln(resErr.Error())
	}

	// Create a client
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(keys.Xpriv.String()),
		buxclient.WithHTTP("localhost:3001"),
		buxclient.WithDebugging(true),
		buxclient.WithSignRequest(true),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", buxClient.IsDebug())
	err = buxClient.NewPaymail(context.Background(), keys.Xpub.String(), "foo@domain.com", "", "Foo", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("paymail added")

}
