package walletclient

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/wif"
	"github.com/pkg/errors"
)

// configurator is the interface for configuring WalletClient
type configurator interface {
	Configure(c *WalletClient)
}

// xPrivConf sets the xPrivString field of a WalletClient
type xPrivConf struct {
	XPrivString string
}

func (w *xPrivConf) Configure(c *WalletClient) {
	var err error
	if c.xPriv, err = bitcoin.GenerateHDKeyFromString(w.XPrivString); err != nil {
		c.xPriv = nil
	}
}

// xPubConf sets the xPubString on the client
type xPubConf struct {
	XPubString string
}

func (w *xPubConf) Configure(c *WalletClient) {
	var err error
	if c.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(w.XPubString); err != nil {
		c.xPub = nil
	}

}

// accessKeyConf sets the accessKeyString on the client
type accessKeyConf struct {
	AccessKeyString string
}

func (w *accessKeyConf) Configure(c *WalletClient) {
	var err error
	if c.accessKey, err = w.initializeAccessKey(); err != nil {
		c.accessKey = nil
	}
}

// adminKeyConf sets the admin key for creating new xpubs
type adminKeyConf struct {
	AdminKeyString string
}

func (w *adminKeyConf) Configure(c *WalletClient) {
	var err error
	c.adminXPriv, err = bitcoin.GenerateHDKeyFromString(w.AdminKeyString)
	if err != nil {
		c.adminXPriv = nil
	}
}

// httpConf sets the URL and httpConf client of a WalletClient
type httpConf struct {
	ServerURL  string
	HTTPClient *http.Client
}

func (w *httpConf) Configure(c *WalletClient) {
	// Ensure the ServerURL ends with a clean base URL
	baseURL, err := validateAndCleanURL(w.ServerURL)
	if err != nil {
		// Handle the error appropriately
		fmt.Println("Invalid URL provided:", err)
		return
	}

	const basePath = "/v1"
	c.server = fmt.Sprintf("%s%s", baseURL, basePath)

	c.httpClient = w.HTTPClient
	if w.HTTPClient != nil {
		c.httpClient = w.HTTPClient
	} else {
		c.httpClient = http.DefaultClient
	}
}

// signRequest configures whether to sign HTTP requests
type signRequest struct {
	Sign bool
}

func (w *signRequest) Configure(c *WalletClient) {
	c.signRequest = w.Sign
}

// initializeAccessKey handles the specific initialization of the access key.
func (w *accessKeyConf) initializeAccessKey() (*bec.PrivateKey, error) {
	var err error
	var privateKey *bec.PrivateKey
	var decodedWIF *wif.WIF

	if decodedWIF, err = wif.DecodeWIF(w.AccessKeyString); err != nil {
		if privateKey, err = bitcoin.PrivateKeyFromString(w.AccessKeyString); err != nil {
			return nil, errors.Wrap(err, "failed to decode access key")
		}
	} else {
		privateKey = decodedWIF.PrivKey
	}

	return privateKey, nil
}

// validateAndCleanURL ensures that the provided URL is valid, and strips it down to just the base URL.
func validateAndCleanURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	// Parse the URL to validate it
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL failed: %w", err)
	}

	// Rebuild the URL with only the scheme and host (and port if included)
	cleanedURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return cleanedURL, nil
	}

	return cleanedURL, nil
}
