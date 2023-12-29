package main

import (
	"github.com/BuxOrg/go-buxclient"
)

func main() {

	//Replace with created access key
	exampleAccessKey := "some_generated_access_key"

	// Create a client
	_, _ = buxclient.New(
		buxclient.WithAccessKey(exampleAccessKey),
		buxclient.WithHTTP("http://localhost:3003/v1"),
		buxclient.WithSignRequest(true),
	)
}
