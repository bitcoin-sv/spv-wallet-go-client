package xkeys

import (
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// Keys is a struct containing the xpriv, xpub and mnemonic
type Keys struct {
	Xpriv    *bip32.ExtendedKey
	Xpub     *bip32.ExtendedKey
	Mnemonic string
}

// Generate generates a random set of keys - xpriv, xpb and mnemonic
func Generate() (*Keys, transports.ResponseError) {
	entropy, err := bip39.GenerateEntropy(160)
	if err != nil {
		return nil, transports.WrapError(err)
	}

	mnemonic, seed, err := bip39.Mnemonic(entropy, "")

	if err != nil {
		return nil, transports.WrapError(err)
	}

	hdXpriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)

	if err != nil {
		return nil, transports.WrapError(err)
	}

	hdXpub, err := hdXpriv.Neuter()
	if err != nil {
		return nil, transports.WrapError(err)
	}

	keys := &Keys{
		Xpriv:    hdXpriv,
		Xpub:     hdXpub,
		Mnemonic: mnemonic,
	}

	return keys, nil
}

// GetPublicKeyFromHDPrivateKey returns the public key from the HD private key
func GetPublicKeyFromHDPrivateKey(hdXpriv string) (*bip32.ExtendedKey, transports.ResponseError) {
	hdKey, err := bip32.NewKeyFromString(hdXpriv)
	hdXpub, err := hdKey.Neuter()
	if err != nil {
		return nil, transports.WrapError(err)
	}
	return hdXpub, nil
}
