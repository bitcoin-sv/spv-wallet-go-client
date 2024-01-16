package transports

import (
	"net/http"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

// WithXPriv will set the xPriv
func WithXPriv(xPriv *bip32.ExtendedKey) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.xPriv = xPriv
		}
	}
}

// WithXPub will set the xPub
func WithXPub(xPub *bip32.ExtendedKey) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.xPub = xPub
		}
	}
}

// WithAccessKey will set the access key
func WithAccessKey(accessKey *bec.PrivateKey) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.accessKey = accessKey
		}
	}
}

// WithHTTP will overwrite the default client with a custom client
func WithHTTP(serverURL string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportHTTP{
				server:      serverURL,
				signRequest: c.signRequest,
				adminXPriv:  c.adminXPriv,
				httpClient:  &http.Client{},
				xPriv:       c.xPriv,
				xPub:        c.xPub,
				accessKey:   c.accessKey,
			})
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportHTTP{
				server:      serverURL,
				signRequest: c.signRequest,
				adminXPriv:  c.adminXPriv,
				httpClient:  httpClient,
				xPriv:       c.xPriv,
				xPub:        c.xPub,
				accessKey:   c.accessKey,
			})
		}
	}
}

// WithAdminKey will set the admin key for admin requests
func WithAdminKey(adminKey string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.adminKey = adminKey
		}
	}
}

// WithSignRequest will set whether to sign all requests
func WithSignRequest(signRequest bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.signRequest = signRequest
			if c.transport != nil {
				c.transport.SetSignRequest(signRequest)
			}
		}
	}
}
