package transports

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/libsv/go-bk/bip32"
)

// XpubService is the xPub related requests
type XpubService interface {
	GetXPub(ctx context.Context) (*bux.Xpub, error)
	NewXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error
	UpdateXPubMetadata(ctx context.Context, metadata *bux.Metadata) (*bux.Xpub, error)
}

// AccessKeyService is the access key related requests
type AccessKeyService interface {
	CreateAccessKey(ctx context.Context, metadata *bux.Metadata) (*bux.AccessKey, error)
	GetAccessKey(ctx context.Context, id string) (*bux.AccessKey, error)
	GetAccessKeys(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.AccessKey, error)
	RevokeAccessKey(ctx context.Context, id string) (*bux.AccessKey, error)
}

// DestinationService is the destination related requests
type DestinationService interface {
	GetDestinationByAddress(ctx context.Context, address string) (*bux.Destination, error)
	GetDestinationByID(ctx context.Context, id string) (*bux.Destination, error)
	GetDestinationByLockingScript(ctx context.Context, lockingScript string) (*bux.Destination, error)
	GetDestinations(ctx context.Context, metadataConditions *bux.Metadata) ([]*bux.Destination, error)
	NewDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error)
	UpdateDestinationMetadataByAddress(ctx context.Context, lockingScript string, metadata *bux.Metadata) (*bux.Destination, error)
	UpdateDestinationMetadataByID(ctx context.Context, id string, metadata *bux.Metadata) (*bux.Destination, error)
	UpdateDestinationMetadataByLockingScript(ctx context.Context, address string, metadata *bux.Metadata) (*bux.Destination, error)
}

// TransactionService is the transaction related requests
type TransactionService interface {
	DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata *bux.Metadata) (*bux.DraftTransaction, error)
	DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig, metadata *bux.Metadata) (*bux.DraftTransaction, error)
	GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error)
	GetTransactions(ctx context.Context, conditions map[string]interface{}, metadataConditions *bux.Metadata) ([]*bux.Transaction, error)
	RecordTransaction(ctx context.Context, hex, referenceID string, metadata *bux.Metadata) (*bux.Transaction, error)
	UpdateTransactionMetadata(ctx context.Context, txID string, metadata *bux.Metadata) (*bux.Transaction, error)
}

// PaymailService is the paymail related requests
type PaymailService interface {
	NewPaymail(ctx context.Context, rawXpub, paymailAddress string, metadata *bux.Metadata) error
}

// TransportService the transport service interface
type TransportService interface {
	AccessKeyService
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
