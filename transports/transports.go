package transports

import (
	"context"
	"errors"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

// TransportType the type of transport being used (http or graphql)
type TransportType string

// BuxUserAgent the bux user agent sent to the bux server
const BuxUserAgent = "BuxClient v1.0.0"

const (
	// BuxTransportHTTP uses the http transport for all bux server actions
	BuxTransportHTTP TransportType = "http"
	// BuxTransportGraphQL uses the graphql transport for all bux server actions
	BuxTransportGraphQL TransportType = "graphql"
	// BuxTransportMock uses the mock transport for all bux server actions
	BuxTransportMock TransportType = "mock"
)

// Client ...
type Client struct {
	accessKey   *bec.PrivateKey
	adminKey    string
	adminXPriv  *bip32.ExtendedKey
	debug       bool
	signRequest bool
	transport   TransportService
	xPriv       *bip32.ExtendedKey
	xPub        *bip32.ExtendedKey
}

// ClientOps ...
type ClientOps func(c *Client)

// addSignature will add the signature to the request
func addSignature(header *http.Header, xPriv *bip32.ExtendedKey, bodyString string) error {
	return bux.SetSignature(header, xPriv, bodyString)
}

// TransportService the transport service interface
type TransportService interface {
	Init() error
	SetAdminKey(adminKey *bip32.ExtendedKey)
	SetDebug(debug bool)
	IsDebug() bool
	SetSignRequest(debug bool)
	IsSignRequest() bool
	RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error
	GetXpub(ctx context.Context, rawXPub string) (*bux.Xpub, error)
	RegisterPaymail(ctx context.Context, rawXpub, paymailAddress string, metadata *bux.Metadata) error
	GetDestination(ctx context.Context, metadata *bux.Metadata) (*bux.Destination, error)
	GetTransaction(ctx context.Context, txID string) (*bux.Transaction, error)
	GetTransactions(ctx context.Context, conditions map[string]interface{}, metadata *bux.Metadata) ([]*bux.Transaction, error)
	DraftToRecipients(ctx context.Context, recipients []*Recipients, metadata *bux.Metadata) (*bux.DraftTransaction, error)
	DraftTransaction(ctx context.Context, transactionConfig *bux.TransactionConfig, metadata *bux.Metadata) (*bux.DraftTransaction, error)
	RecordTransaction(ctx context.Context, hex, referenceID string, metadata *bux.Metadata) (*bux.Transaction, error)
}

// NewTransport create a new transport service object
func NewTransport(opts ...ClientOps) (TransportService, error) {
	client := Client{}

	for _, opt := range opts {
		opt(&client)
	}

	if client.transport == nil {
		return nil, errors.New("no transport client set")
	}

	if err := client.transport.Init(); err != nil {
		return nil, err
	}

	if client.adminKey != "" {
		adminXPriv, err := bip32.NewKeyFromString(client.adminKey)
		if err != nil {
			return nil, err
		}
		client.adminXPriv = adminXPriv
		client.transport.SetAdminKey(adminXPriv)
	}

	return client.transport, nil
}

// NewTransportService create a new transport service interface
func NewTransportService(transportService TransportService) TransportService {
	return transportService
}

func processMetadata(metadata *bux.Metadata) *bux.Metadata {
	if metadata == nil {
		m := make(bux.Metadata)
		metadata = &m
	}

	(*metadata)["user_agent"] = BuxUserAgent

	return metadata
}

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

// WithGraphQL will overwrite the default client with a custom client
func WithGraphQL(serverURL string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportGraphQL{
				debug:       c.debug,
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

// WithHTTP will overwrite the default client with a custom client
func WithHTTP(serverURL string) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportHTTP{
				debug:       c.debug,
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

// WithGraphQLClient will overwrite the default client with a custom client
func WithGraphQLClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportGraphQL{
				debug:       c.debug,
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

// WithHTTPClient will overwrite the default client with a custom client
func WithHTTPClient(serverURL string, httpClient *http.Client) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.transport = NewTransportService(&TransportHTTP{
				debug:       c.debug,
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

// WithDebugging will set whether to turn debugging on
func WithDebugging(debug bool) ClientOps {
	return func(c *Client) {
		if c != nil {
			c.debug = debug
			if c.transport != nil {
				c.transport.SetDebug(debug)
			}
		}
	}
}
