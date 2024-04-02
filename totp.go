package walletclient

import (
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bip32"
)

func (b *WalletClient) GenerateTotpForContact(contact *models.Contact, validationPeriod, digits uint) (string, error) {
	if b.xPriv == nil {
		return "", errors.New("init client with xPriv first")
	}

	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return "", err
	}

	cXpub, err := bip32.NewKeyFromString(contact.PubKey)
	if err != nil {
		return "", fmt.Errorf("contact's PubKey is invalid: %w", err)
	}

	ts := &totp.Service{
		Period: validationPeriod,
		Digits: digits,
	}

	return ts.GenarateTotp(xpriv, cXpub)
}

func (b *WalletClient) ConfirmTotpForContact(contact *models.Contact, passcode string, validationPeriod, digits uint) (bool, error) {
	if b.xPriv == nil {
		return false, errors.New("init client with xPriv first")
	}

	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return false, err
	}

	cXpub, err := bip32.NewKeyFromString(contact.PubKey)
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
	// derive xpriv for current PKI
	// TODO: currently we don't support PKI rotation
	// PKI derivation path: m/0/0/0
	pkiXpriv, err := bitcoin.GetHDKeyByPath(xpriv, utils.ChainExternal, 0)
	if err != nil {
		return nil, err
	}

	return pkiXpriv.Child(0)
}
