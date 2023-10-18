package main

import (
	"context"
	"github.com/BuxOrg/go-buxclient/xkeys"
	"log"

	buxmodels "github.com/BuxOrg/bux-models"
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

	if err = buxClient.NewXpub(
		context.Background(), keys.Xpub.String(), &buxmodels.Metadata{"example_field": "example_data"},
	); err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("registered xPub: " + keys.Xpub.String())
}
