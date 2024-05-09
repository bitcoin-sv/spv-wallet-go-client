package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

const serverURL = "http://localhost:3003/v1"
const adminKey = "you can get admin key from spv-wallet default settings"

func main() {

	// Create a clientAdmin
	clientAdmin, _ := walletclient.New(
		walletclient.WithXPriv(adminKey),
		walletclient.WithAdminKey(adminKey),
		walletclient.WithHTTP(serverURL),
		walletclient.WithSignRequest(true),
	)
	log.Println("is sign req: ", clientAdmin.IsSignRequest())
	ctx := context.Background()

	aliceName := "alice"
	bobName := "bob"
	aliceClient, _ := createUser(ctx, aliceName)
	bobClient, bobPaymail := createUser(ctx, bobName)

	if _, err := aliceClient.UpsertContact(ctx, bobPaymail, bobName, nil); err != nil {
		panic(err)
	}

	fmt.Println()
	log.Println(" **** Admin Get Contacts ****")
	c, err := clientAdmin.AdminGetContacts(ctx, nil, nil, &transports.QueryParams{})
	if err != nil {
		log.Printf("admin error for get contacts - %v", err)
	}
	if len(c) == 0 {
		log.Printf("empty list of contacts")
		return
	}
	log.Printf("admin got contacts: %#v\n", c)

	fmt.Println()
	log.Println(" **** Admin Update Contact ****")
	bobToUpdate := c[0]
	log.Printf("Bob before update [%v]", c[0].FullName)
	updatedContact, err := clientAdmin.AdminUpdateContact(ctx, bobToUpdate.ID, fmt.Sprintf("%s %s", bobName, time.Now().Local().String()), nil)
	if err != nil {
		log.Panicf("error updating contact id: [%s] name: [%s] - %s", bobToUpdate.ID, bobToUpdate.FullName, err)
	}
	log.Printf("updated contact full name [%v]", updatedContact.FullName)

	fmt.Println()
	log.Println(" **** Admin Reject Contact ****")
	aliceContacts, err := bobClient.GetContacts(ctx, nil, nil, nil)
	if err != nil {
		log.Printf("admin error for get contacts - %v", err)
	}
	if len(aliceContacts) == 0 {
		log.Printf("empty list of contacts")
		return
	}
	log.Printf("status should be awaiting == [%v]", aliceContacts[0].Status)

	rejectAliceContact, err := clientAdmin.AdminRejectContact(ctx, aliceContacts[0].ID)
	if err != nil {
		log.Panicf("contact not accepted id: %v - %s", aliceContacts[0].ID, err)
	}
	log.Printf("status should change to rejected == [%v]", rejectAliceContact.Status)

	fmt.Println()
	log.Println(" **** Admin Accept Contact ****")
	if _, err := aliceClient.UpsertContact(ctx, bobPaymail, bobName, nil); err != nil {
		panic(err)
	}

	aliceContacts, err = bobClient.GetContacts(ctx, nil, nil, nil)
	if err != nil {
		log.Printf("admin error for get contacts - %v", err)
	}
	if len(aliceContacts) == 0 {
		log.Printf("empty list of contacts")
		return
	}

	log.Printf("status should be awaiting == [%v]", aliceContacts[0].Status)
	acceptAliceContact, err := clientAdmin.AdminAcceptContact(ctx, aliceContacts[0].ID)
	if err != nil {
		log.Panicf("contact not accepted id: %v - %s", aliceContacts[0].ID, err)
	}
	log.Printf("status should change to unconfirmed == [%v]", acceptAliceContact.Status)

	fmt.Println()
	log.Println(" **** Admin Delete Contact ****")

	err = clientAdmin.AdminDeleteContact(ctx, aliceContacts[0].ID)
	if err != nil {
		log.Panicf("contact not accepted id: %v - %s", aliceContacts[0].ID, err)
	}
	if len(aliceContacts) > 0 {
		log.Println("removed id:", aliceContacts[0].ID)
	} else {
		log.Println("alice contacts are empty")
	}
	aliceContacts, err = bobClient.GetContacts(ctx, nil, nil, nil)
	if err != nil {
		log.Printf("admin error for get contacts - %v", err)
	}
	if len(aliceContacts) == 0 {
		log.Printf("empty list of contacts")
	}
	for _, c := range aliceContacts {
		log.Println("alice contacts id:", c.ID, c.FullName, c.DeletedAt)
	}
}

func createUser(ctx context.Context, name string) (*walletclient.WalletClient, string) {
	keys, _ := xpriv.Generate()

	timestamp := time.Now().UnixMicro()
	examplePaymail := fmt.Sprintf("contacttest_%d_%s@auggie.4chain.space", timestamp, name)

	adminClient, _ := walletclient.New(
		walletclient.WithXPriv(adminKey),
		walletclient.WithAdminKey(adminKey),
		walletclient.WithHTTP(serverURL),
		walletclient.WithSignRequest(true),
	)

	metadata := make(models.Metadata)
	metadata["name"] = name

	err := adminClient.AdminNewXpub(ctx, keys.XPub().String(), &metadata)
	if err != nil {
		panic(err)
	}
	_, err = adminClient.AdminCreatePaymail(context.Background(), keys.XPub().String(), examplePaymail, name, "")
	if err != nil {
		panic(err)
	}

	userClient, _ := walletclient.New(
		walletclient.WithXPriv(keys.XPriv()),
		walletclient.WithHTTP(serverURL),
		walletclient.WithSignRequest(true),
	)

	return userClient, examplePaymail
}
