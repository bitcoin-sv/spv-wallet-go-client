package walletclient

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var ErrClientInitNoXpriv = errors.New("init client with xPriv first")

const (
	// Default number of seconds a TOTP is valid for.
	TotpDefaultPeriod uint = 30
	// Default TOTP length
	TotpDefaultDigits uint = 2
)

// GenerateTotpForContact creates one time-based one-time password based on secret shared between the user and the contact
func (b *WalletClient) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	secret, err := sharedSecret(b, contact)
	if err != nil {
		return "", err
	}

	opts := getTotpOpts(period, digits)
	return totp.GenerateCodeCustom(string(secret), time.Now(), *opts)
}

// ValidateTotpForContact validates one time-based one-time password based on secret shared between the user and the contact
func (b *WalletClient) ValidateTotpForContact(contact *models.Contact, passcode string, period, digits uint) (bool, error) {
	secret, err := sharedSecret(b, contact)
	if err != nil {
		return false, err
	}

	opts := getTotpOpts(period, digits)
	return totp.ValidateCustom(passcode, string(secret), time.Now(), *opts)
}

func sharedSecret(b *WalletClient, c *models.Contact) (string, error) {
	privKey, pubKey, err := getSharedSecretFactors(b, c)
	if err != nil {
		return "", err
	}

	x, _ := bec.S256().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	return base32.StdEncoding.EncodeToString(x.Bytes()), nil
}

func getTotpOpts(period, digits uint) *totp.ValidateOpts {
	if period == 0 {
		period = TotpDefaultPeriod
	}

	if digits == 0 {
		digits = TotpDefaultDigits
	}

	return &totp.ValidateOpts{
		Period: period,
		Digits: otp.Digits(digits),
	}
}

func getSharedSecretFactors(b *WalletClient, c *models.Contact) (*bec.PrivateKey, *bec.PublicKey, error) {
	if b.xPriv == nil {
		return nil, nil, ErrClientInitNoXpriv
	}

	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := xpriv.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := convertPubKey(c.PubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("contact's PubKey is invalid: %w", err)
	}

	return privKey, pubKey, nil
}

func deriveXprivForPki(xpriv *bip32.ExtendedKey) (*bip32.ExtendedKey, error) {
	// PKI derivation path: m/0/0/0
	// NOTICE: we currently do not support PKI rotation; however, adjustments will be made if and when we decide to implement it

	pkiXpriv, err := bitcoin.GetHDKeyByPath(xpriv, utils.ChainExternal, 0)
	if err != nil {
		return nil, err
	}

	return pkiXpriv.Child(0)
}

func convertPubKey(pubKey string) (*bec.PublicKey, error) {
	hex, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	return bec.ParsePubKey(hex, bec.S256())
}
