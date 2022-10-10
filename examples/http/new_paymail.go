package main

import (
	"context"
	"log"

	"github.com/bitcoinschema/go-bitcoin/v2"

	"github.com/BuxOrg/go-buxclient"
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

	log.Printf("client loaded - bux debug: %v", buxClient.IsDebug())
	err = buxClient.NewPaymail(context.Background(), rawXPub, "foo@domain.com", "", "Foo", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("paymail added")

}
