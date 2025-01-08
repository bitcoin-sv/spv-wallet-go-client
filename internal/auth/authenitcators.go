package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

type XpubAuthenticator struct {
	hdKey *bip32.ExtendedKey
}

func (x *XpubAuthenticator) Authenticate(r *resty.Request) error {
	xPub, err := bip32.GetExtendedPublicKey(x.hdKey)
	if err != nil {
		return fmt.Errorf("failed to get extended public key: %w", err)
	}

	r.SetHeader(models.AuthHeader, xPub)
	return nil
}

type XprivAuthenticator struct {
	xpubAuth *XpubAuthenticator
	xpriv    *bip32.ExtendedKey
}

func (x *XprivAuthenticator) Authenticate(r *resty.Request) error {
	err := x.xpubAuth.Authenticate(r)
	if err != nil {
		return fmt.Errorf("failed to set xpub header: %w", err)
	}

	body := bodyString(r)
	header := make(http.Header)
	err = setSignature(&header, x.xpriv, body)
	if err != nil {
		return fmt.Errorf("failed to sign request with xpriv: %w", err)
	}

	r.SetHeaderMultiValues(header)
	return nil
}

type AccessKeyAuthenticator struct {
	priv *ec.PrivateKey
	pub  *ec.PublicKey
}

func (a *AccessKeyAuthenticator) Authenticate(r *resty.Request) error {
	r.Header.Set(models.AuthAccessKey, a.pubKeyHex())
	body := bodyString(r)
	sign, err := createSignatureAccessKey(a.privKeyHex(), body)
	if err != nil {
		return fmt.Errorf("failed to sign request with access key: %w", err)
	}

	setSignatureHeaders(&r.Header, sign)
	return nil
}

func (a *AccessKeyAuthenticator) privKeyHex() string {
	return hex.EncodeToString(a.priv.Serialize())
}

func (a *AccessKeyAuthenticator) pubKeyHex() string {
	return hex.EncodeToString(a.pub.Compressed())
}

func bodyString(r *resty.Request) string {
	switch r.Method {
	case http.MethodGet:
		return ""
	}
	return ""
}

func NewXprivAuthenticator(xpriv string) (*XprivAuthenticator, error) {
	if xpriv == "" {
		return nil, goclienterr.ErrEmptyXprivKey
	}

	hdKey, err := bip32.GenerateHDKeyFromString(xpriv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xpriv key: %w", err)
	}

	return &XprivAuthenticator{
		xpriv:    hdKey,
		xpubAuth: &XpubAuthenticator{hdKey: hdKey},
	}, nil
}

func NewAccessKeyAuthenticator(accessKeyHex string) (*AccessKeyAuthenticator, error) {
	if accessKeyHex == "" {
		return nil, goclienterr.ErrEmptyAccessKey
	}

	privKeyBytes, err := hex.DecodeString(accessKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key string: %w", err)
	}

	privKey, pubKey := ec.PrivateKeyFromBytes(privKeyBytes)
	if privKey == nil || pubKey == nil {
		return nil, errors.New("failed to parse private key: key generation resulted in nil")
	}

	return &AccessKeyAuthenticator{
		priv: privKey,
		pub:  pubKey,
	}, nil
}

func NewXpubOnlyAuthenticator(xpub string) (*XpubAuthenticator, error) {
	if xpub == "" {
		return nil, goclienterr.ErrEmptyPubKey
	}

	xpubKey, err := bip32.GetHDKeyFromExtendedPublicKey(xpub)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xpub key: %w", err)
	}

	return &XpubAuthenticator{hdKey: xpubKey}, nil
}
