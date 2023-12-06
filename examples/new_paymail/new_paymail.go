package main

import (
	"context"
	"github.com/BuxOrg/go-buxclient/logger"
	"github.com/bitcoinschema/go-bitcoin/v2"

	"github.com/BuxOrg/go-buxclient"
)

func main() {
	log := logger.Get()

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
		log.Fatal().Err(err).Str("examples", "new_paymail").Msg(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", buxClient.IsDebug())
	err = buxClient.NewPaymail(context.Background(), rawXPub, "foo@domain.com", "", "Foo", nil)

	if err != nil {
		log.Fatal().Err(err).Str("examples", "new_paymail").Msg(err.Error())
	}
	log.Printf("paymail added")

}
