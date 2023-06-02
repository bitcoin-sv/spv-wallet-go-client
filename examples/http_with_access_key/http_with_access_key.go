package main

import (
	"log"

	"github.com/BuxOrg/go-buxclient"
)

func main() {

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
		log.Fatalln(err.Error())
	}

	log.Printf("client loaded - bux debug: %v", client.IsDebug())
}
