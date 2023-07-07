package transports

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	buxmodels "github.com/BuxOrg/bux-models"
	buxerrors "github.com/BuxOrg/bux-models/bux-errors"
	"github.com/BuxOrg/go-buxclient/utils"
	"github.com/libsv/go-bk/bip32"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	xPrivString = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	xPubString  = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
)

type TransportGraphQLMock struct {
	TransportGraphQL
	client *GraphQLMockClient
}

func (t *TransportGraphQLMock) Init() error {
	t.client = &GraphQLMockClient{}
	return nil
}

type GraphQLMockClient struct {
	Response interface{}
	Request  *graphql.Request
	Error    error
}

func (g *GraphQLMockClient) Run(_ context.Context, req *graphql.Request, resp interface{}) error {
	j, _ := json.Marshal(g.Response) //nolint:errchkjson // used for testing only
	_ = json.Unmarshal(j, &resp)
	g.Request = req
	return g.Error
}

// TestNewXpub will test the NewXpub method
func TestNewXpub(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	// xPub, _ := xPriv.Neuter()

	t.Run("no admin key", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{},
		}
		err := client.NewXpub(context.Background(), xPubString, nil)
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
		err := client.NewXpub(context.Background(), xPubString, nil)
		assert.ErrorIs(t, err, errTestTerror)
	})

	t.Run("return success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: buxmodels.Xpub{
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
		err := client.NewXpub(context.Background(), xPubString, nil)
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
		destination, err := client.NewDestination(context.Background(), nil)
		assert.ErrorIs(t, err, buxerrors.ErrMissingXPriv)
		assert.ErrorIs(t, err, nil)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DestinationData{
				Destination: &buxmodels.Destination{
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
		destination, err := client.NewDestination(context.Background(), nil)
		assert.NoError(t, err)
		assert.IsType(t, &buxmodels.Destination{}, destination)
		assert.Equal(t, "test-address", destination.Address)
		checkAuthHeaders(t, graphqlClient)
	})

	t.Run("no signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DestinationData{
				Destination: &buxmodels.Destination{
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
		destination, err := client.NewDestination(context.Background(), nil)
		assert.NoError(t, err)
		assert.IsType(t, &buxmodels.Destination{}, destination)
		assert.Equal(t, "test-address", destination.Address)
		assert.Len(t, graphqlClient.Request.Header, 1)
		assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Xpub")
	})
}

// TestDraftTransaction will test the DraftTransaction method
func TestDraftTransaction(t *testing.T) {
	xPriv, _ := bip32.NewKeyFromString(xPrivString)
	xPub, _ := xPriv.Neuter()
	config := &buxmodels.TransactionConfig{
		SendAllTo: &buxmodels.TransactionOutput{
			To: "bux@bux.org",
		},
	}

	t.Run("missing xpriv", func(t *testing.T) {
		client := TransportGraphQLMock{
			TransportGraphQL: TransportGraphQL{
				signRequest: true,
				client:      &GraphQLMockClient{},
			},
		}
		destination, err := client.DraftTransaction(context.Background(), config, nil)
		assert.ErrorIs(t, err, buxerrors.ErrMissingXPriv)
		assert.ErrorIs(t, err, nil)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DraftTransactionData{
				NewTransaction: &buxmodels.DraftTransaction{
					Status: buxmodels.DraftStatusDraft,
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
		draftTransaction, err := client.DraftTransaction(context.Background(), config, nil)
		assert.NoError(t, err)
		assert.IsType(t, &buxmodels.DraftTransaction{}, draftTransaction)
		checkAuthHeaders(t, graphqlClient)
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
		destination, err := client.DraftToRecipients(context.Background(), recipients, nil)
		assert.ErrorIs(t, err, buxerrors.ErrMissingXPriv)
		assert.ErrorIs(t, err, nil)
		assert.Nil(t, destination)
	})

	t.Run("signRequest success", func(t *testing.T) {
		graphqlClient := GraphQLMockClient{
			Response: DraftTransactionData{
				NewTransaction: &buxmodels.DraftTransaction{
					Status: buxmodels.DraftStatusDraft,
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
		draftTransaction, err := client.DraftToRecipients(context.Background(), recipients, nil)
		assert.NoError(t, err)
		assert.IsType(t, &buxmodels.DraftTransaction{}, draftTransaction)
		checkAuthHeaders(t, graphqlClient)
	})
}

func checkAuthHeaders(t *testing.T, graphqlClient GraphQLMockClient) {
	assert.Len(t, graphqlClient.Request.Header, 5)
	assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Hash")
	assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Nonce")
	assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Signature")
	assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Time")
	assert.Contains(t, graphqlClient.Request.Header, "Bux-Auth-Xpub")
}
