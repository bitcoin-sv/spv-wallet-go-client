package walletclient

import (
	"fmt"
	"net/http"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"

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
	if c.xPriv, err = bip32.GenerateHDKeyFromString(w.XPrivString); err != nil {
		c.xPriv = nil
	}
}

// xPubConf sets the xPubString on the client
type xPubConf struct {
	XPubString string
}

func (w *xPubConf) Configure(c *WalletClient) {
	var err error
	if c.xPub, err = bip32.GetHDKeyFromExtendedPublicKey(w.XPubString); err != nil {
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
	c.adminXPriv, err = bip32.GenerateHDKeyFromString(w.AdminKeyString)
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

	const basePath = "/api/v1"
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

func (w *accessKeyConf) initializeAccessKey() (*ec.PrivateKey, error) {
	privateKey, err := ec.PrivateKeyFromWif(w.AccessKeyString)
	if err != nil {
		privateKey, _ = ec.PrivateKeyFromHex(w.AccessKeyString)
		if privateKey == nil {
			return nil, errors.New("failed to decode access key")
		}
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
