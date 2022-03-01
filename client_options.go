package buxclient

import (
	"net/http"

	"github.com/BuxOrg/go-buxclient/transports"
)

// WithXPriv will set xPrivString on the client
func WithXPriv(xPrivString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.xPrivString = xPrivString
		}
	}
}

// WithXPub will set xPubString on the client
func WithXPub(xPubString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.xPubString = xPubString
		}
	}
}

// WithAccessKey will set accessKey on the client
func WithAccessKey(accessKeyString string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.accessKeyString = accessKeyString
		}
	}
}

// WithHTTP will overwrite the default client with a custom client
func WithHTTP(serverURL string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTP(serverURL))
		}
	}
}

// WithGraphQL will overwrite the default client with a custom client
func WithGraphQL(serverURL string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithGraphQL(serverURL))
		}
	}
}

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithHTTPClient(serverURL, httpClient))
		}
	}
}

// WithGraphQLClient will overwrite the default client with a custom client
func WithGraphQLClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithGraphQLClient(serverURL, httpClient))
		}
	}
}

// WithAdminKey will set the admin key for admin requests
func WithAdminKey(adminKey string) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithAdminKey(adminKey))
		}
	}
}

// WithSignRequest will set whether to sign all requests
func WithSignRequest(signRequest bool) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithSignRequest(signRequest))
		}
	}
}

// WithDebugging will set whether to turn debugging on
func WithDebugging(debug bool) ClientOps {
	return func(c *BuxClient) {
		if c != nil {
			c.transportOptions = append(c.transportOptions, transports.WithDebugging(debug))
		}
	}
}
