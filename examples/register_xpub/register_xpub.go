package main

import (
	"context"
	"log"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient"
	"github.com/bitcoinschema/go-bitcoin/v2"
)

func main() {

	// Example xPub
	masterKey, _ := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
	rawXPub, _ := bitcoin.GetExtendedPublicKey(masterKey)

	// Create a client
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(masterKey.String()),
		buxclient.WithHTTP("localhost:3001"),
		buxclient.WithDebugging(true),
		buxclient.WithSignRequest(true),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err = buxClient.NewXpub(
		context.Background(), rawXPub, &bux.Metadata{"example_field": "example_data"},
	); err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("registered xPub: " + rawXPub)
}
