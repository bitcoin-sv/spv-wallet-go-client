package walletclient

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

var ErrClientInitNoXpriv = errors.New("init client with xPriv first")

// GenerateTotpForContact creates one time-based one-time password based on secret shared between the user and the contact
func (b *WalletClient) GenerateTotpForContact(contact *models.Contact, validationPeriod, digits uint) (string, error) {
	if b.xPriv == nil {
		return "", ErrClientInitNoXpriv
	}

	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return "", err
	}

	cXpub, err := convertPubKey(contact.PubKey)
	if err != nil {
		return "", fmt.Errorf("contact's PubKey is invalid: %w", err)
	}

	ts := &totp.Service{
		Period: validationPeriod,
		Digits: digits,
	}

	return ts.GenerateTotp(xpriv, cXpub)
}

// ValidateTotpForContact validates one time-based one-time password based on secret shared between the user and the contact
func (b *WalletClient) ValidateTotpForContact(contact *models.Contact, passcode string, validationPeriod, digits uint) (bool, error) {
	if b.xPriv == nil {
		return false, ErrClientInitNoXpriv
	}

	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return false, err
	}

	cXpub, err := convertPubKey(contact.PubKey)
	if err != nil {
		return false, fmt.Errorf("contact's PubKey is invalid: %w", err)
	}

	ts := &totp.Service{
		Period: validationPeriod,
		Digits: digits,
	}

	return ts.ValidateTotp(xpriv, cXpub, passcode)
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
