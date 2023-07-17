package transports

import (
	"context"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/libsv/go-bk/bip32"
)

// XpubService is the xPub related requests
type XpubService interface {
	GetXPub(ctx context.Context) (*buxmodels.Xpub, error)
	NewXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) error
	RegisterXpub(ctx context.Context, rawXPub string, metadata *buxmodels.Metadata) error
	UpdateXPubMetadata(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Xpub, error)
}

// AccessKeyService is the access key related requests
type AccessKeyService interface {
	CreateAccessKey(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.AccessKey, error)
	GetAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, error)
	GetAccessKeys(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.AccessKey, error)
	RevokeAccessKey(ctx context.Context, id string) (*buxmodels.AccessKey, error)
}

// DestinationService is the destination related requests
type DestinationService interface {
	GetDestinationByAddress(ctx context.Context, address string) (*buxmodels.Destination, error)
	GetDestinationByID(ctx context.Context, id string) (*buxmodels.Destination, error)
	GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*buxmodels.Destination, error)
	GetDestinations(ctx context.Context, metadataConditions *buxmodels.Metadata) ([]*buxmodels.Destination, error)
	NewDestination(ctx context.Context, metadata *buxmodels.Metadata) (*buxmodels.Destination, error)
	UpdateDestinationMetadataByAddress(ctx context.Context, lockingScript string, metadata *buxmodels.Metadata) (*buxmodels.Destination, error)
	UpdateDestinationMetadataByID(ctx context.Context, id string, metadata *buxmodels.Metadata) (*buxmodels.Destination, error)
	UpdateDestinationMetadataByLockingScript(ctx context.Context, address string, metadata *buxmodels.Metadata) (*buxmodels.Destination, error)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata *buxmodels.Metadata) (*buxmodels.DraftTransaction, error)
	DraftTransaction(ctx context.Context, transactionConfig *buxmodels.TransactionConfig, metadata *buxmodels.Metadata) (*buxmodels.DraftTransaction, error)
	GetTransaction(ctx context.Context, txID string) (*buxmodels.Transaction, error)
	GetTransactions(ctx context.Context, conditions map[string]interface{}, metadataConditions *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Transaction, error)
	GetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	RecordTransaction(ctx context.Context, hex, referenceID string, metadata *buxmodels.Metadata) (*buxmodels.Transaction, error)
	UpdateTransactionMetadata(ctx context.Context, txID string, metadata *buxmodels.Metadata) (*buxmodels.Transaction, error)
}

// PaymailService is the paymail related requests
type PaymailService interface {
	NewPaymail(ctx context.Context, rawXpub, paymailAddress, avatar, publicName string, metadata *buxmodels.Metadata) error
}

// AdminService is the admin related requests
type AdminService interface {
	AdminGetStatus(ctx context.Context) (bool, error)
	AdminGetStats(ctx context.Context) (*buxmodels.AdminStats, error)
	AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.AccessKey, error)
	AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.BlockHeader, error)
	AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminGetDestinations(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Destination, error)
	AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminGetPaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, error)
	AdminGetPaymails(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.PaymailAddress, error)
	AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminCreatePaymail(ctx context.Context, xPubID string, address string, publicName string, avatar string) (*buxmodels.PaymailAddress, error)
	AdminDeletePaymail(ctx context.Context, address string) (*buxmodels.PaymailAddress, error)
	AdminGetTransactions(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Transaction, error)
	AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminGetUtxos(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Utxo, error)
	AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminGetXPubs(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata, queryParams *QueryParams) ([]*buxmodels.Xpub, error)
	AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{}, metadata *buxmodels.Metadata) (int64, error)
	AdminRecordTransaction(ctx context.Context, hex string) (*buxmodels.Transaction, error)
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
