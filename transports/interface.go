package transports

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/libsv/go-bk/bip32"
)

// XpubService is the xPub related requests
type XpubService interface {
	GetXPub(ctx context.Context) (*buxmodels.Xpub, ResponseError)
	NewXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) ResponseError
	RegisterXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) ResponseError
	UpdateXPubMetadata(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Xpub, ResponseError)
}

// AccessKeyService is the access key related requests
type AccessKeyService interface {
	CreateAccessKey(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.AccessKey, ResponseError)
	GetAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, ResponseError)
	GetAccessKeys(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.AccessKey, ResponseError)
	RevokeAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, ResponseError)
}

// DestinationService is the destination related requests
type DestinationService interface {
	GetDestinationByAddress(ctx context.Context, address string) (*buxmodels.Destination, ResponseError)
	GetDestinationByID(ctx context.Context, id string) (*buxmodels.Destination, ResponseError)
	GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*buxmodels.Destination, ResponseError)
	GetDestinations(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.Destination, ResponseError)
	NewDestination(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Destination, ResponseError)
	UpdateDestinationMetadataByAddress(ctx context.Context, lockingScript string, metadata *buxmodels.Metadata) (*buxmodels.Destination, ResponseError)
	UpdateDestinationMetadataByID(ctx context.Context, id string, metadata *buxmodels.Metadata) (*buxmodels.Destination, ResponseError)
	UpdateDestinationMetadataByLockingScript(ctx context.Context, address string, metadata *buxmodels.Metadata) (*buxmodels.Destination, ResponseError)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata *buxmodels.Metadata) (*buxmodels.DraftTransaction, ResponseError)
	DraftTransaction(ctx context.Context, transactionConfig *buxmodels.TransactionConfig, metadata *buxmodels.Metadata) (*buxmodels.DraftTransaction, ResponseError)
	GetTransaction(ctx context.Context, txID string) (*buxmodels.Transaction, ResponseError)
	GetTransactions(ctx context.Context, conditions map[string]interface{}, metadataConditions *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Transaction, ResponseError)
	GetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	RecordTransaction(ctx context.Context, hex, referenceID string, metadata *buxmodels.Metadata) (*buxmodels.Transaction, ResponseError)
	UpdateTransactionMetadata(ctx context.Context, txID string, metadata *buxmodels.Metadata) (*buxmodels.Transaction, ResponseError)
	UnreserveUtxos(ctx context.Context, referenceID string) ResponseError
}

// PaymailService is the paymail related requests
type PaymailService interface {
	NewPaymail(ctx context.Context, rawXpub, paymailAddress, avatar, publicName string, metadata *buxmodels.Metadata) ResponseError
}

// AdminService is the admin related requests
type AdminService interface {
	AdminGetStatus(ctx context.Context) (bool, ResponseError)
	AdminGetStats(ctx context.Context) (*buxmodels.AdminStats, ResponseError)
	AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.AccessKey, ResponseError)
	AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.BlockHeader, ResponseError)
	AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminGetDestinations(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Destination, ResponseError)
	AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminGetPaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, ResponseError)
	AdminGetPaymails(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.PaymailAddress, ResponseError)
	AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminCreatePaymail(ctx context.Context, xPubID string, address string, publicName string, avatar string) (*buxmodels.PaymailAddress, ResponseError)
	AdminDeletePaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, ResponseError)
	AdminGetTransactions(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Transaction, ResponseError)
	AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminGetUtxos(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Utxo, ResponseError)
	AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminGetXPubs(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Xpub, ResponseError)
	AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, ResponseError)
	AdminRecordTransaction(ctx context.Context, hex string) (*buxmodels.Transaction, ResponseError)
}

// TransportService the transport service interface
type TransportService interface {
	AccessKeyService
	AdminService
	DestinationService
	PaymailService
	TransactionService
	XpubService
	Init() error
	IsDebug() bool
	IsSignRequest() bool
	SetAdminKey(adminKey *bip32.ExtendedKey)
	SetDebug(debug bool)
	SetSignRequest(debug bool)
}
