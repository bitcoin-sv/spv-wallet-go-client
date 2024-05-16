package walletclient

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

// WalletClient is the spv wallet Go client representation.
type WalletClient struct {
	signRequest bool
	server      string
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
func NewWithXPriv(serverURL, xPriv string) *WalletClient {
	return makeClient(
		&xPrivConf{XPrivString: xPriv},
		&httpConf{ServerURL: serverURL},
		&signRequest{Sign: true},
	)
}

// NewWalletClientWithXPublic creates a new WalletClient instance using a public key (xPub).
// This client is configured for operations that require a public key, such as verifying signatures or receiving transactions.
// - `xPub`: The extended public key used for cryptographic verification and other public operations.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithXPub(serverURL, xPub string) *WalletClient {
	return makeClient(
		&xPubConf{XPubString: xPub},
		&httpConf{ServerURL: serverURL},
		&signRequest{Sign: false},
	)
}

// NewWalletClientWithAdminKey creates a new WalletClient using an administrative key for advanced operations.
// This configuration is typically used for administrative tasks such as managing sub-wallets or configuring system-wide settings.
// - `adminKey`: The extended private key used for administrative operations.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithAdminKey(serverURL, adminKey string) *WalletClient {
	return makeClient(
		&adminKeyConf{AdminKeyString: adminKey},
		&httpConf{ServerURL: serverURL},
		&signRequest{Sign: true},
	)
}

// NewWalletClientWithAccessKey creates a new WalletClient configured with an access key for API authentication.
// This method is useful for scenarios where the client needs to authenticate using a less sensitive key than an xPriv.
// - `accessKey`: The access key used for API authentication.
// - `serverURL`: The URL of the server the client will interact with.
func NewWithAccessKey(serverURL, accessKey string) *WalletClient {
	return makeClient(
		&accessKeyConf{AccessKeyString: accessKey},
		&httpConf{ServerURL: serverURL},
		&signRequest{Sign: true},
	)
}

// makeClient creates a new WalletClient using the provided configuration options.
func makeClient(configurators ...Configurator) *WalletClient {
	client := &WalletClient{}

	for _, configurator := range configurators {
		configurator.Configure(client)
	}

	return client
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
