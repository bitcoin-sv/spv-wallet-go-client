package walletclient

import (
	"net/http"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/wif"
	"github.com/pkg/errors"
)

// WalletClientConfigurator is the interface for configuring WalletClient
type WalletClientConfigurator interface {
	Configure(c *WalletClient)
}

// XPriv sets the xPrivString field of a WalletClient
type XPriv struct {
	XPrivString *string
}

func (w *XPriv) Configure(c *WalletClient) {
	var err error
	if c.xPriv, err = bitcoin.GenerateHDKeyFromString(*w.XPrivString); err != nil {
		c.xPriv = nil
	}
}

// XPub sets the xPubString on the client
type XPub struct {
	XPubString *string
}

func (w *XPub) Configure(c *WalletClient) {
	var err error
	if c.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(*w.XPubString); err != nil {
		w.XPubString = nil
	}

}

// AccessKey sets the accessKeyString on the client
type AccessKey struct {
	AccessKeyString *string
}

func (w *AccessKey) Configure(c *WalletClient) {
	var err error
	if c.accessKey, err = w.initializeAccessKey(); err != nil {
		c.accessKey = nil
	}
}

// AdminKey sets the admin key for creating new xpubs
type AdminKey struct {
	AdminKeyString *string
}

func (w *AdminKey) Configure(c *WalletClient) {
	var err error
	c.adminXPriv, err = bitcoin.GenerateHDKeyFromString(*w.AdminKeyString)
	if err != nil {
		c.adminXPriv = nil
	}
}

// HTTP sets the URL and HTTP client of a WalletClient
type HTTP struct {
	ServerURL  *string
	HTTPClient *http.Client
}

func (w *HTTP) Configure(c *WalletClient) {
	c.server = *w.ServerURL
	c.httpClient = w.HTTPClient
	if w.HTTPClient != nil {
		c.httpClient = w.HTTPClient
	} else {
		c.httpClient = http.DefaultClient
	}
}

// SignRequest configures whether to sign HTTP requests
type SignRequest struct {
	Sign *bool
}

func (w *SignRequest) Configure(c *WalletClient) {
	c.signRequest = *w.Sign
}

// initializeAccessKey handles the specific initialization of the access key.
func (c *AccessKey) initializeAccessKey() (*bec.PrivateKey, error) {
	var err error
	var privateKey *bec.PrivateKey
	var decodedWIF *wif.WIF

	if decodedWIF, err = wif.DecodeWIF(*c.AccessKeyString); err != nil {
		if privateKey, err = bitcoin.PrivateKeyFromString(*c.AccessKeyString); err != nil {
			return nil, errors.Wrap(err, "failed to decode access key")
		}
	} else {
		privateKey = decodedWIF.PrivKey
	}

	return privateKey, nil
}
