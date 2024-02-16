package main

import (
	"context"
	"fmt"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	// Replace with your admin keys
	adminXpriv := "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXpub := "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

	// Create a client
	walletClient, _ := walletclient.New(
		walletclient.WithXPriv(adminXpriv),
		walletclient.WithHTTP("localhost:3003/v1"),
		walletclient.WithSignRequest(true),
	)

	ctx := context.Background()

	_ = walletClient.AdminNewXpub(ctx, adminXpub, &models.Metadata{"example_field": "example_data"})

	xpubKey, err := walletClient.GetXPub(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(xpubKey)
}
