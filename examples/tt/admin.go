package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

const serverURL = "http://localhost:3003/v1"
const adminKey = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

func main() {
	clientAdmin, err := setupAdminClient()
	if err != nil {
		log.Fatalf("Failed to setup admin client: %v", err)
	}
	log.Println("is sign req: ", clientAdmin.IsSignRequest())

	ctx := context.Background()

	// Example of creating users and interacting with wallet functions
	handleUsers(ctx, clientAdmin)
}

func setupAdminClient() (*walletclient.WalletClient, error) {

	// Set up a client with administrative privileges
	return walletclient.NewWalletClientWithAdminKey(adminKey, serverURL, true)
}

func handleUsers(ctx context.Context, clientAdmin *walletclient.WalletClient) {
	aliceName, bobName := "alice", "bob"
	aliceClient, _, err := createUser(ctx, aliceName, clientAdmin)
	if err != nil {
		log.Fatalf("Failed to create user %s: %v", aliceName, err)
	}
	bobClient, bobPaymail, err := createUser(ctx, bobName, clientAdmin)
	if err != nil {
		log.Fatalf("Failed to create user %s: %v", bobName, err)
	}

	fmt.Println()
	log.Println(" **** Admin Get Contacts ****")
	c, err := clientAdmin.AdminGetContacts(ctx, nil, nil, &walletclient.QueryParams{})
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

func createUser(ctx context.Context, name string, adminClient *walletclient.WalletClient) (*walletclient.WalletClient, string, error) {
	keys, err := xpriv.Generate()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate keys: %v", err)
	}

	timestamp := time.Now().UnixMicro()
	examplePaymail := fmt.Sprintf("contacttest_%d_%s@auggie.4chain.space", timestamp, name)

	metadata := make(models.Metadata)
	metadata["name"] = name

	if err := adminClient.AdminNewXpub(ctx, keys.XPub().String(), &metadata); err != nil {
		return nil, "", fmt.Errorf("failed to create new xpub: %v", err)
	}
	if _, err := adminClient.AdminCreatePaymail(ctx, keys.XPub().String(), examplePaymail, name, ""); err != nil {
		return nil, "", fmt.Errorf("failed to create paymail: %v", err)
	}

	userClient, err := walletclient.NewWalletClientWithXPrivate(keys.XPriv(), serverURL, true)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user client: %v", err)
	}

	return userClient, examplePaymail, nil
}
