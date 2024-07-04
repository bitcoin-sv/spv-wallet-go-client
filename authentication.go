package walletclient

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
)

// SetSignature will set the signature on the header for the request
func setSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) error {
	// Create the signature
	authData, err := createSignature(xPriv, bodyString)
	if err != nil {
		return WrapError(err)
	}

	// Set the auth header
	header.Set(models.AuthHeader, authData.XPub)

	setSignatureHeaders(header, authData)

	return nil
}

// GetSignedHex will sign all the inputs using the given xPriv key
func GetSignedHex(dt *models.DraftTransaction, xPriv *bip32.ExtendedKey) (signedHex string, err error) {
	var tx *bt.Tx
	if tx, err = bt.NewTxFromString(dt.Hex); err != nil {
		return
	}

	// Enrich inputs
	for index, draftInput := range dt.Configuration.Inputs {
		tx.Inputs[index].PreviousTxSatoshis = draftInput.Satoshis

		dst := draftInput.Destination
		if err = setPreviousTxScript(tx, uint32(index), &dst); err != nil {
			return
		}

		if err = setUnlockingScript(tx, uint32(index), xPriv, &dst); err != nil {
			return
		}
	}

	// Return the signed hex
	signedHex = tx.String()
	return
}

func setPreviousTxScript(tx *bt.Tx, inputIndex uint32, dst *models.Destination) (err error) {
	var ls *bscript.Script
	if ls, err = bscript.NewFromHexString(dst.LockingScript); err != nil {
		return
	}

	tx.Inputs[inputIndex].PreviousTxScript = ls
	return
}

func setUnlockingScript(tx *bt.Tx, inputIndex uint32, xPriv *bip32.ExtendedKey, dst *models.Destination) (err error) {
	var key *bec.PrivateKey
	if key, err = getDerivedKeyForDestination(xPriv, dst); err != nil {
		return
	}

	var s *bscript.Script
	if s, err = getUnlockingScript(tx, inputIndex, key); err != nil {
		return
	}

	tx.Inputs[inputIndex].UnlockingScript = s
	return
}

func getDerivedKeyForDestination(xPriv *bip32.ExtendedKey, dst *models.Destination) (key *bec.PrivateKey, err error) {
	// Derive the child key (m/chain/num)
	var derivedKey *bip32.ExtendedKey
	if derivedKey, err = bitcoin.GetHDKeyByPath(xPriv, dst.Chain, dst.Num); err != nil {
		return
	}

	// Derive key for paymail destination (m/chain/num/paymailNum)
	if dst.PaymailExternalDerivationNum != nil {
		if derivedKey, err = derivedKey.Child(
			*dst.PaymailExternalDerivationNum,
		); err != nil {
			return
		}
	}

	if key, err = bitcoin.GetPrivateKeyFromHDKey(derivedKey); err != nil {
		return
	}

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
func createSignature(xPriv *bip32.ExtendedKey, bodyString string) (payload *models.AuthPayload, err error) {
	// No key?
	if xPriv == nil {
		err = ErrMissingXpriv
		return
	}

	// Get the xPub
	payload = new(models.AuthPayload)
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
func createSignatureCommon(payload *models.AuthPayload, bodyString string, privateKey *bec.PrivateKey) (*models.AuthPayload, error) {
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
func getSigningMessage(xPub string, auth *models.AuthPayload) string {
	return fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
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
