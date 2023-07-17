package buxclient

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	buxerrors "github.com/BuxOrg/bux-models/bux-errors"
	"github.com/BuxOrg/go-buxclient/transports"
)

// GetTransaction get a transaction by id
func (b *BuxClient) GetTransaction(ctx context.Context, txID string) (*buxmodels.Transaction, error) {
	return b.transport.GetTransaction(ctx, txID)
}

// GetTransactions get all transactions matching search criteria
func (b *BuxClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata, queryParams *transports.QueryParams,
) ([]*buxmodels.Transaction, error) {
	return b.transport.GetTransactions(ctx, conditions, metadata, queryParams)
}

// GetTransactionsCount get number of user transactions
func (b *BuxClient) GetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *buxmodels.Metadata,
) (int64, error) {
	return b.transport.GetTransactionsCount(ctx, conditions, metadata)
}

// DraftToRecipients initialize a new P2PKH draft transaction to a list of recipients
func (b *BuxClient) DraftToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *buxmodels.Metadata,
) (*buxmodels.DraftTransaction, error) {
	return b.transport.DraftToRecipients(ctx, recipients, metadata)
}

// DraftTransaction initialize a new draft transaction
func (b *BuxClient) DraftTransaction(ctx context.Context, transactionConfig *buxmodels.TransactionConfig,
	metadata *buxmodels.Metadata,
) (*buxmodels.DraftTransaction, error) {
	return b.transport.DraftTransaction(ctx, transactionConfig, metadata)
}

// RecordTransaction record a new transaction
func (b *BuxClient) RecordTransaction(ctx context.Context, hex, draftID string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Transaction, error) {
	return b.transport.RecordTransaction(ctx, hex, draftID, metadata)
}

// UpdateTransactionMetadata update the metadata of a transaction
func (b *BuxClient) UpdateTransactionMetadata(ctx context.Context, txID string,
	metadata *buxmodels.Metadata,
) (*buxmodels.Transaction, error) {
	return b.transport.UpdateTransactionMetadata(ctx, txID, metadata)
}

// FinalizeTransaction will finalize the transaction
func (b *BuxClient) FinalizeTransaction(draft *buxmodels.DraftTransaction) (string, error) {
	return transports.SignInputs(draft, b.xPriv)
}

// SendToRecipients send to recipients
func (b *BuxClient) SendToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *buxmodels.Metadata,
) (*buxmodels.Transaction, error) {
	draft, err := b.DraftToRecipients(ctx, recipients, metadata)
	if err != nil {
		return nil, err
	} else if draft == nil {
		return nil, buxerrors.ErrDraftNotFound
	}

	var hex string
	if hex, err = b.FinalizeTransaction(draft); err != nil {
		return nil, err
	}

	return b.RecordTransaction(ctx, hex, draft.ID, metadata)
}
