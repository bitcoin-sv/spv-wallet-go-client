// Package buxclient is a Go client for interacting with Bux Servers
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By MrZ (https://github.com/mrz1836)
// By icellan (https://github.com/icellan)
package buxclient

import (
	"context"
	"errors"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/BuxOrg/go-buxclient/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
)

// ClientOps ...
type ClientOps func(c *BuxClient)

// BuxClient the bux client
type BuxClient struct {
	// private fields
	debug            bool
	transport        transports.TransportService
	xPrivString      string
	xPubString       string
	accessKeyString  string
	xPriv            *bip32.ExtendedKey
	xPub             *bip32.ExtendedKey
	accessKey        *bec.PrivateKey
	transportOptions []transports.ClientOps
}

// New create a new bux client
func New(opts ...ClientOps) (*BuxClient, error) {
	client := &BuxClient{}

	for _, opt := range opts {
		opt(client)
	}

	var err error
	if client.xPrivString != "" {
		client.xPriv, err = bip32.NewKeyFromString(client.xPrivString)
		if err != nil {
			return nil, err
		}
		client.xPub, err = client.xPriv.Neuter()
		if err != nil {
			return nil, err
		}
	} else if client.xPubString != "" {
		client.xPriv = nil
		client.xPub, err = bip32.NewKeyFromString(client.xPubString)
		if err != nil {
			return nil, err
		}
	} else if client.accessKeyString != "" {
		client.xPriv = nil
		client.xPub = nil
		var decodedWIF *wif.WIF
		decodedWIF, err = wif.DecodeWIF(client.accessKeyString)
		if err != nil {
			return nil, err
		}
		client.accessKey = decodedWIF.PrivKey
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

	client.transport, err = transports.NewTransport(transportOptions...)
	if err != nil {
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

// GetTransport returns the current transport service
func (b *BuxClient) GetTransport() *transports.TransportService {
	return &b.transport
}

// RegisterXpub registers a new xpub - admin key needed
func (b *BuxClient) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {
	return b.transport.RegisterXpub(ctx, rawXPub, metadata)
}

// DraftTransaction initialize a new draft transaction
func (b *BuxClient) DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig, metadata *bux.Metadata) (*bux.DraftTransaction, error) {
	return b.transport.DraftTransaction(ctx, transactionConfig, metadata)
}

// DraftToRecipients initialize a new P2PKH draft transaction to a list of recipients
func (b *BuxClient) DraftToRecipients(ctx context.Context, recipients []*transports.Recipients, metadata *bux.Metadata) (*bux.DraftTransaction, error) {
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

// RecordTransaction record a new transaction
func (b *BuxClient) RecordTransaction(ctx context.Context, hex, draftID string, metadata *bux.Metadata) (string, error) {
	return b.transport.RecordTransaction(ctx, hex, draftID, metadata)
}

// SendToRecipients send to recipients
func (b *BuxClient) SendToRecipients(ctx context.Context, recipients []*transports.Recipients, metadata *bux.Metadata) (string, error) {
	draft, err := b.DraftToRecipients(ctx, recipients, metadata)
	if err != nil {
		return "", err
	}

	var hex string
	if hex, err = b.FinalizeTransaction(draft); err != nil {
		return "", err
	}

	return b.RecordTransaction(ctx, hex, draft.ID, metadata)
}

// WithXPriv will set xPrivString on the client
func WithXPriv(xPrivString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.xPrivString = xPrivString
		}
	}
}

// WithXPub will set xPubString on the client
func WithXPub(xPubString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.xPubString = xPubString
		}
	}
}

// WithAccessKey will set accessKey on the client
func WithAccessKey(accessKeyString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.accessKeyString = accessKeyString
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(serverURL))
		}
	}
}

// WithGraphQLClient will overwrite the default client with a custom client
func WithGraphQLClient(serverURL string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithGraphQLClient(serverURL))
		}
	}
}

// WithClient will overwrite the default client with a custom client
func WithClient(transportClient transports.TransportService) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithClient(transportClient))
		}
	}
}

// WithAdminKey will set the admin key for admin requests
func WithAdminKey(adminKey string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithAdminKey(adminKey))
		}
	}
}

// WithSignRequest will set whether to sign all requests
func WithSignRequest(signRequest bool) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithSignRequest(signRequest))
		}
	}
}

// WithDebugging will set whether to turn debugging on
func WithDebugging(debug bool) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithDebugging(debug))
		}
	}
}
