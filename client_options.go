package walletclient

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// WithXPriv will set xPrivString on the client
func WithXPriv(xPrivString string) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.xPrivString = xPrivString
		}
	}
}

// WithXPub will set xPubString on the client
func WithXPub(xPubString string) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.xPubString = xPubString
		}
	}
}

// WithAccessKey will set the access key on the client
func WithAccessKey(accessKeyString string) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.accessKeyString = accessKeyString
		}
	}
}

// WithHTTP will overwrite the default client with a custom client
func WithHTTP(serverURL string) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTP(serverURL))
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(serverURL, httpClient))
		}
	}
}

// WithAdminKey will set the admin key for admin requests
func WithAdminKey(adminKey string) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithAdminKey(adminKey))
		}
	}
}

// WithSignRequest will set whether to sign all requests
func WithSignRequest(signRequest bool) ClientOps {
	return func(c *WalletClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithSignRequest(signRequest))
		}
	}
}
