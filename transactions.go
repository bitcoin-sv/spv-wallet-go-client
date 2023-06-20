package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/mrz1836/go-datastore"
)

// GetTransaction get a transaction by id
func (b *BuxClient) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {
	return b.transport.GetTransaction(ctx, txID)
}

// GetTransactions get all transactions matching search criteria
func (b *BuxClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams,
) ([]*bux.Transaction, error) {
	return b.transport.GetTransactions(ctx, conditions, metadata, queryParams)
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
func (b *BuxClient) UpdateTransactionMetadata(ctx context.Context, txID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {
	return b.transport.UpdateTransactionMetadata(ctx, txID, metadata)
}

// FinalizeTransaction will finalize the transaction
func (b *BuxClient) FinalizeTransaction(draft *bux.DraftTransaction) (string, error) {
	return draft.SignInputs(b.xPriv)
}

// SendToRecipients send to recipients
func (b *BuxClient) SendToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *bux.Metadata) (*bux.Transaction, error) {
	draft, err := b.DraftToRecipients(ctx, recipients, metadata)
	if err != nil {
		return nil, err
	} else if draft == nil {
		return nil, bux.ErrDraftNotFound
	}

	var hex string
	if hex, err = b.FinalizeTransaction(draft); err != nil {
		return nil, err
	}

	return b.RecordTransaction(ctx, hex, draft.ID, metadata)
}
