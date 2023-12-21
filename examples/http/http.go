package main

import (
	"github.com/BuxOrg/go-buxclient/xpriv"

	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/logger"
)

func main() {
	log := logger.Get()

	// Generate keys
	keys, resErr := xpriv.Generate()
	if resErr != nil {
		log.Fatal().Msg(resErr.Error())
	}

	// Create a client
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(keys.XPriv()),
		buxclient.WithHTTP("localhost:3001"),
		buxclient.WithDebugging(true),
		buxclient.WithSignRequest(true),
	)
	if err != nil {
		log.Fatal().Stack().Msg(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", buxClient.IsDebug())
}
