package transports

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/BitcoinSchema/xapi/bux"
	"github.com/BitcoinSchema/xapi/bux/authentication"
	"github.com/BitcoinSchema/xapi/bux/utils"
	"github.com/libsv/go-bk/bip32"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	xPrivString = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	xPubString  = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	serverURL   = "https://example.com/"
)

// TransportGraphQLMock ...
type TransportGraphQLMock struct {
	TransportGraphQL
	client *GraphQLMockClient
}

// Init() ...
func (t *TransportGraphQLMock) Init() error {
	t.client = &GraphQLMockClient{}
	return nil
}

// GraphQLMockClient ...
type GraphQLMockClient struct {
	Response interface{}
	Request  *graphql.Request
	Error    error
}

// Run ...
func (g *GraphQLMockClient) Run(_ context.Context, req *graphql.Request, resp interface{}) error {
	j, _ := json.Marshal(g.Response)
	_ = json.Unmarshal(j, &resp)
	g.Request = req
	return g.Error
}

// TestRegisterXpub will test the RegisterXpub method
func TestRegisterXpub(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	//xPub, _ := xPriv.Neuter()

	t.Run("no admin key", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{},
		}
		err := client.RegisterXpub(context.Background(), xPubString)
		assert.ErrorIs(t, err, ErrAdminKey)
	})

	t.Run("return error", func(t *testing.T) {
		errTestTerror := errors.New("test error")
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				adminXPriv: xPriv,
				client: &GraphQLMockClient{
					Error: errTestTerror,
				},
			},
		}
		err := client.RegisterXpub(context.Background(), xPubString)
		assert.ErrorIs(t, err, errTestTerror)
	})

	t.Run("return success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: bux.Xpub{
				ID:              utils.Hash(xPubString),
				CurrentBalance:  0,
				NextInternalNum: 0,
				NextExternalNum: 0,
			},
		}
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				adminXPriv: xPriv,
				client:     &graphqlClient,
			},
		}
		err := client.RegisterXpub(context.Background(), xPubString)
		assert.NoError(t, err)
	})
}

// TestGetDestination will test the GetDestination method
func TestGetDestination(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	xPub, _ := xPriv.Neuter()

	t.Run("missing xpriv", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				signRequest: true,
				client:      &GraphQLMockClient{},
			},
		}
		destination, err := client.GetDestination(context.Background())
		assert.ErrorIs(t, err, authentication.ErrMissingXPriv)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DestinationData{
				Destination: &bux.Destination{
					Address: "test-address",
				},
			},
		}
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				xPriv:       xPriv,
				xPub:        xPub,
				signRequest: true,
				client:      &graphqlClient,
			},
		}
		destination, err := client.GetDestination(context.Background())
		assert.NoError(t, err)
		assert.IsType(t, &bux.Destination{}, destination)
		assert.Equal(t, "test-address", destination.Address)
		assert.Len(t, graphqlClient.Request.Header, 5)
		assert.Contains(t, graphqlClient.Request.Header, "Auth_hash")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_nonce")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_signature")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_time")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_xpub")
	})

	t.Run("no signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DestinationData{
				Destination: &bux.Destination{
					Address: "test-address",
				},
			},
		}
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				xPriv:  xPriv,
				xPub:   xPub,
				client: &graphqlClient,
			},
		}
		destination, err := client.GetDestination(context.Background())
		assert.NoError(t, err)
		assert.IsType(t, &bux.Destination{}, destination)
		assert.Equal(t, "test-address", destination.Address)
		assert.Len(t, graphqlClient.Request.Header, 1)
		assert.Contains(t, graphqlClient.Request.Header, "Auth_xpub")
	})
}

// TestDraftTransaction will test the DraftTransaction method
func TestDraftTransaction(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	xPub, _ := xPriv.Neuter()
	config := &bux.TransactionConfig{
		SendAllTo: "bux@bux.org",
	}

	t.Run("missing xpriv", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				signRequest: true,
				client:      &GraphQLMockClient{},
			},
		}
		destination, err := client.DraftTransaction(context.Background(), config)
		assert.ErrorIs(t, err, authentication.ErrMissingXPriv)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DraftTransactionData{
				NewTransaction: &bux.DraftTransaction{
					Status: bux.DraftStatusDraft,
				},
			},
		}
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				xPriv:       xPriv,
				xPub:        xPub,
				client:      &graphqlClient,
				signRequest: true,
			},
		}
		draftTransaction, err := client.DraftTransaction(context.Background(), config)
		assert.NoError(t, err)
		assert.IsType(t, &bux.DraftTransaction{}, draftTransaction)
		assert.Len(t, graphqlClient.Request.Header, 5)
		assert.Contains(t, graphqlClient.Request.Header, "Auth_hash")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_nonce")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_signature")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_time")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_xpub")
	})
}

// TestDraftToRecipients will test the DraftToRecipients method
func TestDraftToRecipients(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	xPub, _ := xPriv.Neuter()
	recipients := []*Recipients{{
		To:       "bux@bux.org",
		Satoshis: 12125,
	}}

	t.Run("missing xpriv", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				signRequest: true,
				client:      &GraphQLMockClient{},
			},
		}
		destination, err := client.DraftToRecipients(context.Background(), recipients)
		assert.ErrorIs(t, err, authentication.ErrMissingXPriv)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DraftTransactionData{
				NewTransaction: &bux.DraftTransaction{
					Status: bux.DraftStatusDraft,
				},
			},
		}
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				xPriv:       xPriv,
				xPub:        xPub,
				client:      &graphqlClient,
				signRequest: true,
			},
		}
		draftTransaction, err := client.DraftToRecipients(context.Background(), recipients)
		assert.NoError(t, err)
		assert.IsType(t, &bux.DraftTransaction{}, draftTransaction)
		assert.Len(t, graphqlClient.Request.Header, 5)
		assert.Contains(t, graphqlClient.Request.Header, "Auth_hash")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_nonce")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_signature")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_time")
		assert.Contains(t, graphqlClient.Request.Header, "Auth_xpub")
	})
}
