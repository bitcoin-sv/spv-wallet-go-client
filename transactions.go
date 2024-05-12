package walletclient

// import (
// 	"context"
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// 

// 	"github.com/bitcoin-sv/spv-wallet/models"

// 	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
// )

// // GetTransaction get a transaction by id
// func (b *WalletClient) GetTransaction(ctx context.Context, txID string) (*models.Transaction, transports.ResponseError) {
// 	return b.transport.GetTransaction(ctx, txID)
// }

// // GetTransactions get all transactions matching search criteria
// func (b *WalletClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
// 	metadata *models.Metadata, queryParams *transports.QueryParams,
// ) ([]*models.Transaction, transports.ResponseError) {
// 	return b.transport.GetTransactions(ctx, conditions, metadata, queryParams)
// }

// // GetTransactionsCount get number of user transactions
// func (b *WalletClient) GetTransactionsCount(ctx context.Context, conditions map[string]interface{},
// 	metadata *models.Metadata,
// ) (int64, transports.ResponseError) {
// 	return b.transport.GetTransactionsCount(ctx, conditions, metadata)
// }

// // DraftToRecipients initialize a new P2PKH draft transaction to a list of recipients
// func (b *WalletClient) DraftToRecipients(ctx context.Context, recipients []*transports.Recipients,
// 	metadata *models.Metadata,
// ) (*models.DraftTransaction, transports.ResponseError) {
// 	return b.transport.DraftToRecipients(ctx, recipients, metadata)
// }

// // DraftTransaction initialize a new draft transaction
// func (b *WalletClient) DraftTransaction(ctx context.Context, transactionConfig *models.TransactionConfig,
// 	metadata *models.Metadata,
// ) (*models.DraftTransaction, transports.ResponseError) {
// 	return b.transport.DraftTransaction(ctx, transactionConfig, metadata)
// }

// // RecordTransaction record a new transaction
// func (b *WalletClient) RecordTransaction(ctx context.Context, hex, draftID string,
// 	metadata *models.Metadata,
// ) (*models.Transaction, transports.ResponseError) {
// 	return b.transport.RecordTransaction(ctx, hex, draftID, metadata)
// }

// // UpdateTransactionMetadata update the metadata of a transaction
// func (b *WalletClient) UpdateTransactionMetadata(ctx context.Context, txID string,
// 	metadata *models.Metadata,
// ) (*models.Transaction, transports.ResponseError) {
// 	return b.transport.UpdateTransactionMetadata(ctx, txID, metadata)
// }

// // FinalizeTransaction will finalize the transaction
// func (b *WalletClient) FinalizeTransaction(draft *models.DraftTransaction) (string, transports.ResponseError) {
// 	res, err := transports.GetSignedHex(draft, b.xPriv)
// 	if err != nil {
// 		return "", transports.WrapError(err)
// 	}

// 	return res, nil
// }

// // SendToRecipients send to recipients
// func (b *WalletClient) SendToRecipients(ctx context.Context, recipients []*transports.Recipients,
// 	metadata *models.Metadata,
// ) (*models.Transaction, transports.ResponseError) {
// 	draft, err := b.DraftToRecipients(ctx, recipients, metadata)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var hex string
// 	if hex, err = b.FinalizeTransaction(draft); err != nil {
// 		return nil, err
// 	}

// 	return b.RecordTransaction(ctx, hex, draft.ID, metadata)
// }
