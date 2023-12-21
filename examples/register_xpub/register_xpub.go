package main

import (
	"context"
	"github.com/BuxOrg/go-buxclient/xpriv"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/logger"
)

func main() {
	log := logger.Get()

	// Generate keys
	keys, resErr := xpriv.Generate()
	if resErr != nil {
		log.Fatal().Stack().Msg(resErr.Error())
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

	if err = buxClient.NewXpub(
		context.Background(), keys.XPub().String(), &buxmodels.Metadata{"example_field": "example_data"},
	); err != nil {
		log.Fatal().Stack().Msg(err.Error())
	}

	log.Printf("registered xPub: %s", keys.XPub().String())
}
