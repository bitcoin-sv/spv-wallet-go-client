package xpriv_test

import (
	"fmt"

	"github.com/BuxOrg/go-buxclient/xpriv"
)

func ExampleGenerate() {
	keys, _ := xpriv.Generate()

	fmt.Println("xpriv:", keys.XPriv())
	fmt.Println("xpub:", keys.XPub().String())
}

func ExampleFromMnemonic() {
	keys, _ := xpriv.FromMnemonic("absorb corn ostrich order sing boost just harvest enable make detail future desert bus adult")

	fmt.Println("mnemonic:", keys.Mnemonic())
	fmt.Println("xpriv:", keys.XPriv())
	fmt.Println("xpub:", keys.XPub().String())

	// Output:
	// mnemonic: absorb corn ostrich order sing boost just harvest enable make detail future desert bus adult
	// xpriv: xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si
	// xpub: xpub661MyMwAqRbcFpmY3fFdD4V6ueUBTcaCi49XDCPbRTs5XtDomZpzxAS3LUb2hMfUVphDsSPxfjietmsBRFkLDY9Xa3P4jbgNDMnDK3UqJe2
}
