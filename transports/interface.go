package transports

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bip32"
)

// XpubService is the xPub related requests
type XpubService interface {
	GetXPub(ctx context.Context) (*models.Xpub, ResponseError)
	UpdateXPubMetadata(ctx context.Context, metadata *models.Metadata) (*models.Xpub, ResponseError)
}

// AccessKeyService is the access key related requests
type AccessKeyService interface {
	CreateAccessKey(ctx context.Context, metadata *models.Metadata) (*models.AccessKey, ResponseError)
	GetAccessKey(ctx context.Context, id string) (*models.AccessKey, ResponseError)
	GetAccessKeys(ctx context.Context, metadataConditions *models.Metadata) ([]*models.AccessKey, ResponseError)
	GetAccessKeysCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	RevokeAccessKey(ctx context.Context, id string) (*models.AccessKey, ResponseError)
}

// DestinationService is the destination related requests
type DestinationService interface {
	GetDestinationByAddress(ctx context.Context, address string) (*models.Destination, ResponseError)
	GetDestinationByID(ctx context.Context, id string) (*models.Destination, ResponseError)
	GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*models.Destination, ResponseError)
	GetDestinations(ctx context.Context, metadataConditions *models.Metadata) ([]*models.Destination, ResponseError)
	GetDestinationsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	NewDestination(ctx context.Context, metadata *models.Metadata) (*models.Destination, ResponseError)
	UpdateDestinationMetadataByAddress(ctx context.Context, lockingScript string, metadata *models.Metadata) (*models.Destination, ResponseError)
	UpdateDestinationMetadataByID(ctx context.Context, id string, metadata *models.Metadata) (*models.Destination, ResponseError)
	UpdateDestinationMetadataByLockingScript(ctx context.Context, address string, metadata *models.Metadata) (*models.Destination, ResponseError)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata *models.Metadata) (*models.DraftTransaction, ResponseError)
	DraftTransaction(ctx context.Context, transactionConfig *models.TransactionConfig, metadata *models.Metadata) (*models.DraftTransaction, ResponseError)
	GetTransaction(ctx context.Context, txID string) (*models.Transaction, ResponseError)
	GetTransactions(ctx context.Context, conditions map[string]interface{}, metadataConditions *models.Metadata, queryParams *QueryParams) ([]*models.Transaction, ResponseError)
	GetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	RecordTransaction(ctx context.Context, hex, referenceID string, metadata *models.Metadata) (*models.Transaction, ResponseError)
	UpdateTransactionMetadata(ctx context.Context, txID string, metadata *models.Metadata) (*models.Transaction, ResponseError)
	GetUtxo(ctx context.Context, txID string, outputIndex uint32) (*models.Utxo, ResponseError)
	GetUtxos(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Utxo, ResponseError)
	GetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
}

// ContactService is the contact related requests
type ContactService interface {
	AcceptContact(ctx context.Context, paymail string) ResponseError
	RejectContact(ctx context.Context, paymail string) ResponseError
}

// AdminService is the admin related requests
type AdminService interface {
	AdminGetStatus(ctx context.Context) (bool, ResponseError)
	AdminGetStats(ctx context.Context) (*models.AdminStats, ResponseError)
	AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.AccessKey, ResponseError)
	AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.BlockHeader, ResponseError)
	AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminGetDestinations(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Destination, ResponseError)
	AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminGetPaymail(ctx context.Context, address string) (*models.PaymailAddress, ResponseError)
	AdminGetPaymails(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.PaymailAddress, ResponseError)
	AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminCreatePaymail(ctx context.Context, rawXPub string, address string, publicName string, avatar string) (*models.PaymailAddress, ResponseError)
	AdminDeletePaymail(ctx context.Context, address string) ResponseError
	AdminGetTransactions(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Transaction, ResponseError)
	AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminGetUtxos(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Utxo, ResponseError)
	AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminNewXpub(ctx context.Context, rawXPub string, metadata *models.Metadata) ResponseError
	AdminGetXPubs(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *QueryParams) ([]*models.Xpub, ResponseError)
	AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata) (int64, ResponseError)
	AdminRecordTransaction(ctx context.Context, hex string) (*models.Transaction, ResponseError)
	AdminGetSharedConfig(ctx context.Context) (*models.SharedConfig, ResponseError)
}

// TransportService the transport service interface
type TransportService interface {
	AccessKeyService
	AdminService
	ContactService
	DestinationService
	TransactionService
	XpubService
	Init() error
	IsSignRequest() bool
	SetAdminKey(adminKey *bip32.ExtendedKey)
	SetSignRequest(signRequest bool)
}
