package main

import (
	"github.com/BuxOrg/go-buxclient/xpriv"
	"log"

	"github.com/BuxOrg/go-buxclient"
)

func main() {

	// Generate keys
	keys, resErr := xpriv.Generate()
	if resErr != nil {
		log.Fatalln(resErr.Error())
	}

	// Create a client
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(keys.XPriv()),
		buxclient.WithGraphQL("localhost:3001"),
		buxclient.WithDebugging(true),
		buxclient.WithSignRequest(true),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", buxClient.IsDebug())
}
