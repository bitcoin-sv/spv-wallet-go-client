package auth

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bsm "github.com/bitcoin-sv/go-sdk/compat/bsm"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/cryptoutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func setSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) error {
	// Create the signature
	authData, err := createSignature(xPriv, bodyString)
	if err != nil {
		return fmt.Errorf("failed to create signature: %w", err)
	}

	// Set the auth header
	header.Set(models.AuthHeader, authData.XPub)
	setSignatureHeaders(header, authData)
	return nil
}

func createSignature(xPriv *bip32.ExtendedKey, bodyString string) (payload *models.AuthPayload, err error) {
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

func getSigningMessage(xPub string, auth *models.AuthPayload) []byte {
	message := fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
	return []byte(message)
}

func setSignatureHeaders(header *http.Header, authData *models.AuthPayload) {
	header.Set(models.AuthHeaderHash, authData.AuthHash)
	header.Set(models.AuthHeaderNonce, authData.AuthNonce)
	header.Set(models.AuthHeaderTime, fmt.Sprintf("%d", authData.AuthTime))
	header.Set(models.AuthSignature, authData.Signature)
}

func createSignatureAccessKey(privateKeyHex, bodyString string) (payload *models.AuthPayload, err error) {
	privateKey, err := ec.PrivateKeyFromHex(privateKeyHex)
	if err != nil {
		return
	}

	publicKey := privateKey.PubKey()

	// Get the AccessKey
	payload = new(models.AuthPayload)
	payload.AccessKey = hex.EncodeToString(publicKey.SerializeCompressed())

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	payload.AuthNonce, err = cryptoutil.RandomHex(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random hexadecimal string: %w", err)
	}

	return createSignatureCommon(payload, bodyString, privateKey)
}
