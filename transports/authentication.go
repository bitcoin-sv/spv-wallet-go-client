package transports

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	buxmodels "github.com/BuxOrg/bux-models"
	buxerrors "github.com/BuxOrg/bux-models/bux-errors"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"

	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
)

// SetSignature will set the signature on the header for the request
func setSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) ResponseError {
	// Create the signature
	authData, err := createSignature(xPriv, bodyString)
	if err != nil {
		return WrapError(err)
	}

	// Set the auth header
	header.Set(buxmodels.AuthHeader, authData.XPub)

	return setSignatureHeaders(header, authData)
}

// SignInputs will sign all the inputs using the given xPriv key
func SignInputs(dt *buxmodels.DraftTransaction, xPriv *bip32.ExtendedKey) (signedHex string, resError ResponseError) {
	var err error
	// Start a bt draft transaction
	var txDraft *bt.Tx
	if txDraft, err = bt.NewTxFromString(dt.Hex); err != nil {
		resError = WrapError(err)
		return
	}

	// Sign the inputs
	for index, input := range dt.Configuration.Inputs {

		// Get the locking script
		var ls *bscript.Script
		if ls, err = bscript.NewFromHexString(
			input.Destination.LockingScript,
		); err != nil {
			resError = WrapError(err)
			return
		}
		txDraft.Inputs[index].PreviousTxScript = ls
		txDraft.Inputs[index].PreviousTxSatoshis = input.Satoshis

		// Derive the child key (chain)
		var chainKey *bip32.ExtendedKey
		if chainKey, err = xPriv.Child(
			input.Destination.Chain,
		); err != nil {
			resError = WrapError(err)
			return
		}

		// Derive the child key (num)
		var numKey *bip32.ExtendedKey
		if numKey, err = chainKey.Child(
			input.Destination.Num,
		); err != nil {
			resError = WrapError(err)
			return
		}

		// Get the private key
		var privateKey *bec.PrivateKey
		if privateKey, err = bitcoin.GetPrivateKeyFromHDKey(
			numKey,
		); err != nil {
			resError = WrapError(err)
			return
		}

		// Get the unlocking script
		var s *bscript.Script
		if s, err = getUnlockingScript(
			txDraft, uint32(index), privateKey,
		); err != nil {
			resError = WrapError(err)
			return
		}

		// Insert the locking script
		if err = txDraft.InsertInputUnlockingScript(
			uint32(index), s,
		); err != nil {
			resError = WrapError(err)
			return
		}
	}

	// Return the signed hex
	signedHex = txDraft.String()
	return
}

// GetUnlockingScript will generate an unlocking script
func getUnlockingScript(tx *bt.Tx, inputIndex uint32, privateKey *bec.PrivateKey) (*bscript.Script, error) {
	sigHashFlags := sighash.AllForkID

	sigHash, err := tx.CalcInputSignatureHash(inputIndex, sigHashFlags)
	if err != nil {
		return nil, err
	}

	var sig *bec.Signature
	if sig, err = privateKey.Sign(sigHash); err != nil {
		return nil, err
	}

	pubKey := privateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	var script *bscript.Script
	if script, err = bscript.NewP2PKHUnlockingScript(pubKey, signature, sigHashFlags); err != nil {
		return nil, err
	}

	return script, nil
}

// createSignature will create a signature for the given key & body contents
func createSignature(xPriv *bip32.ExtendedKey, bodyString string) (payload *buxmodels.AuthPayload, err error) {
	// No key?
	if xPriv == nil {
		err = buxerrors.ErrMissingXPriv
		return
	}

	// Get the xPub
	payload = new(buxmodels.AuthPayload)
	if payload.XPub, err = bitcoin.GetExtendedPublicKey(
		xPriv,
	); err != nil { // Should never error if key is correct
		return
	}

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	if payload.AuthNonce, err = utils.RandomHex(32); err != nil { // Should never error if key is correct
		return
	}

	// Derive the address for signing
	var key *bip32.ExtendedKey
	if key, err = utils.DeriveChildKeyFromHex(
		xPriv, payload.AuthNonce,
	); err != nil {
		return
	}

	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.GetPrivateKeyFromHDKey(key); err != nil {
		return // Should never error if key is correct
	}

	return createSignatureCommon(payload, bodyString, privateKey)
}

// createSignatureCommon will create a signature
func createSignatureCommon(payload *buxmodels.AuthPayload, bodyString string, privateKey *bec.PrivateKey) (*buxmodels.AuthPayload, error) {
	// Create the auth header hash
	payload.AuthHash = utils.Hash(bodyString)

	// auth_time is the current time and makes sure a request can not be sent after 30 secs
	payload.AuthTime = time.Now().UnixMilli()

	key := payload.XPub
	if key == "" && payload.AccessKey != "" {
		key = payload.AccessKey
	}

	// Signature, using bitcoin signMessage
	var err error
	if payload.Signature, err = bitcoin.SignMessage(
		hex.EncodeToString(privateKey.Serialise()),
		getSigningMessage(key, payload),
		true,
	); err != nil {
		return nil, err
	}

	return payload, nil
}

// getSigningMessage will build the signing message string
func getSigningMessage(xPub string, auth *buxmodels.AuthPayload) string {
	return fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
}

func setSignatureHeaders(header *http.Header, authData *buxmodels.AuthPayload) ResponseError {
	// Create the auth header hash
	header.Set(buxmodels.AuthHeaderHash, authData.AuthHash)

	// Set the nonce
	header.Set(buxmodels.AuthHeaderNonce, authData.AuthNonce)

	// Set the time
	header.Set(buxmodels.AuthHeaderTime, fmt.Sprintf("%d", authData.AuthTime))

	// Set the signature
	header.Set(buxmodels.AuthSignature, authData.Signature)

	return nil
}
