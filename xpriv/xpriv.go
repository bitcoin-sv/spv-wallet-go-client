package xpriv

import (
	"fmt"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// Keys is a struct containing the xpriv, xpub and mnemonic
type Keys struct {
	xpriv    string
	xpub     PublicKey
	mnemonic string
}

// PublicKey is a struct containing public key information
type PublicKey struct {
	Xpub string
}

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
	return k.Xpub
}

// Generate generates a random set of keys - xpriv, xpb and mnemonic
func Generate() (KeyWithMnemonic, error) {
	entropy, err := bip39.GenerateEntropy(160)
	if err != nil {
		return nil, fmt.Errorf("generate method: key generation error when creating entropy: %w", err)
	}

	mnemonic, seed, err := bip39.Mnemonic(entropy, "")

	if err != nil {
		return nil, fmt.Errorf("generate method: key generation error when creating mnemonic: %w", err)
	}

	hdXpriv, hdXpub, err := createXPrivAndXPub(seed)
	if err != nil {
		return nil, err
	}

	keys := &Keys{
		xpriv:    hdXpriv.String(),
		xpub:     PublicKey{hdXpub.String()},
		mnemonic: mnemonic,
	}

	return keys, nil
}

// FromMnemonic generates Keys based on given mnemonic
func FromMnemonic(mnemonic string) (KeyWithMnemonic, error) {
	seed, err := bip39.MnemonicToSeed(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("FromMnemonic method: error when creating seed: %w", err)
	}

	hdXpriv, hdXpub, err := createXPrivAndXPub(seed)
	if err != nil {
		return nil, fmt.Errorf("FromMnemonic method: %w", err)
	}

	keys := &Keys{
		xpriv:    hdXpriv.String(),
		xpub:     PublicKey{hdXpub.String()},
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
		xpub:  PublicKey{hdXpub.String()},
	}

	return keys, nil
}

func createXPrivAndXPub(seed []byte) (*bip32.ExtendedKey, *bip32.ExtendedKey, error) {
	hdXpriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, nil, fmt.Errorf("key generation error when creating hd private key: %w", err)
	}

	hdXpub, err := hdXpriv.Neuter()
	if err != nil {
		return nil, nil, fmt.Errorf("key generation error when creating hd public hey: %w", err)
	}
	return hdXpriv, hdXpub, nil
}
