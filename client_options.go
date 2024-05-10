package walletclient

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// // WithXPriv will set xPrivString on the client
// func WithXPriv(xPrivString string) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.xPrivString = xPrivString
// 		}
// 	}
// }

// // WithXPub will set xPubString on the client
// func WithXPub(xPubString string) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.xPubString = xPubString
// 		}
// 	}
// }

// WithAccessKey will set the access key on the client
// func WithAccessKey(accessKeyString string) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.accessKeyString = accessKeyString
// 		}
// 	}
// }

// // WithHTTP will overwrite the default client with a custom client
// func WithHTTP(serverURL string) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.transportOptions = append(c.transportOptions, transports.WithHTTP(serverURL))
// 		}
// 	}
// }

// // WithHTTPClient will overwrite the default client with a custom client
// func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(serverURL, httpClient))
// 		}
// 	}
// }

// // WithAdminKey will set the admin key for admin requests
// func WithAdminKey(adminKey string) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.transportOptions = append(c.transportOptions, transports.WithAdminKey(adminKey))
// 		}
// 	}
// }

// // WithSignRequest will set whether to sign all requests
// func WithSignRequest(signRequest bool) ClientOps {
// 	return func(c *WalletClient) {
// 		if c != nil {
// 			c.transportOptions = append(c.transportOptions, transports.WithSignRequest(signRequest))
// 		}
// 	}
// }

// WalletClientConfigurator is the interface for configuring WalletClient
type WalletClientConfigurator interface {
	Configure(c *WalletClient)
}

// WithXPriv sets the xPrivString field of a WalletClient
type WithXPriv struct {
	XPrivString string
}

// NewConfigWithXPriv creates new configuration for configurator with xpriv key
func NewConfigWithXPriv(xPrivString string) *WithXPriv {
	return &WithXPriv{xPrivString}
}

// WithXPriv will set xPrivString on the client
func (w *WithXPriv) Configure(c *WalletClient) {
	c.xPrivString = w.XPrivString
}

// WithHTTP sets the URL for the HTTP transport of a WalletClient
type WithHTTP struct {
	ServerURL string
}

func (w *WithHTTP) Configure(c *WalletClient) {
	if c.transportOptions == nil {
		c.transportOptions = []transports.ClientOps{}
	}
	c.transportOptions = append(c.transportOptions, transports.WithHTTP(w.ServerURL))
}

// WithAdminKey sets the admin key for creating new xpubs
type WithAdminKey struct {
	AdminKeyString string
}

// WithAdminKey will set the admin key for admin requests
func (w *WithAdminKey) Configure(c *WalletClient) {
	if c.transportOptions == nil {
		c.transportOptions = []transports.ClientOps{}
	}
	c.transportOptions = append(c.transportOptions, transports.WithAdminKey(w.AdminKeyString))
}

// WithSignRequest configures whether to sign HTTP requests
type WithSignRequest struct {
	Sign bool
}

// WithSignRequest will set whether to sign all requests
func (w *WithSignRequest) Configure(c *WalletClient) {
	if c.transportOptions == nil {
		c.transportOptions = []transports.ClientOps{}
	}
	c.transportOptions = append(c.transportOptions, transports.WithSignRequest(w.Sign))
}

// WithXPub sets the xPubString on the client
type WithXPub struct {
	XPubString string
}

// WithXPub will set xPubString on the client
func (w *WithXPub) Configure(c *WalletClient) {
	c.xPubString = w.XPubString
}

// WithAccessKey sets the accessKeyString on the client.
type WithAccessKey struct {
	AccessKeyString string
}

// WithAccessKey will set the access key on the client
func (w *WithAccessKey) Configure(c *WalletClient) {
	c.accessKeyString = w.AccessKeyString
}

// WithHTTPClient sets a custom HTTP client and server URL for the transport of a WalletClient.
type WithHTTPClient struct {
	ServerURL  string
	HTTPClient *http.Client
}

// WithHTTPClient will overwrite the default client with a custom client
func (w *WithHTTPClient) Configure(c *WalletClient) {
	if c != nil {
		if c.transportOptions == nil {
			c.transportOptions = []transports.ClientOps{}
		}
		// Append the custom HTTP client configuration to the transport options
		c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(w.ServerURL, w.HTTPClient))
	}
}
