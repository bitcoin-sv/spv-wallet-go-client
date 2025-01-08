package transactions

import (
	"errors"
	"fmt"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	walleterrors "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type noopTransactionSigner struct {
}

func (*noopTransactionSigner) TransactionSignedHex(dt *response.DraftTransaction) (string, error) {
	return "", nil
}

type xPrivTransactionSigner struct {
	xPriv *bip32.ExtendedKey
}

func NewXPrivTransactionSigner(xPriv string) (*xPrivTransactionSigner, error) {
	hdKey, err := bip32.GenerateHDKeyFromString(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPriv str: %w", err)
	}

	return &xPrivTransactionSigner{xPriv: hdKey}, nil
}

func (ts *xPrivTransactionSigner) TransactionSignedHex(dt *response.DraftTransaction) (string, error) {
	// Create transaction from hex
	tx, err := trx.NewTransactionFromHex(dt.Hex)
	if err != nil {
		return "", errors.Join(walleterrors.ErrFailedToParseHex, err)
	}
	// we need to reset the inputs as we are going to add them via tx.AddInputFrom (ts-sdk method) and then sign
	tx.Inputs = make([]*trx.TransactionInput, 0)

	// Enrich inputs
	for _, draftInput := range dt.Configuration.Inputs {
		lockingScript, err := script.NewFromHex(draftInput.Destination.LockingScript)
		if err != nil {
			return "", errors.Join(walleterrors.ErrCreateLockingScript, err)
		}

		// prepare unlocking script
		key, err := getDerivedKeyForDestination(ts.xPriv, &draftInput.Destination)
		if err != nil {
			return "", errors.Join(walleterrors.ErrGetDerivedKeyForDestination, err)
		}
		sigHashFlags := sighash.AllForkID
		unlockScript, err := p2pkh.Unlock(key, &sigHashFlags)
		if err != nil {
			return "", errors.Join(walleterrors.ErrCreateUnlockingScript, err)
		}

		err = tx.AddInputFrom(draftInput.TransactionID, draftInput.OutputIndex, lockingScript.String(), draftInput.Satoshis, unlockScript)
		if err != nil {
			return "", errors.Join(walleterrors.ErrAddInputsToTransaction, err)
		}
	}

	err = tx.Sign()
	if err != nil {
		return "", errors.Join(walleterrors.ErrSignTransaction, err)
	}

	return tx.String(), nil
}

func getDerivedKeyForDestination(xPriv *bip32.ExtendedKey, dst *response.Destination) (*ec.PrivateKey, error) {
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
