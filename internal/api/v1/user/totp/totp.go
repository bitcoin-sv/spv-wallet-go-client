package totp

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	utils "github.com/bitcoin-sv/spv-wallet-go-client/internal/cryptoutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	// DefaultPeriod - Default number of seconds a TOTP is valid for.
	DefaultPeriod uint = 30
	// DefaultDigits - Default TOTP length
	DefaultDigits uint = 2
)

// Client handles TOTP generation and validation.
type Client struct {
	xPriv *bip32.ExtendedKey
}

// New creates a new TOTP WalletClient.
func New(xPriv *bip32.ExtendedKey) *Client {
	return &Client{xPriv: xPriv}
}

// GenerateTotpForContact generates a time-based one-time password (TOTP) for a contact.
func (b *Client) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	sharedSecret, err := b.makeSharedSecret(contact)
	if err != nil {
		return "", fmt.Errorf("generateTotpForContact: error when making shared: %w", err)
	}

	opts := getTotpOpts(period, digits)
	passcode, err := totp.GenerateCodeCustom(directedSecret(sharedSecret, contact.Paymail), time.Now(), *opts)
	if err != nil {
		return "", fmt.Errorf("generateTotpForContact: error when generating TOTP: %w", err)
	}
	return passcode, nil
}

// ValidateTotpForContact validates a TOTP for a contact.
func (b *Client) ValidateTotpForContact(contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	sharedSecret, err := b.makeSharedSecret(contact)
	if err != nil {
		return fmt.Errorf("ValidateTotpForContact: error when making shared secret: %w", err)
	}

	opts := getTotpOpts(period, digits)
	valid, err := totp.ValidateCustom(passcode, directedSecret(sharedSecret, requesterPaymail), time.Now(), *opts)
	if err != nil {
		return fmt.Errorf("ValidateTotpForContact: error when validating TOTP: %w", err)
	}
	if !valid {
		return fmt.Errorf("ValidateTotpForContact: TOTP is invalid")
	}
	return nil
}

func (b *Client) makeSharedSecret(contact *models.Contact) ([]byte, error) {
	privKey, pubKey, err := b.getSharedSecretFactors(contact)
	if err != nil {
		return nil, fmt.Errorf("makeSharedSecret: error when getting shared secret factors: %w", err)
	}

	x, _ := ec.S256().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	return x.Bytes(), nil
}

func (b *Client) getSharedSecretFactors(contact *models.Contact) (*ec.PrivateKey, *ec.PublicKey, error) {
	if b.xPriv == nil {
		return nil, nil, errors.ErrMissingXpriv
	}

	// Derive private key from xPriv for PKI operations.
	xpriv, err := deriveXprivForPki(b.xPriv)
	if err != nil {
		return nil, nil, fmt.Errorf("getSharedSecretFactors: error when deriving xpriv for PKI: %w", err)
	}

	privKey, err := xpriv.ECPrivKey()
	if err != nil {
		return nil, nil, fmt.Errorf("getSharedSecretFactors: error when deriving private key: %w", err)
	}

	// Convert contact's public key.
	pubKey, err := convertPubKey(contact.PubKey)
	if err != nil {
		return nil, nil, errors.ErrContactPubKeyInvalid
	}

	return privKey, pubKey, nil
}

func deriveXprivForPki(xpriv *bip32.ExtendedKey) (*bip32.ExtendedKey, error) {
	pkiXpriv, err := bip32.GetHDKeyByPath(xpriv, utils.ChainExternal, 0)
	if err != nil {
		return nil, fmt.Errorf("deriveXprivForPki: error when deriving xpriv for PKI: %w", err)
	}
	pki, err := pkiXpriv.Child(0)
	if err != nil {
		return nil, fmt.Errorf("deriveXprivForPki: error when deriving xpriv for PKI: %w", err)
	}
	return pki, nil
}

func convertPubKey(pubKey string) (*ec.PublicKey, error) {
	decoded, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, fmt.Errorf("convertPubKey: error when decoding public key: %w", err)
	}

	parsedPubKey, err := ec.ParsePubKey(decoded)
	if err != nil {
		return nil, fmt.Errorf("convertPubKey: error when parsing public key: %w", err)
	}
	return parsedPubKey, nil
}

func getTotpOpts(period, digits uint) *totp.ValidateOpts {
	if period == 0 {
		period = DefaultPeriod
	}

	if digits == 0 {
		digits = DefaultDigits
	}

	return &totp.ValidateOpts{
		Period: period,
		Digits: otp.Digits(digits), //nolint: gosec
	}
}

// directedSecret appends a paymail to the shared secret and encodes it as base32.
func directedSecret(sharedSecret []byte, paymail string) string {
	return base32.StdEncoding.EncodeToString(append(sharedSecret, []byte(paymail)...))
}
