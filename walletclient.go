package walletclient

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/wif"
	"github.com/pkg/errors"
)

// WalletClient is the spv wallet Go client representation.
type WalletClient struct {
	accessKeyString string
	xPrivString     string
	xPubString      string
	accessKey       *bec.PrivateKey
	adminXPriv      *bip32.ExtendedKey
	httpClient      *http.Client
	server          string
	signRequest     bool
	xPriv           *bip32.ExtendedKey
	xPub            *bip32.ExtendedKey
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

// processMetadata will process the metadata
func processMetadata(metadata *models.Metadata) *models.Metadata {
	if metadata == nil {
		m := make(models.Metadata)
		metadata = &m
	}

	return metadata
}

// addSignature will add the signature to the request
func addSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) ResponseError {
	return setSignature(header, xPriv, bodyString)
}
