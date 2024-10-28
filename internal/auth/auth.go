// Package auth is responsible for handling the authentication of the requests to the server.
// nolint: unused,nolintlint // FIXME: this file will be used soon
package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bsm "github.com/bitcoin-sv/go-sdk/compat/bsm"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/cryptoutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// GetSignedHex will sign all the inputs using the given xPriv key
func GetSignedHex(dt *models.DraftTransaction, xPriv *bip32.ExtendedKey) (string, error) {
	// Create transaction from hex
	tx, err := trx.NewTransactionFromHex(dt.Hex)
	// we need to reset the inputs as we are going to add them via tx.AddInputFrom (ts-sdk method) and then sign
	tx.Inputs = make([]*trx.TransactionInput, 0)
	if err != nil {
		return "", fmt.Errorf("failed to parse hex, %w", err)
	}

	// Enrich inputs
	for _, draftInput := range dt.Configuration.Inputs {
		lockingScript, err := prepareLockingScript(&draftInput.Destination)
		if err != nil {
			return "", fmt.Errorf("failed to prepare locking script, %w", err)
		}
		unlockScript, err := prepareUnlockingScript(xPriv, &draftInput.Destination)
		if err != nil {
			return "", fmt.Errorf("failed to prepare unlocking script, %w", err)
		}
		err = tx.AddInputFrom(draftInput.TransactionID, draftInput.OutputIndex, lockingScript.String(), draftInput.Satoshis, unlockScript)
		if err != nil {
			return "", fmt.Errorf("failed to add inputs to transaction, %w", err)
		}
	}

	err = tx.Sign()
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction, %w", err)
	}
	return tx.String(), nil
}

// SetSignature will set the signature on the header for the request
func setSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) error {
	// Create the signature
	authData, err := createSignature(xPriv, bodyString)
	if err != nil {
		return fmt.Errorf("create signature op failure: %w", err)
	}

	// Set the auth header
	header.Set(models.AuthHeader, authData.XPub)
	setSignatureHeaders(header, authData)
	return nil
}

func prepareLockingScript(dst *models.Destination) (*script.Script, error) {
	lockingScript, err := script.NewFromHex(dst.LockingScript)
	if err != nil {
		return nil, fmt.Errorf("failed to create locking script from hex for destination: %w", err)
	}
	return lockingScript, nil
}

func prepareUnlockingScript(xPriv *bip32.ExtendedKey, dst *models.Destination) (*p2pkh.P2PKH, error) {
	key, err := getDerivedKeyForDestination(xPriv, dst)
	if err != nil {
		return nil, err
	}
	return getUnlockingScript(key)
}

func getDerivedKeyForDestination(xPriv *bip32.ExtendedKey, dst *models.Destination) (*ec.PrivateKey, error) {
	// Derive the child key (m/chain/num)
	derivedKey, err := bip32.GetHDKeyByPath(xPriv, dst.Chain, dst.Num)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key for unlocking input, %w", err)
	}
	// Handle paymail destination derivation if applicable
	if dst.PaymailExternalDerivationNum != nil {
		derivedKey, err = derivedKey.Child(*dst.PaymailExternalDerivationNum)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key for unlocking paymail input, %w", err)
		}
	}
	// Get the private key from the derived key
	priv, err := bip32.GetPrivateKeyFromHDKey(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key for unlocking paymail input, %w", err)
	}
	return priv, nil
}

// Generate unlocking script using private key
func getUnlockingScript(privateKey *ec.PrivateKey) (*p2pkh.P2PKH, error) {
	sigHashFlags := sighash.AllForkID
	unlocked, err := p2pkh.Unlock(privateKey, &sigHashFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to create unlocking script, %w", err)
	}
	return unlocked, nil
}

// createSignature will create a signature for the given key & body contents
func createSignature(xPriv *bip32.ExtendedKey, bodyString string) (payload *models.AuthPayload, err error) {
	// No key?
	if xPriv == nil {
		err = ErrMissingXpriv
		return
	}

	// Get the xPub
	payload = new(models.AuthPayload)
	if payload.XPub, err = bip32.GetExtendedPublicKey(xPriv); err != nil { // Should never error if key is correct
		return
	}

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	if payload.AuthNonce, err = cryptoutil.RandomHex(32); err != nil { // Should never error if key is correct
		return
	}

	// Derive the address for signing
	var key *bip32.ExtendedKey
	if key, err = cryptoutil.DeriveChildKeyFromHex(xPriv, payload.AuthNonce); err != nil {
		return
	}

	var privateKey *ec.PrivateKey
	if privateKey, err = bip32.GetPrivateKeyFromHDKey(key); err != nil {
		return // Should never error if key is correct
	}
	return createSignatureCommon(payload, bodyString, privateKey)
}

// createSignatureCommon will create a signature
func createSignatureCommon(payload *models.AuthPayload, bodyString string, privateKey *ec.PrivateKey) (*models.AuthPayload, error) {
	// Create the auth header hash
	payload.AuthHash = cryptoutil.Hash(bodyString)
	// auth_time is the current time and makes sure a request can not be sent after 30 secs
	payload.AuthTime = time.Now().UnixMilli()

	key := payload.XPub
	if key == "" && payload.AccessKey != "" {
		key = payload.AccessKey
	}
	// Signature, using bitcoin signMessage
	sigBytes, err := bsm.SignMessage(privateKey, getSigningMessage(key, payload))
	if err != nil {
		return nil, fmt.Errorf("failed to sign message, %w", err)
	}

	payload.Signature = base64.StdEncoding.EncodeToString(sigBytes)
	return payload, nil
}

// getSigningMessage will build the signing message byte array
func getSigningMessage(xPub string, auth *models.AuthPayload) []byte {
	message := fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
	return []byte(message)
}

func setSignatureHeaders(header *http.Header, authData *models.AuthPayload) {
	// Create the auth header hash
	header.Set(models.AuthHeaderHash, authData.AuthHash)
	// Set the nonce
	header.Set(models.AuthHeaderNonce, authData.AuthNonce)
	// Set the time
	header.Set(models.AuthHeaderTime, fmt.Sprintf("%d", authData.AuthTime))
	// Set the signature
	header.Set(models.AuthSignature, authData.Signature)
}

// ErrMissingXpriv is when xpriv is missing
var ErrMissingXpriv = errors.New("xpriv is missing")
