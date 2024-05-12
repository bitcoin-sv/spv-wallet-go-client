package walletclient

// import (
// 	"context"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// 

// 	"github.com/bitcoin-sv/spv-wallet/models"
// 	"github.com/bitcoinschema/go-bitcoin/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/bitcoin-sv/spv-wallet-go-client/fixtures"
// 	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
// )

// // localRoundTripper is an http.RoundTripper that executes HTTP transactions
// // by using handler directly, instead of going over an HTTP connection.
// type localRoundTripper struct {
// 	handler http.Handler
// }

// func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	w := httptest.NewRecorder()
// 	l.handler.ServeHTTP(w, req)
// 	return w.Result(), nil
// }

// func mustWrite(w io.Writer, s string) {
// 	_, err := io.WriteString(w, s)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// type testTransportHandler struct {
// 	ClientURL string
// 	Client    func(serverURL string, httpClient *http.Client) ClientOps
// 	Path      string
// 	Queries   []*testTransportHandlerRequest
// 	Result    string
// 	Type      string
// }

// type testTransportHandlerRequest struct {
// 	Path   string
// 	Result func(w http.ResponseWriter, req *http.Request)
// }

// // TestNewWalletClient will test the TestNewWalletClient method
// func TestNewWalletClient(t *testing.T) {
// 	t.Run("no keys", func(t *testing.T) {
// 		client, err := New()
// 		assert.Error(t, err)
// 		assert.Nil(t, client)
// 	})

// 	t.Run("empty xpriv", func(t *testing.T) {
// 		client, err := New(
// 			WithXPriv(""),
// 		)
// 		assert.Error(t, err)
// 		assert.Nil(t, client)
// 	})

// 	t.Run("invalid xpriv", func(t *testing.T) {
// 		client, err := New(
// 			WithXPriv("invalid-xpriv"),
// 		)
// 		assert.Error(t, err)
// 		assert.Nil(t, client)
// 	})

// 	t.Run("valid client", func(t *testing.T) {
// 		client, err := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		require.NoError(t, err)
// 		assert.IsType(t, WalletClient{}, *client)
// 	})

// 	t.Run("valid xPub client", func(t *testing.T) {
// 		client, err := New(
// 			WithXPub(fixtures.XPubString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		require.NoError(t, err)
// 		assert.IsType(t, WalletClient{}, *client)
// 	})

// 	t.Run("invalid xPub client", func(t *testing.T) {
// 		client, err := New(
// 			WithXPub("invalid-xpub"),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		assert.Error(t, err)
// 		assert.Nil(t, client)
// 	})

// 	t.Run("valid access keys", func(t *testing.T) {
// 		client, err := New(
// 			WithAccessKey(fixtures.AccessKeyString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		require.NoError(t, err)
// 		assert.IsType(t, WalletClient{}, *client)
// 	})

// 	t.Run("invalid access keys", func(t *testing.T) {
// 		client, err := New(
// 			WithAccessKey("invalid-access-key"),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		assert.Error(t, err)
// 		assert.Nil(t, client)
// 	})

// 	t.Run("valid access key WIF", func(t *testing.T) {
// 		wifKey, _ := bitcoin.PrivateKeyToWif(fixtures.AccessKeyString)
// 		client, err := New(
// 			WithAccessKey(wifKey.String()),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		require.NoError(t, err)
// 		assert.IsType(t, WalletClient{}, *client)
// 	})
// }

// // TestSetAdminKey will test the admin key setter
// func TestSetAdminKey(t *testing.T) {
// 	t.Run("invalid", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		err := client.SetAdminKey("")
// 		assert.Error(t, err)
// 	})

// 	t.Run("valid", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		err := client.SetAdminKey(fixtures.XPrivString)
// 		assert.NoError(t, err)
// 	})

// 	t.Run("invalid with", func(t *testing.T) {
// 		_, err := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithAdminKey("rest"),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		assert.Error(t, err)
// 	})

// 	t.Run("valid with", func(t *testing.T) {
// 		_, err := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithAdminKey(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		assert.NoError(t, err)
// 	})
// }

// // TestSetSignRequest will test the sign request setter
// func TestSetSignRequest(t *testing.T) {
// 	t.Run("true", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		client.SetSignRequest(true)
// 		assert.True(t, client.IsSignRequest())
// 	})

// 	t.Run("false", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		client.SetSignRequest(false)
// 		assert.False(t, client.IsSignRequest())
// 	})

// 	t.Run("false by default", func(t *testing.T) {
// 		client, err := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		require.NoError(t, err)
// 		assert.False(t, client.IsSignRequest())
// 	})
// }

// // TestGetTransport will test the GetTransport method
// func TestGetTransport(t *testing.T) {
// 	t.Run("GetTransport", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 		)
// 		transport := client.GetTransport()
// 		assert.IsType(t, &transports.TransportHTTP{}, *transport)
// 	})

// 	t.Run("client GetTransport", func(t *testing.T) {
// 		client, _ := New(
// 			WithXPriv(fixtures.XPrivString),
// 			WithHTTP(fixtures.ServerURL),
// 			WithAdminKey(fixtures.XPrivString),
// 			WithSignRequest(false),
// 		)
// 		transport := client.GetTransport()
// 		assert.IsType(t, &transports.TransportHTTP{}, *transport)
// 	})
// }

// func TestAuthenticationWithOnlyAccessKey(t *testing.T) {
// 	anyConditions := make(map[string]interface{}, 0)
// 	var anyMetadataConditions *models.Metadata
// 	anyParam := "sth"

// 	testCases := []struct {
// 		caseTitle    string
// 		path         string
// 		clientMethod func(*WalletClient) (any, error)
// 	}{
// 		{
// 			caseTitle:    "GetXPub",
// 			path:         "/xpub",
// 			clientMethod: func(c *WalletClient) (any, error) { return c.GetXPub(context.Background()) },
// 		},
// 		{
// 			caseTitle:    "GetAccessKey",
// 			path:         "/access-key",
// 			clientMethod: func(c *WalletClient) (any, error) { return c.GetAccessKey(context.Background(), anyParam) },
// 		},
// 		{
// 			caseTitle: "GetAccessKeys",
// 			path:      "/access-key",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetAccessKeys(context.Background(), anyMetadataConditions)
// 			},
// 		},
// 		{
// 			caseTitle:    "GetDestinationByID",
// 			path:         "/destination",
// 			clientMethod: func(c *WalletClient) (any, error) { return c.GetDestinationByID(context.Background(), anyParam) },
// 		},
// 		{
// 			caseTitle: "GetDestinationByAddress",
// 			path:      "/destination",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetDestinationByAddress(context.Background(), anyParam)
// 			},
// 		},
// 		{
// 			caseTitle: "GetDestinationByLockingScript",
// 			path:      "/destination",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetDestinationByLockingScript(context.Background(), anyParam)
// 			},
// 		},
// 		{
// 			caseTitle: "GetDestinations",
// 			path:      "/destination/search",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetDestinations(context.Background(), nil)
// 			},
// 		},
// 		{
// 			caseTitle: "GetTransaction",
// 			path:      "/transaction",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetTransaction(context.Background(), fixtures.Transaction.ID)
// 			},
// 		},
// 		{
// 			caseTitle: "GetTransactions",
// 			path:      "/transaction/search",
// 			clientMethod: func(c *WalletClient) (any, error) {
// 				return c.GetTransactions(context.Background(), anyConditions, anyMetadataConditions, &transports.QueryParams{})
// 			},
// 		},
// 	}

// 	for _, test := range testCases {
// 		t.Run(test.caseTitle, func(t *testing.T) {
// 			transportHandler := testTransportHandler{
// 				Type: fixtures.RequestType,
// 				Queries: []*testTransportHandlerRequest{{
// 					Path: test.path,
// 					Result: func(w http.ResponseWriter, req *http.Request) {
// 						assertAuthHeaders(t, req)
// 						w.Header().Set("Content-Type", "application/json")
// 						mustWrite(w, "{}")
// 					},
// 				}},
// 				ClientURL: fixtures.ServerURL,
// 				Client:    WithHTTPClient,
// 			}

// 			client := getTestWalletClientWithOpts(transportHandler, WithAccessKey(fixtures.AccessKeyString))

// 			_, err := test.clientMethod(client)
// 			if err != nil {
// 				t.Log(err)
// 			}
// 		})
// 	}
// }

// func assertAuthHeaders(t *testing.T, req *http.Request) {
// 	assert.Empty(t, req.Header.Get("x-auth-xpub"), "Header value x-auth-xpub should be empty")
// 	assert.NotEmpty(t, req.Header.Get("x-auth-key"), "Header value x-auth-key should not be empty")
// 	assert.NotEmpty(t, req.Header.Get("x-auth-time"), "Header value x-auth-time should not be empty")
// 	assert.NotEmpty(t, req.Header.Get("x-auth-hash"), "Header value x-auth-hash should not be empty")
// 	assert.NotEmpty(t, req.Header.Get("x-auth-nonce"), "Header value x-auth-nonce should not be empty")
// 	assert.NotEmpty(t, req.Header.Get("x-auth-signature"), "Header value x-auth-signature should not be empty")
// }

// func getTestWalletClient(transportHandler testTransportHandler, adminKey bool) *WalletClient {
// 	opts := []ClientOps{
// 		WithXPriv(fixtures.XPrivString),
// 	}
// 	if adminKey {
// 		opts = append(opts, WithAdminKey(fixtures.XPrivString))
// 	}

// 	return getTestWalletClientWithOpts(transportHandler, opts...)
// }

// func getTestWalletClientWithOpts(transportHandler testTransportHandler, options ...ClientOps) *WalletClient {
// 	mux := http.NewServeMux()
// 	if transportHandler.Queries != nil {
// 		for _, query := range transportHandler.Queries {
// 			mux.HandleFunc(query.Path, query.Result)
// 		}
// 	} else {
// 		mux.HandleFunc(transportHandler.Path, func(w http.ResponseWriter, req *http.Request) {
// 			w.Header().Set("Content-Type", "application/json")
// 			mustWrite(w, transportHandler.Result)
// 		})
// 	}
// 	httpclient := &http.Client{Transport: localRoundTripper{handler: mux}}

// 	opts := []ClientOps{
// 		transportHandler.Client(transportHandler.ClientURL, httpclient),
// 	}

// 	opts = append(opts, options...)

// 	client, _ := New(opts...)

// 	return client
// }
