package spvwallettest

import (
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	spvwallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
)

const TestAPIAddr = "http://localhost:3003"

const (
	UserXPriv         = "xprv9s21ZrQH143K3fqNnUmXmgfT9ToMtiq5cuKsVBG4E5UqVh4psHDY2XKsEfZKuV4FSZcPS9CYgEQiLUpW2xmHqHFyp23SvTkTCE153cCdwaj"
	UserXPub          = "xpub661MyMwAqRbcG9uqtWJY8pcBhVdrJBYvz8FUHZffnR1pNVPyQpXnaKeM5w2FyH5Wwhf5Cf15mFDVRZnuK9sEHDqqd39qWz36UDoobrzLyFM"
	UserPrivAccessKey = "03a446ede05f04fd92d2707599a80b67ad76f63b3958706819c76308bfc7c1143d"
	UserPubAccessKey  = "0239a60e37d62b0217ac86881caba194ab943e18099c080de70c173daf75d917b2"
	PubKey            = "034252e5359a1de3b8ec08e6c29b80594e88fb47e6ae9ce65ee5a94f0d371d2cde"

	AliceXPriv = "xprv9s21ZrQH143K4JFXqGhBzdrthyNFNuHPaMUwvuo8xvpHwWXprNK7T4JPj1w53S1gojQncyj8JhSh8qouYPZpbocsq934cH5G1t1DRBfgbod"
	AliceXPub  = "xpub661MyMwAqRbcGnKzwJECMmodG1CjnN1EwaQYjJCkXGMGpJryPudMzrcsaK6frwUxXqFxRJwPkKvJh6myJEpQPJS9N67jhZWr24biGe277DH"
	BobXPriv   = "xprv9s21ZrQH143K4VneY3UWCF1o5Kk2tmgGrGtMtsrThCTsHsszEZ6H1iP37ZTwuUBvMwudG68SRkcfTjeu8h3rkayfyqkjKAStFBkuNsBnAkS"
	BobXPub    = "xpub661MyMwAqRbcGys7e51WZNxXdMaXJEQ8DVoxhGG5FXzrAgD8n6QXZWhWxrm2yMzH8e9fxV8TYxmkL9sivVEEoPfDpg4u5mrp2VTqvfGT1Us"
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
	t.Helper()
	key, err := ec.PrivateKeyFromHex(UserPrivAccessKey)
	if err != nil {
		t.Fatalf("test helper - ec private key from hex: %s", err)
	}

	return key
}

func GivenSPVUserAPI(t *testing.T) (*spvwallet.UserAPI, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	cfg := config.Config{
		Addr:      TestAPIAddr,
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	spv, err := spvwallet.NewUserAPIWithXPriv(cfg, UserXPriv)
	if err != nil {
		t.Fatalf("test helper - spv wallet client with xpriv: %s", err)
	}

	return spv, transport
}

func GivenSPVAdminAPI(t *testing.T) (*spvwallet.AdminAPI, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	cfg := config.Config{
		Addr:      TestAPIAddr,
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	api, err := spvwallet.NewAdminAPIWithXPriv(cfg, UserXPriv)
	if err != nil {
		t.Fatalf("test helper - admin api with xPub: %s", err)
	}

	return api, transport
}

func MockPKI(t *testing.T, xpub string) string {
	t.Helper()
	xPub, _ := bip32.NewKeyFromString(xpub)
	var err error
	for i := 0; i < 3; i++ { //magicNumberOfInheritance is 3 -> 2+1; 2: because of the way spv-wallet stores xpubs in db; 1: to make a PKI
		xPub, err = xPub.Child(0)
		if err != nil {
			t.Fatalf("test helper - retrieve a derived child extended key at index 0 failed: %s", err)
		}
	}

	pubKey, err := xPub.ECPubKey()
	if err != nil {
		t.Fatalf("test helper - ec public key from xpub: %s", err)
	}

	return hex.EncodeToString(pubKey.SerializeCompressed())
}

func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       models.UnknownErrorCode,
	}
}

func NewInternalServerSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Code:       models.UnknownErrorCode,
	}
}

func ParseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}

func Ptr[T any](value T) *T {
	return &value
}
