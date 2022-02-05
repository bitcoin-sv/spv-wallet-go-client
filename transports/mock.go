package transports

import (
	"context"

	"github.com/BitcoinSchema/xapi/bux"
	"github.com/libsv/go-bk/bip32"
)

// TransportMock is the struct for Mock transport that can be used for testing
type TransportMock struct {
	adminXPriv  *bip32.ExtendedKey
	debug       bool
	server      string
	signRequest bool
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
	callback    func(interface{}, interface{}, interface{}) (interface{}, error)
}

// Init will initialize
func (m *TransportMock) Init() error {
	return nil
}

// SetCallback set the callback function when a function is called
func (m *TransportMock) SetCallback(callback func(interface{}, interface{}, interface{}) (interface{}, error)) {
	m.callback = callback
}

// SetDebug turn the debugging on or off
func (m *TransportMock) SetDebug(debug bool) {
	m.debug = debug
}

// SetAdminKey set the admin key
func (m *TransportMock) SetAdminKey(adminKey *bip32.ExtendedKey) {
	m.adminXPriv = adminKey
}

// RegisterXpub will register an xPub
func (m *TransportMock) RegisterXpub(ctx context.Context, rawXPub string) error {
	_, err := m.callback(ctx, rawXPub, nil)
	return err
}

// GetDestination will get a destination
func (m *TransportMock) GetDestination(ctx context.Context) (*bux.Destination, error) {
	destination, err := m.callback(ctx, nil, nil)
	return destination.(*bux.Destination), err
}

// DraftTransaction is a draft transaction
func (m *TransportMock) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig) (*bux.DraftTransaction, error) {
	draftTransaction, err := m.callback(ctx, transactionConfig, nil)
	return draftTransaction.(*bux.DraftTransaction), err
}

// DraftToRecipients is a draft transaction to a slice of recipients
func (m *TransportMock) DraftToRecipients(ctx context.Context, recipients []*Recipients) (*bux.DraftTransaction, error) {
	draftTransaction, err := m.callback(ctx, recipients, nil)
	return draftTransaction.(*bux.DraftTransaction), err
}

// RecordTransaction will record a transaction
func (m *TransportMock) RecordTransaction(ctx context.Context, hex, referenceID string) (string, error) {
	txID, err := m.callback(ctx, hex, referenceID)
	return txID.(string), err
}
