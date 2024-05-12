package walletclient

import (
	"fmt"
	"net/http"
)

// WalletClientConfigurator is the interface for configuring WalletClient
type WalletClientConfigurator interface {
	Configure(c *WalletClient)
}

// WithXPriv sets the xPrivString field of a WalletClient
type WithXPriv struct {
	XPrivString string
}

func (w *WithXPriv) Configure(c *WalletClient) {
	fmt.Printf("withXpriv configure: %#v\n", w)
	c.xPrivString = w.XPrivString
}

// WithXPub sets the xPubString on the client
type WithXPub struct {
	XPubString string
}

func (w *WithXPub) Configure(c *WalletClient) {
	c.xPubString = w.XPubString
}

// WithAccessKey sets the accessKeyString on the client
type WithAccessKey struct {
	AccessKeyString string
}

func (w *WithAccessKey) Configure(c *WalletClient) {
	c.accessKeyString = w.AccessKeyString
}

// WithAdminKey sets the admin key for creating new xpubs
type WithAdminKey struct {
	AdminKeyString string
}

func (w *WithAdminKey) Configure(c *WalletClient) {
	fmt.Printf("withAdminKey configure: %#v\n", w)
	fmt.Printf("withAdminKey configure adminxpriv: %v  \n", w.AdminKeyString)
	c.adminXPriv = w.AdminKeyString
}

// WithHTTP sets the URL and HTTP client of a WalletClient
type WithHTTP struct {
	ServerURL  string
	HTTPClient *http.Client
}

func (w *WithHTTP) Configure(c *WalletClient) {
	c.server = w.ServerURL
	c.httpClient = w.HTTPClient
}

// WithSignRequest configures whether to sign HTTP requests
type WithSignRequest struct {
	Sign bool
}

func (w *WithSignRequest) Configure(c *WalletClient) {
	c.signRequest = w.Sign
}
