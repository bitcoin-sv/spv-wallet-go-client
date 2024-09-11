// Package xpriv manges keys
package xpriv

// "github.com/libsv/go-bk/bip39" - no replacements

import (
	"fmt"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bip39 "github.com/bitcoin-sv/go-sdk/compat/bip39"
	chaincfg "github.com/bitcoin-sv/go-sdk/transaction/chaincfg"
)

// Keys is a struct containing the xpriv, xpub and mnemonic
type Keys struct {
	xpriv    string
	xpub     PublicKey
	mnemonic string
}

// PublicKey is a struct containing public key information
type PublicKey string

// Key represents basic key methods
type Key interface {
	XPriv() string
	XPub() PubKey
}

// PubKey represents public key methods
type PubKey interface {
	String() string
}

// KeyWithMnemonic represents methods for generated keys
type KeyWithMnemonic interface {
	Key
	Mnemonic() string
}

// XPub return hierarchical struct which contain xpub info
func (k *Keys) XPub() PubKey {
	return k.xpub
}

// XPriv return hierarchical deterministic private key
func (k *Keys) XPriv() string {
	return k.xpriv
}

// Mnemonic return mnemonic from which keys where generated
func (k *Keys) Mnemonic() string {
	return k.mnemonic
}

// String return hierarchical deterministic publick ey
func (k PublicKey) String() string {
	return string(k)
}

// Generate generates a random set of keys - xpriv, xpb and mnemonic
func Generate() (KeyWithMnemonic, error) {
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		return nil, fmt.Errorf("generate method: key generation error when creating entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)

	if err != nil {
		return nil, fmt.Errorf("generate method: key generation error when creating mnemonic: %w", err)
	}

	hdKey, err := bip32.GenerateHDKeyFromMnemonic(mnemonic, "", &chaincfg.MainNet)
	if err != nil {
		return nil, err
	}

	hdXpriv := hdKey.String()
	hdXpub, err := bip32.GetExtendedPublicKey(hdKey)
	if err != nil {
		return nil, err
	}

	keys := &Keys{
		xpriv:    hdXpriv,
		xpub:     PublicKey(hdXpub),
		mnemonic: mnemonic,
	}

	return keys, nil
}

// FromMnemonic generates Keys based on given mnemonic
func FromMnemonic(mnemonic string) (KeyWithMnemonic, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("FromMnemonic method: error when creating seed: %w", err)
	}

	hdXpriv, hdXpub, err := createXPrivAndXPub(seed)
	if err != nil {
		return nil, fmt.Errorf("FromMnemonic method: %w", err)
	}

	keys := &Keys{
		xpriv:    hdXpriv.String(),
		xpub:     PublicKey(hdXpub.String()),
		mnemonic: mnemonic,
	}

	return keys, nil
}

// FromString generates keys from given xpriv
func FromString(xpriv string) (Key, error) {
	hdXpriv, err := bip32.NewKeyFromString(xpriv)
	if err != nil {
		return nil, fmt.Errorf("FromString method: key generation error when creating hd private key: %w", err)
	}

	hdXpub, err := hdXpriv.Neuter()
	if err != nil {
		return nil, fmt.Errorf("FromString method: key generation error when creating hd public hey: %w", err)
	}

	keys := &Keys{
		xpriv: hdXpriv.String(),
		xpub:  PublicKey(hdXpub.String()),
	}

	return keys, nil
}

func createXPrivAndXPub(seed []byte) (hdXpriv *bip32.ExtendedKey, hdXpub *bip32.ExtendedKey, err error) {
	hdXpriv, err = bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, nil, fmt.Errorf("key generation error when creating hd private key: %w", err)
	}

	hdXpub, err = hdXpriv.Neuter()
	if err != nil {
		return nil, nil, fmt.Errorf("key generation error when creating hd public hey: %w", err)
	}
	return hdXpriv, hdXpub, nil
}
