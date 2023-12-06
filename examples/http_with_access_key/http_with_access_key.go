package main

import (
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/logger"
)

func main() {
	log := logger.Get()

	//Replace with created access key
	exampleAccessKey := "some_generated_access_key"

	// Create a client
	client, err := buxclient.New(
		buxclient.WithAccessKey(exampleAccessKey),
		buxclient.WithHTTP("http://localhost:3003/v1"),
		buxclient.WithDebugging(true),
		buxclient.WithSignRequest(true),
	)
	if err != nil {
		log.Fatal().Stack().Msg(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", client.IsDebug())
}
