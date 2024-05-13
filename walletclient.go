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
	signRequest     *bool
	server          *string
	accessKeyString *string
	xPrivString     *string
	xPubString      *string
	httpClient      *http.Client
	accessKey       *bec.PrivateKey
	adminXPriv      *bip32.ExtendedKey
	xPriv           *bip32.ExtendedKey
	xPub            *bip32.ExtendedKey
}

// NewWalletClientWithXPrivate creates a new WalletClient instance using a private key (xPriv).
// It configures the client with a specific server URL and a flag indicating whether requests should be signed.
// - `xPriv`: The extended private key used for cryptographic operations.
// - `serverURL`: The URL of the server the client will interact with.
// - `sign`: A boolean flag to determine if the outgoing requests should be signed.
func NewWalletClientWithXPrivate(xPriv, serverURL string, sign bool) (*WalletClient, error) {
	return newWalletClient(
		&WithXPriv{XPrivString: &xPriv},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: &sign},
	)
}

// NewWalletClientWithXPublic creates a new WalletClient instance using a public key (xPub).
// This client is configured for operations that require a public key, such as verifying signatures or receiving transactions.
// - `xPub`: The extended public key used for cryptographic verification and other public operations.
// - `serverURL`: The URL of the server the client will interact with.
// - `sign`: A boolean flag to determine if the outgoing requests should be signed.
func NewWalletClientWithXPublic(xPub, serverURL string, sign bool) (*WalletClient, error) {
	return newWalletClient(
		&WithXPub{XPubString: &xPub},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: &sign},
	)
}

// NewWalletClientWithAdminKey creates a new WalletClient using an administrative key for advanced operations.
// This configuration is typically used for administrative tasks such as managing sub-wallets or configuring system-wide settings.
// - `adminKey`: The extended private key used for administrative operations.
// - `serverURL`: The URL of the server the client will interact with.
// - `sign`: A boolean flag to determine if the outgoing requests should be signed.
func NewWalletClientWithAdminKey(adminKey, serverURL string, sign bool) (*WalletClient, error) {
	return newWalletClient(
		&WithXPriv{XPrivString: &adminKey},
		&WithAdminKey{AdminKeyString: &adminKey},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: &sign},
	)
}

// NewWalletClientWithAccessKey creates a new WalletClient configured with an access key for API authentication.
// This method is useful for scenarios where the client needs to authenticate using a less sensitive key than an xPriv.
// - `accessKey`: The access key used for API authentication.
// - `serverURL`: The URL of the server the client will interact with.
// - `sign`: A boolean flag to determine if the outgoing requests should be signed.
func NewWalletClientWithAccessKey(accessKey, serverURL string, sign bool) (*WalletClient, error) {
	return newWalletClient(
		&WithAccessKey{AccessKeyString: &accessKey},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: &sign},
	)
}

// newWalletClient creates a new WalletClient using the provided configuration options.
func newWalletClient(configurators ...WalletClientConfigurator) (*WalletClient, error) {
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
	case c.xPrivString != nil:
		if c.xPriv, err = bitcoin.GenerateHDKeyFromString(*c.xPrivString); err != nil {
			return err
		}
		if c.xPub, err = c.xPriv.Neuter(); err != nil {
			return err
		}
	case c.xPubString != nil:
		if c.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(*c.xPubString); err != nil {
			return err
		}
	case c.accessKeyString != nil:
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

	if decodedWIF, err = wif.DecodeWIF(*c.accessKeyString); err != nil {
		if privateKey, err = bitcoin.PrivateKeyFromString(*c.accessKeyString); err != nil {
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
