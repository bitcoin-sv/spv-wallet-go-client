package walletclient

import (
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/wif"
	"github.com/pkg/errors"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// WalletClient is the spv wallet Go client representation.
type WalletClient struct {
	transports.TransportService
	accessKey        *bec.PrivateKey
	accessKeyString  string
	transport        transports.TransportService
	transportOptions []transports.ClientOps
	xPriv            *bip32.ExtendedKey
	xPrivString      string
	xPub             *bip32.ExtendedKey
	xPubString       string
}

// New creates a new WalletClient using the provided configuration options.
func New(configurators ...WalletClientConfigurator) (*WalletClient, error) {
	client := &WalletClient{}

	for _, configurator := range configurators {
		configurator.Configure(client)
	}

	// Initialize keys based on provided strings
	if err := client.initializeKeys(); err != nil {
		return nil, err
	}

	// Setup transport based on initialized keys
	if err := client.setupTransport(); err != nil {
		return nil, err
	}

	client.TransportService = client.transport

	return client, nil
}

// initializeKeys handles the initialization of keys based on the existing fields.
func (c *WalletClient) initializeKeys() error {
	var err error
	switch {
	case c.xPrivString != "":
		if c.xPriv, err = bitcoin.GenerateHDKeyFromString(c.xPrivString); err != nil {
			return err
		}
		if c.xPub, err = c.xPriv.Neuter(); err != nil {
			return err
		}
	case c.xPubString != "":
		if c.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(c.xPubString); err != nil {
			return err
		}
	case c.accessKeyString != "":
		return c.initializeAccessKey()
	default:
		return errors.New("no keys provided for initialization")
	}
	return nil
}

// initializeAccessKey handles the specific initialization of the access key.
func (c *WalletClient) initializeAccessKey() error {
	var err error
	var privateKey *bec.PrivateKey
	var decodedWIF *wif.WIF

	if decodedWIF, err = wif.DecodeWIF(c.accessKeyString); err != nil {
		if privateKey, err = bitcoin.PrivateKeyFromString(c.accessKeyString); err != nil {
			return errors.Wrap(err, "failed to decode access key")
		}
	} else {
		privateKey = decodedWIF.PrivKey
	}

	c.accessKey = privateKey
	return nil
}

// setupTransport configures the transport service based on the available keys.
func (c *WalletClient) setupTransport() error {
	var err error
	transportOptions := make([]transports.ClientOps, 0)

	if c.xPriv != nil {
		transportOptions = append(transportOptions, transports.WithXPriv(c.xPriv))
		transportOptions = append(transportOptions, transports.WithXPub(c.xPub))
	} else if c.xPub != nil {
		transportOptions = append(transportOptions, transports.WithXPub(c.xPub))
	} else if c.accessKey != nil {
		transportOptions = append(transportOptions, transports.WithAccessKey(c.accessKey))
	}

	if len(c.transportOptions) > 0 {
		transportOptions = append(transportOptions, c.transportOptions...)
	}

	if c.transport, err = transports.NewTransport(transportOptions...); err != nil {
		return err
	}

	return nil
}
