package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/BuxOrg/go-buxclient/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
)

// GetTransaction get a transaction by id
func (b *BuxClient) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {
	return b.transport.GetTransaction(ctx, txID)
}

// GetTransactions get all transactions matching search criteria
func (b *BuxClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) ([]*bux.Transaction, error) {

	return b.transport.GetTransactions(ctx, conditions, metadata)
}

// DraftToRecipients initialize a new P2PKH draft transaction to a list of recipients
func (b *BuxClient) DraftToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	return b.transport.DraftToRecipients(ctx, recipients, metadata)
}

// DraftTransaction initialize a new draft transaction
func (b *BuxClient) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	return b.transport.DraftTransaction(ctx, transactionConfig, metadata)
}

// RecordTransaction record a new transaction
func (b *BuxClient) RecordTransaction(ctx context.Context, hex, draftID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	return b.transport.RecordTransaction(ctx, hex, draftID, metadata)
}

// UpdateTransactionMetadata update the metadata of a transaction
func (b *BuxClient) UpdateTransactionMetadata(ctx context.Context, txID string, metadata *bux.Metadata) (*bux.Transaction, error) {
	return b.transport.UpdateTransactionMetadata(ctx, txID, metadata)
}

// FinalizeTransaction will finalize the transaction
func (b *BuxClient) FinalizeTransaction(draft *bux.DraftTransaction) (string, error) {
	txDraft, err := bt.NewTxFromString(draft.Hex)
	if err != nil {
		return "", err
	}

	// sign the inputs
	for index, input := range draft.Configuration.Inputs {
		var ls *bscript.Script
		ls, err = bscript.NewFromHexString(input.Destination.LockingScript)
		if err != nil {
			return "", err
		}
		txDraft.Inputs[index].PreviousTxScript = ls

		var chainKey *bip32.ExtendedKey
		chainKey, err = b.xPriv.Child(input.Destination.Chain)
		if err != nil {
			return "", err
		}

		var numKey *bip32.ExtendedKey
		numKey, err = chainKey.Child(input.Destination.Num)
		if err != nil {
			return "", err
		}

		var privateKey *bec.PrivateKey
		privateKey, err = bitcoin.GetPrivateKeyFromHDKey(numKey)
		if err != nil {
			return "", err
		}

		var s *bscript.Script
		s, err = utils.GetUnlockingScript(txDraft, uint32(index), privateKey)
		if err != nil {
			return "", err
		}

		err = txDraft.InsertInputUnlockingScript(uint32(index), s)
		if err != nil {
			return "", err
		}
	}

	return txDraft.String(), nil
}

// SendToRecipients send to recipients
func (b *BuxClient) SendToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	draft, err := b.DraftToRecipients(ctx, recipients, metadata)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		return nil, bux.ErrDraftNotFound
	}

	var hex string
	if hex, err = b.FinalizeTransaction(draft); err != nil {
		return nil, err
	}

	return b.RecordTransaction(ctx, hex, draft.ID, metadata)
}
