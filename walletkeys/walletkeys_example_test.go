package walletkeys_test

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func ExampleRandomKeysWithMnemonic() {
	keys, err := walletkeys.RandomKeysWithMnemonic()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mnemonic: ", keys.Mnemonic())
	fmt.Println("xPriv: ", keys.Keys.XPriv())
	fmt.Println("XPub: ", keys.Keys.XPub())
}

func ExampleXPrivFromString() {
	key := "xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si"
	xPriv, err := walletkeys.XPrivFromString(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("xPriv:", xPriv)

	// Output:
	// xPriv: xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si
}

func ExampleXPrivFromMnemonic() {
	mnemonic := "absorb corn ostrich order sing boost just harvest enable make detail future desert bus adult"
	xPriv, err := walletkeys.XPrivFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("xPriv:", xPriv)

	// Output:
	// xPriv: xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si
}
