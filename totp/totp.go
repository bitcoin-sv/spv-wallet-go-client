package totp

import (
	"encoding/base32"
	"time"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// Time-base one-time password service
type Service struct {
	Period uint
	Digits uint
}

// GenerateTotp creates one time-based one-time password based on secrets calculated from the keys
func (s *Service) GenerateTotp(xPriv *bip32.ExtendedKey, pubKey *bec.PublicKey) (string, error) {
	secret, err := sharedSecret(xPriv, pubKey)
	if err != nil {
		return "", err
	}

	return totp.GenerateCodeCustom(string(secret), time.Now(), s.getOpts())
}

// ValidateTotp checks if given one-time password is valid
func (s *Service) ValidateTotp(xPriv *bip32.ExtendedKey, pubKey *bec.PublicKey, passcode string) (bool, error) {
	secret, err := sharedSecret(xPriv, pubKey)
	if err != nil {
		return false, err
	}

	return totp.ValidateCustom(passcode, string(secret), time.Now(), s.getOpts())
}

func (s *Service) getOpts() totp.ValidateOpts {
	return totp.ValidateOpts{
		Period: s.Period,
		Digits: otp.Digits(s.Digits),
	}
}

func sharedSecret(xPriv *bip32.ExtendedKey, pubKey *bec.PublicKey) (string, error) {
	privKey, err := xPriv.ECPrivKey()
	if err != nil {
		return "", err
	}

	x, _ := bec.S256().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())

	return base32.StdEncoding.EncodeToString(x.Bytes()), nil
}
