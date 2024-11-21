package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/invitations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/users"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

// Config holds configuration settings for establishing a connection and handling
// request details in the application.
type Config struct {
	Addr      string            // The base address of the SPV Wallet API.
	Timeout   time.Duration     // The HTTP requests timeout duration.
	Transport http.RoundTripper // Custom HTTP transport, allowing optional customization of the HTTP client behavior.
}

// NewDefaultConfig returns a default configuration for connecting to the SPV Wallet API,
// setting a one-minute timeout, using the default HTTP transport, and applying the
// base API address as the addr value.
func NewDefaultConfig(addr string) Config {
	return Config{
		Addr:      addr,
		Timeout:   1 * time.Minute,
		Transport: http.DefaultTransport,
	}
}

// Client provides methods for user-related and admin-related APIs.
// This struct is designed to abstract and simplify the process of making HTTP calls
// to the relevant endpoints. By utilizing this Client struct, developers can easily
// interact with both user and admin APIs without needing to manage the details
// of the HTTP requests and responses directly.
type Client struct {
	xpubAPI         *users.XPubAPI
	accessKeyAPI    *users.AccessKeyAPI
	configsAPI      *configs.API
	merkleRootsAPI  *merkleroots.API
	contactsAPI     *contacts.API
	invitationsAPI  *invitations.API
	transactionsAPI *transactions.API
	utxosAPI        *utxos.API

	totp *totp.Client //only available when using xPriv
}

// NewWithXPub creates a new client instance using an extended public key (xPub).
// Requests made with this instance will not be signed, that's why we strongly recommend to use `WithXPriv` or `WithAccessKey` option instead.
func NewWithXPub(cfg Config, xPub string) (*Client, error) {
	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xpub authenticator: %w", err)
	}
	client, err := newClient(cfg, authenticator)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}
	return client, nil
}

// NewWithXPriv creates a new client instance using an extended private key (xPriv).
// Generates an HD key from the provided xPriv and sets up the client instance to sign requests
// by setting the SignRequest flag to true. The generated HD key can be used for secure communications.
func NewWithXPriv(cfg Config, xPriv string) (*Client, error) {
	key, err := bip32.GenerateHDKeyFromString(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xpriv: %w", err)
	}

	authenticator, err := auth.NewXprivAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xpriv authenticator: %w", err)
	}

	client, err := newClient(cfg, authenticator)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	client.totp = totp.New(key)

	return client, nil
}

// NewWithAccessKey creates a new client instance using an access key.
// Function attempts to convert the provided access key from either hex or WIF format
// to a PrivateKey. The resulting PrivateKey is used to sign requests made by the client instance
// by setting the SignRequest flag to true.
func NewWithAccessKey(cfg Config, accessKey string) (*Client, error) {
	key, err := privateKeyFromHexOrWIF(accessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to return private key from hex or WIF: %w", err)
	}

	authenticator, err := auth.NewAccessKeyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized access key authenticator: %w", err)
	}

	return newClient(cfg, authenticator)
}

// Contacts retrieves a paginated list of user contacts from the user contacts API.
// The API response includes user contacts along with pagination details, such as
// the current page number, sort order, and the field used for sorting (sortBy).
//
// Optional query parameters can be provided via query options. The response is
// unmarshaled into a *queries.UserContactsPage struct. If the API request fails
// or the response cannot be decoded, an error is returned.
func (c *Client) Contacts(ctx context.Context, contactOpts ...queries.ContactQueryOption) (*queries.UserContactsPage, error) {
	res, err := c.contactsAPI.Contacts(ctx, contactOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contacts from the user contacts API: %w", err)
	}

	return res, nil
}

// ContactWithPaymail retrieves a specific user contact by their paymail address.
// The response is unmarshaled into a *response.Contact struct. If the API request
// fails or the response cannot be decoded, an error is returned.
func (c *Client) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	res, err := c.contactsAPI.ContactWithPaymail(ctx, paymail)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contact by paymail from the user contacts API: %w", err)
	}

	return res, nil
}

// UpsertContact adds or updates a user contact through the user contacts API.
// The response is unmarshaled into a *response.Contact struct. If the API request
// fails or the response cannot be decoded, an error is returned.
func (c *Client) UpsertContact(ctx context.Context, cmd commands.UpsertContact) (*response.Contact, error) {
	res, err := c.contactsAPI.UpsertContact(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert contact using the user contacts API: %w", err)
	}

	return res, nil
}

// RemoveContact deletes a user contact using the user contacts API.
// If the API request fails, an error is returned.
func (c *Client) RemoveContact(ctx context.Context, paymail string) error {
	err := c.contactsAPI.RemoveContact(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to remove contact using the user contacts API: %w", err)
	}

	return nil
}

// ConfirmContact checks the TOTP code and if it's ok, confirms user's contact using the user contacts API.
func (c *Client) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	if err := c.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits); err != nil {
		return fmt.Errorf("failed to validate TOTP for contact: %w", err)
	}

	err := c.contactsAPI.ConfirmContact(ctx, contact.Paymail)
	if err != nil {
		return fmt.Errorf("failed to confirm contact using the user contacts API: %w", err)
	}

	return nil
}

// UnconfirmContact unconfirms a user contact using the user contacts API.
// If the API request fails, an error is returned.
func (c *Client) UnconfirmContact(ctx context.Context, paymail string) error {
	err := c.contactsAPI.UnconfirmContact(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to unconfirm contact using the user contacts API: %w", err)
	}

	return nil
}

// AcceptInvitation accepts a contact invitation using the user invitations API.
// If the API request fails, an error is returned.
func (c *Client) AcceptInvitation(ctx context.Context, paymail string) error {
	err := c.invitationsAPI.AcceptInvitation(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to accept invitation using the user invitations API: %w", err)
	}

	return nil
}

// RejectInvitation rejects a contact invitation using the user invitations API.
// If the API request fails, an error is returned.
func (c *Client) RejectInvitation(ctx context.Context, paymail string) error {
	err := c.invitationsAPI.RejectInvitation(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to reject invitation using the user invitations API: %w", err)
	}

	return nil
}

// SharedConfig retrieves the shared configuration from the user configurations API.
// This method constructs an HTTP GET request to the "api/v1/configs/shared" endpoint and expects
// a response that can be unmarshaled into the response.SharedConfig struct. If the request fails
// or the response cannot be decoded, an error will be returned.
func (c *Client) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	res, err := c.configsAPI.SharedConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shared configuration from user configs API: %w", err)
	}

	return res, nil
}

// DraftTransaction creates a new draft transaction using the user transactions API.
// This method sends an HTTP POST request to the "/draft" endpoint and expects
// a response that can be unmarshaled into a response.DraftTransaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) DraftTransaction(ctx context.Context, cmd *commands.DraftTransaction) (*response.DraftTransaction, error) {
	res, err := c.transactionsAPI.DraftTransaction(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create a draft transaction by calling the user transactions API: %w", err)
	}

	return res, nil
}

// RecordTransaction submits a transaction for recording using the user transactions API.
// This method sends an HTTP POST request to the "/transactions" endpoint, expecting
// a response that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) RecordTransaction(ctx context.Context, cmd *commands.RecordTransaction) (*response.Transaction, error) {
	res, err := c.transactionsAPI.RecordTransaction(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to record a transaction with reference ID: %s by calling the user transactions API: %w", cmd.ReferenceID, err)
	}

	return res, nil
}

// UpdateTransactionMetadata updates the metadata of a transaction using the user transactions API.
// This method sends an HTTP PATCH request with updated metadata and expects a response
// that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) UpdateTransactionMetadata(ctx context.Context, cmd *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	res, err := c.transactionsAPI.UpdateTransactionMetadata(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to update a transaction metadata by calling the user user transactions API: %w", err)
	}

	return res, nil
}

// Transactions retrieves a paginated list of transactions from the user transactions API.
// The returned response includes transactions and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to unmarshal into a *response.PageModel[response.Transaction] struct.
// If the API request fails or the response cannot be decoded successfully, an error is returned.
func (c *Client) Transactions(ctx context.Context, opts ...queries.TransactionsQueryOption) (*queries.TransactionPage, error) {
	res, err := c.transactionsAPI.Transactions(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions page from the user transactions API: %w", err)
	}

	return res, nil
}

// Transaction retrieves a specific transaction by its ID using the user transactions API.
// This method expects a response that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	res, err := c.transactionsAPI.Transaction(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction with ID: %s from the user transactions API: %w", ID, err)
	}

	return res, nil
}

// XPub retrieves the complete xpub information for the current user.
// The server's response is expected to be unmarshaled into a *response.Xpub struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) XPub(ctx context.Context) (*response.Xpub, error) {
	res, err := c.xpubAPI.XPub(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve xpub information from the users API: %w", err)
	}

	return res, nil
}

// UpdateXPubMetadata updates the metadata associated with the current user's xpub.
// The server's response is expected to be unmarshaled into a *response.Xpub struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) UpdateXPubMetadata(ctx context.Context, cmd *commands.UpdateXPubMetadata) (*response.Xpub, error) {
	res, err := c.xpubAPI.UpdateXPubMetadata(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to update xpub metadata using the users API: %w", err)
	}

	return res, nil
}

// GenerateAccessKey creates a new access key associated with the current user's xpub.
// The server's response is expected to be unmarshaled into a *response.AccessKey struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) GenerateAccessKey(ctx context.Context, cmd *commands.GenerateAccessKey) (*response.AccessKey, error) {
	res, err := c.accessKeyAPI.GenerateAccessKey(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access key using the user access key API: %w", err)
	}

	return res, nil
}

// AccessKeys retrieves a paginated list of access keys from the user access keys API.
// The response includes access keys and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to unmarshal into a *queries.AccessKeyPage struct.
// If the API request fails or the response cannot be decoded successfully, an error is returned.
func (c *Client) AccessKeys(ctx context.Context, accessKeyOpts ...queries.AccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	res, err := c.accessKeyAPI.AccessKeys(ctx, accessKeyOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access keys page from the user access key API: %w", err)
	}

	return res, nil
}

// AccessKey retrieves the access key associated with the specified ID.
// The server's response is expected to be unmarshaled into a *response.AccessKey struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) AccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	res, err := c.accessKeyAPI.AccessKey(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access key using the user access key API: %w", err)
	}

	return res, nil
}

// RevokeAccessKey revokes the access key associated with the given ID.
// If the request fails or the response cannot be processed, an error is returned.
func (c *Client) RevokeAccessKey(ctx context.Context, ID string) error {
	err := c.accessKeyAPI.RevokeAccessKey(ctx, ID)
	if err != nil {
		return fmt.Errorf("failed to revoke access key using the users API: %w", err)
	}

	return nil
}

// UTXOs fetches a paginated list of UTXOs from the user UTXOs API.
// The response includes UTXOs along with pagination details, such as page number,
// sort order, and sorting field.
//
// Optional query parameters can be applied using the provided query options.
// The response is unmarshaled into a *queries.UtxosPage struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (c *Client) UTXOs(ctx context.Context, opts ...queries.UtxoQueryOption) (*queries.UtxosPage, error) {
	res, err := c.utxosAPI.UTXOs(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve UTXOs page from the user UTXOs API: %w", err)
	}

	return res, nil
}

// MerkleRoots retrieves a paginated list of Merkle roots from the user Merkle roots API.
// The API response includes Merkle roots along with pagination details, such as the current
// page number, sort order, and sorting field (sortBy).
//
// This method supports optional query parameters, which can be specified using the provided
// query options. These options customize the behavior of the API request, such as setting
// batch size or applying filters for pagination.
//
// The response is unmarshaled into a *queries.MerkleRootPage struct. If the API request fails
// or the response cannot be successfully decoded, an error is returned.
func (c *Client) MerkleRoots(ctx context.Context, opts ...queries.MerkleRootsQueryOption) (*queries.MerkleRootPage, error) {
	res, err := c.merkleRootsAPI.MerkleRoots(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Merkle roots from the API: %w", err)
	}

	return res, nil
}

// GenerateTotpForContact generates a TOTP code for the specified contact.
func (c *Client) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	if c.totp == nil {
		return "", errors.New("totp client not initialized - xPriv authentication required")
	}
	totp, err := c.totp.GenerateTotpForContact(contact, period, digits)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP for contact: %w", err)
	}
	return totp, nil
}

// ValidateTotpForContact validates a TOTP code for the specified contact.
func (c *Client) ValidateTotpForContact(contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	if c.totp == nil {
		return errors.New("totp client not initialized - xPriv authentication required")
	}
	if err := c.totp.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits); err != nil {
		return fmt.Errorf("failed to validate TOTP for contact: %w", err)
	}
	return nil
}

func privateKeyFromHexOrWIF(s string) (*ec.PrivateKey, error) {
	pk, err1 := ec.PrivateKeyFromWif(s)
	if err1 == nil {
		return pk, nil
	}

	pk, err2 := ec.PrivateKeyFromHex(s)
	if err2 != nil {
		return nil, errors.Join(err1, err2)
	}

	return pk, nil
}

type authenticator interface {
	Authenticate(r *resty.Request) error
}

func newClient(cfg Config, auth authenticator) (*Client, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr to url.URL: %w", err)
	}

	httpClient := newRestyClient(cfg, auth)
	return &Client{
		merkleRootsAPI:  merkleroots.NewAPI(url, httpClient),
		configsAPI:      configs.NewAPI(url, httpClient),
		transactionsAPI: transactions.NewAPI(url, httpClient),
		utxosAPI:        utxos.NewAPI(url, httpClient),
		accessKeyAPI:    users.NewAccessKeyAPI(url, httpClient),
		xpubAPI:         users.NewXPubAPI(url, httpClient),
		contactsAPI:     contacts.NewAPI(url, httpClient),
		invitationsAPI:  invitations.NewAPI(url, httpClient),
	}, nil
}

func newRestyClient(cfg Config, auth authenticator) *resty.Client {
	return resty.New().
		SetTransport(cfg.Transport).
		SetBaseURL(cfg.Addr).
		SetTimeout(cfg.Timeout).
		OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
			return auth.Authenticate(r)
		}).
		SetError(&models.SPVError{}).
		OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
			if r.IsSuccess() {
				return nil
			}

			if spvError, ok := r.Error().(*models.SPVError); ok && len(spvError.Code) > 0 {
				return spvError
			}

			return fmt.Errorf("%w: %s", goclienterr.ErrUnrecognizedAPIResponse, r.Body())
		})
}
