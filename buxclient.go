// Package buxclient is a Go client for interacting with Bux Servers
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BuxOrg (https://github.com/BuxOrg)
package buxclient

import (
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/wif"
	"github.com/pkg/errors"
)

// ClientOps are used for client options
type ClientOps func(c *BuxClient)

// BuxClient is the go-buxclient
type BuxClient struct {
	transports.TransportService
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
		if client.xPriv, err = bitcoin.GenerateHDKeyFromString(client.xPrivString); err != nil {
			return nil, err
		}
		if client.xPub, err = client.xPriv.Neuter(); err != nil {
			return nil, err
		}
	} else if client.xPubString != "" {
		client.xPriv = nil
		if client.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(client.xPubString); err != nil {
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

	client.TransportService = client.transport

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
