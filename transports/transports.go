// Package transports encapsulates the different ways to communicate with SPV Wallet
package transports

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

// Client is the transport client
type Client struct {
	accessKey   *bec.PrivateKey
	adminKey    string
	adminXPriv  *bip32.ExtendedKey
	signRequest bool
	transport   TransportService
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// ClientOps are the client options functions
type ClientOps func(c *Client)

// addSignature will add the signature to the request
func addSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) ResponseError {
	return setSignature(header, xPriv, bodyString)
}

// NewTransport create a new transport service object
func NewTransport(opts ...ClientOps) (TransportService, error) {
	client := Client{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.transport == nil {
		return nil, ErrNoClientSet
	}

	if err := client.transport.Init(); err != nil {
		return nil, err
	}

	if client.adminKey != "" {
		adminXPriv, err := bip32.NewKeyFromString(client.adminKey)
		if err != nil {
			return nil, err
		}
		client.adminXPriv = adminXPriv
		client.transport.SetAdminKey(adminXPriv)
	}

	return client.transport, nil
}

// NewTransportService create a new transport service interface
func NewTransportService(transportService TransportService) TransportService {
	return transportService
}

// processMetadata will process the metadata
func processMetadata(metadata *models.Metadata) *models.Metadata {
	if metadata == nil {
		m := make(models.Metadata)
		metadata = &m
	}

	return metadata
}
