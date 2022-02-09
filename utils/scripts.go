package utils

import (
	"github.com/BuxOrg/bux/utils"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
)

// Hash will generate a hash of the given string (used for xPub:hash)
func Hash(data string) string {
	return utils.Hash(data)
}

// RandomHex returns a random hex string and error
func RandomHex(n int) (string, error) {
	return utils.RandomHex(n)
}

// ValidateXPub will check the xPub key for length & validation
func ValidateXPub(rawKey string) (*bip32.ExtendedKey, error) {
	return utils.ValidateXPub(rawKey)
}

// DeriveAddresses will derive the internal and external address from a key
func DeriveAddresses(hdKey *bip32.ExtendedKey, num uint32) (external, internal string, err error) {
	return utils.DeriveAddresses(hdKey, num)
}

// DerivePublicKey will derive the internal and external address from a key
func DerivePublicKey(hdKey *bip32.ExtendedKey, chain uint32, num uint32) (*bec.PublicKey, error) {
	return utils.DerivePublicKey(hdKey, chain, num)
}

// StringInSlice check whether the string already is in the slice
func StringInSlice(a string, list []string) bool {
	return utils.StringInSlice(a, list)
}

// GetTransactionIDFromHex get the transaction ID from the given transaction hex
func GetTransactionIDFromHex(hex string) (string, error) {
	return utils.GetTransactionIDFromHex(hex)
}

// GetUnlockingScript will generate an unlocking script
func GetUnlockingScript(tx *bt.Tx, inputIndex uint32, privateKey *bec.PrivateKey) (*bscript.Script, error) {
	shf := sighash.AllForkID

	sh, err := tx.CalcInputSignatureHash(inputIndex, shf)
	if err != nil {
		return nil, err
	}

	var sig *bec.Signature
	sig, err = privateKey.Sign(bt.ReverseBytes(sh))
	if err != nil {
		return nil, err
	}

	pubKey := privateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	var s *bscript.Script
	s, err = bscript.NewP2PKHUnlockingScript(pubKey, signature, shf)
	if err != nil {
		return nil, err
	}

	return s, nil
}
