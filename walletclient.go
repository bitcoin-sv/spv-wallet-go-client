package walletclient

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
	"github.com/pkg/errors"
)

// WalletClient is the spv wallet Go client representation.
type WalletClient struct {
	signRequest *bool
	server      *string
	httpClient  *http.Client
	accessKey   *bec.PrivateKey
	adminXPriv  *bip32.ExtendedKey
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// NewWalletClientWithXPrivate creates a new WalletClient instance using a private key (xPriv).
// It configures the client with a specific server URL and a flag indicating whether requests should be signed.
// - `xPriv`: The extended private key used for cryptographic operations.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithXPriv(xPriv, serverURL string) (*WalletClient, error) {
	return newWalletClient(
		&WithXPriv{XPrivString: &xPriv},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: Ptr(true)},
	)
}

// NewWalletClientWithXPublic creates a new WalletClient instance using a public key (xPub).
// This client is configured for operations that require a public key, such as verifying signatures or receiving transactions.
// - `xPub`: The extended public key used for cryptographic verification and other public operations.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithXPub(xPub, serverURL string) (*WalletClient, error) {
	return newWalletClient(
		&WithXPub{XPubString: &xPub},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: Ptr(false)},
	)
}

// NewWalletClientWithAdminKey creates a new WalletClient using an administrative key for advanced operations.
// This configuration is typically used for administrative tasks such as managing sub-wallets or configuring system-wide settings.
// - `adminKey`: The extended private key used for administrative operations.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithAdminKey(adminKey, serverURL string) (*WalletClient, error) {
	return newWalletClient(
		&WithAdminKey{AdminKeyString: &adminKey},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: Ptr(true)},
	)
}

// NewWalletClientWithAccessKey creates a new WalletClient configured with an access key for API authentication.
// This method is useful for scenarios where the client needs to authenticate using a less sensitive key than an xPriv.
// - `accessKey`: The access key used for API authentication.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithAccessKey(accessKey, serverURL string) (*WalletClient, error) {
	return newWalletClient(
		&WithAccessKey{AccessKeyString: &accessKey},
		&WithHTTP{ServerURL: &serverURL},
		&WithSignRequest{Sign: Ptr(true)},
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
	switch {
	case c.xPriv != nil:
		return nil
	case c.xPub != nil:
		return nil
	case c.accessKey != nil:
		return nil
	case c.adminXPriv != nil:
		return nil
	default:
		return errors.New("no keys provided for initialization")
	}
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

func Ptr[T any](obj T) *T {
	return &obj
}

// getPrivKey retrieves the client's private key. If the primary key is not set,
func (wc *WalletClient) getPrivKey() *bip32.ExtendedKey {
	if wc.adminXPriv != nil {
		return wc.adminXPriv
	}
	return wc.xPriv
}
