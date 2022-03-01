// Package buxclient is a Go client for interacting with Bux Servers
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BuxOrg (https://github.com/BuxOrg)
package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/BuxOrg/go-buxclient/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/pkg/errors"
)

// ClientOps are used for client options
type ClientOps func(c *BuxClient)

// BuxClient is the bux client
type BuxClient struct {
	accessKey        *bec.PrivateKey
	accessKeyString  string
	debug            bool
	transport        transports.TransportService
	transportOptions []transports.ClientOps
	xPriv            *bip32.ExtendedKey
	xPrivString      string
	xPub             *bip32.ExtendedKey
	xPubString       string
}

// New create a new bux client
func New(opts ...ClientOps) (*BuxClient, error) {
	client := &BuxClient{}

	for _, opt := range opts {
		opt(client)
	}

	var err error
	if client.xPrivString != "" {
		if client.xPriv, err = bip32.NewKeyFromString(client.xPrivString); err != nil {
			return nil, err
		}
		if client.xPub, err = client.xPriv.Neuter(); err != nil {
			return nil, err
		}
	} else if client.xPubString != "" {
		client.xPriv = nil
		if client.xPub, err = bip32.NewKeyFromString(client.xPubString); err != nil {
			return nil, err
		}
	} else if client.accessKeyString != "" {
		client.xPriv = nil
		client.xPub = nil

		var privateKey *bec.PrivateKey
		var decodedWIF *wif.WIF
		if decodedWIF, err = wif.DecodeWIF(client.accessKeyString); err != nil {
			// try as a hex string
			var errHex error
			if privateKey, errHex = bitcoin.PrivateKeyFromString(client.accessKeyString); errHex != nil {
				return nil, errors.Wrap(err, errHex.Error())
			}
		} else {
			privateKey = decodedWIF.PrivKey
		}
		client.accessKey = privateKey
	} else {
		return nil, errors.New("no keys available")
	}

	transportOptions := make([]transports.ClientOps, 0)
	if client.xPriv != nil {
		transportOptions = append(transportOptions, transports.WithXPriv(client.xPriv))
		transportOptions = append(transportOptions, transports.WithXPub(client.xPub))
	} else if client.xPub != nil {
		transportOptions = append(transportOptions, transports.WithXPub(client.xPub))
	} else if client.accessKey != nil {
		transportOptions = append(transportOptions, transports.WithAccessKey(client.accessKey))
	}
	if len(client.transportOptions) > 0 {
		transportOptions = append(transportOptions, client.transportOptions...)
	}

	if client.transport, err = transports.NewTransport(transportOptions...); err != nil {
		return nil, err
	}

	return client, nil
}

// SetAdminKey set the admin key to use to create new xpubs
func (b *BuxClient) SetAdminKey(adminKeyString string) error {
	adminKey, err := bip32.NewKeyFromString(adminKeyString)
	if err != nil {
		return err
	}

	b.transport.SetAdminKey(adminKey)

	return nil
}

// SetDebug turn the debugging on or off
func (b *BuxClient) SetDebug(debug bool) {
	b.debug = debug
	b.transport.SetDebug(debug)
}

// IsDebug return the debugging status
func (b *BuxClient) IsDebug() bool {
	return b.transport.IsDebug()
}

// SetSignRequest turn the signing of the http request on or off
func (b *BuxClient) SetSignRequest(signRequest bool) {
	b.transport.SetSignRequest(signRequest)
}

// IsSignRequest return whether to sign all requests
func (b *BuxClient) IsSignRequest() bool {
	return b.transport.IsSignRequest()
}

// GetTransport returns the current transport service
func (b *BuxClient) GetTransport() *transports.TransportService {
	return &b.transport
}

// RegisterXpub registers a new xpub - admin key needed
func (b *BuxClient) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {
	return b.transport.RegisterXpub(ctx, rawXPub, metadata)
}

// DraftTransaction initialize a new draft transaction
func (b *BuxClient) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	return b.transport.DraftTransaction(ctx, transactionConfig, metadata)
}

// DraftToRecipients initialize a new P2PKH draft transaction to a list of recipients
func (b *BuxClient) DraftToRecipients(ctx context.Context, recipients []*transports.Recipients,
	metadata *bux.Metadata) (*bux.DraftTransaction, error) {

	return b.transport.DraftToRecipients(ctx, recipients, metadata)
}

// GetDestination get new fresh destination
func (b *BuxClient) GetDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error) {
	return b.transport.GetDestination(ctx, metadata)
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

// GetTransaction get a transaction by id
func (b *BuxClient) GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error) {
	return b.transport.GetTransaction(ctx, txID)
}

// GetTransactions get all transactions matching search criteria
func (b *BuxClient) GetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) ([]*bux.Transaction, error) {

	return b.transport.GetTransactions(ctx, conditions, metadata)
}

// RecordTransaction record a new transaction
func (b *BuxClient) RecordTransaction(ctx context.Context, hex, draftID string,
	metadata *bux.Metadata) (*bux.Transaction, error) {

	return b.transport.RecordTransaction(ctx, hex, draftID, metadata)
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
