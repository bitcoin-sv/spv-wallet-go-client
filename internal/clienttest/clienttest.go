package clienttest

import (
	"testing"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/jarcoal/httpmock"
)

const TestAPIAddr = "http://localhost:3003"

const (
	UserXPriv         = "xprv9s21ZrQH143K3fqNnUmXmgfT9ToMtiq5cuKsVBG4E5UqVh4psHDY2XKsEfZKuV4FSZcPS9CYgEQiLUpW2xmHqHFyp23SvTkTCE153cCdwaj"
	UserXPub          = "xpub661MyMwAqRbcG9uqtWJY8pcBhVdrJBYvz8FUHZffnR1pNVPyQpXnaKeM5w2FyH5Wwhf5Cf15mFDVRZnuK9sEHDqqd39qWz36UDoobrzLyFM"
	UserPrivAccessKey = "03a446ede05f04fd92d2707599a80b67ad76f63b3958706819c76308bfc7c1143d"
	UserPubAccessKey  = "0239a60e37d62b0217ac86881caba194ab943e18099c080de70c173daf75d917b2"
)

func ExtendedKey(t *testing.T) *bip32.ExtendedKey {
	t.Helper()
	key, err := bip32.GenerateHDKeyFromString(UserXPriv)
	if err != nil {
		t.Fatalf("test helper - bip32 generate hd key from string: %s", err)
	}

	return key
}

func PrivateKey(t *testing.T) *ec.PrivateKey {
	key, err := ec.PrivateKeyFromHex(UserPrivAccessKey)
	if err != nil {
		t.Fatalf("test helper - ec private key from hex: %s", err)
	}

	return key
}

func GivenSPVWalletClient(t *testing.T) (*client.Client, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	cfg := client.Config{
		Addr:      TestAPIAddr,
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	spv, err := client.NewWithXPriv(cfg, UserXPriv)
	if err != nil {
		t.Fatalf("test helper - spv wallet client with xpriv: %s", err)
	}

	return spv, transport
}
