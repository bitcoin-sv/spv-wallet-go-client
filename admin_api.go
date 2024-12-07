package spvwallet

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/accesskeys"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/invitations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/paymails"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/stats"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/status"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/utxos"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/webhooks"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/xpubs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// AdminAPI provides a simplified interface for interacting with admin-related APIs.
// It abstracts the complexities of making HTTP requests and handling responses,
// allowing developers to easily interact with admin API endpoints.
//
// A zero-value AdminAPI is not usable. Use the NewAdminAPI function to create
// a properly initialized instance.
//
// Methods may return wrapped errors, including models.SPVError or
// ErrUnrecognizedAPIResponse, depending on the behavior of the SPV Wallet API.
type AdminAPI struct {
	xpubsAPI        *xpubs.API
	paymailsAPI     *paymails.API
	accessKeyAPI    *accesskeys.API
	transactionsAPI *transactions.API
	utxosAPI        *utxos.API
	contactsAPI     *contacts.API
	invitationsAPI  *invitations.API
	webhooksAPI     *webhooks.API
	statusAPI       *status.API
	statsAPI        *stats.API
}

// CreateXPub creates a new XPub record via the Admin XPubs API.
// The provided command contains the necessary parameters to define the XPub record.
//
// The API response is unmarshaled into a *response.Xpub struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) CreateXPub(ctx context.Context, cmd *commands.CreateUserXpub) (*response.Xpub, error) {
	res, err := a.xpubsAPI.CreateXPub(ctx, cmd)
	if err != nil {
		return nil, xpubs.HTTPErrorFormatter("create XPub", err).FormatPostErr()
	}

	return res, nil
}

// XPubs retrieves a paginated list of user XPubs via the Admin XPubs API.
// The response includes user XPubs along with pagination metadata, such as
// the current page number, sort order, and the field used for sorting (sortBy).
//
// Query parameters can be configured using optional query options. These options allow
// filtering based on metadata, pagination settings, or specific XPub attributes.
//
// The API response is unmarshaled into a *queries.XPubPage struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) XPubs(ctx context.Context, opts ...queries.XPubQueryOption) (*queries.XPubPage, error) {
	res, err := a.xpubsAPI.XPubs(ctx, opts...)
	if err != nil {
		return nil, xpubs.HTTPErrorFormatter("retrieve XPubs page", err).FormatGetErr()
	}

	return res, nil
}

// Contacts retrieves a paginated list of user contacts from the admin contacts API.
//
// The response includes contact data along with pagination details, such as the
// current page, sort order, and sortBy field. Optional query parameters can be
// provided using query options. The result is unmarshaled into a *queries.UserContactsPage.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) Contacts(ctx context.Context, opts ...queries.ContactQueryOption) (*queries.UserContactsPage, error) {
	res, err := a.contactsAPI.Contacts(ctx, opts...)
	if err != nil {
		return nil, contacts.HTTPErrorFormatter("retrieve user contacts page", err).FormatGetErr()
	}

	return res, nil
}

// ContactUpdate updates a user's contact information through the admin contacts API.
//
// This method uses the `UpdateContact` command to specify the details of the contact to update.
// It sends the update request to the API, unmarshals the response into a `*response.Contact`,
// and returns the updated contact. If the API request fails or the response cannot be decoded,
// an error is returned.
func (a *AdminAPI) ContactUpdate(ctx context.Context, cmd *commands.UpdateContact) (*response.Contact, error) {
	res, err := a.contactsAPI.UpdateContact(ctx, cmd)
	if err != nil {
		msg := fmt.Sprintf("update contact with ID: %s", cmd.ID)
		return nil, contacts.HTTPErrorFormatter(msg, err).FormatPutErr()
	}

	return res, nil
}

// DeleteContact deletes a user contact with the given ID via the admin contacts API.
// Returns an error if the API request fails or the response cannot be decoded.
// A nil error indicates the deleting contact was successful.
func (a *AdminAPI) DeleteContact(ctx context.Context, ID string) error {
	err := a.contactsAPI.DeleteContact(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("delete contact with ID: %s", ID)
		return contacts.HTTPErrorFormatter(msg, err).FormatDeleteErr()
	}

	return nil
}

// AcceptInvitation processes and accepts a user contact invitation using the given ID via the admin invitations API.
// Returns an error if the API request fails. A nil error indicates the invitation was successfully accepted.
func (a *AdminAPI) AcceptInvitation(ctx context.Context, ID string) error {
	err := a.invitationsAPI.AcceptInvitation(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("accept invitation with ID: %s", ID)
		return invitations.HTTPErrorFormatter(msg, err).FormatDeleteErr()
	}

	return nil
}

// RejectInvitation processes and rejects a user contact invitation using the given ID via the admin invitations API.
// Returns an error if the API request fails. A nil error indicates the invitation was successfully rejected.
func (a *AdminAPI) RejectInvitation(ctx context.Context, ID string) error {
	err := a.invitationsAPI.RejectInvitation(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("delete invitation with ID: %s", ID)
		return invitations.HTTPErrorFormatter(msg, err).FormatDeleteErr()
	}

	return nil
}

// Transactions retrieves a paginated list of transactions via the Admin transactions API.
// The returned response includes transactions and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to be to unmarshal into a *queries.TransactionPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (a *AdminAPI) Transactions(ctx context.Context, opts ...queries.TransactionsQueryOption) (*queries.TransactionPage, error) {
	res, err := a.transactionsAPI.Transactions(ctx, opts...)
	if err != nil {
		return nil, transactions.HTTPErrorFormatter("retrieve transactions page", err).FormatGetErr()
	}

	return res, nil
}

// Transaction retrieves a specific transaction by its ID via the Admin transactions API.
// The response is expected to be unmarshaled into a *response.Transaction struct.
// Returns an error if the request fails or the response cannot be decoded.
func (a *AdminAPI) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	res, err := a.transactionsAPI.Transaction(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("retrieve a transaction with ID: %s", ID)
		return nil, transactions.HTTPErrorFormatter(msg, err).FormatGetErr()
	}

	return res, nil
}

// AccessKeys retrieves a paginated list of access keys via the Admin XPubs API.
// The response includes access keys and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to unmarshal into a *queries.AccessKeyPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (a *AdminAPI) AccessKeys(ctx context.Context, accessKeyOpts ...queries.AdminAccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	res, err := a.accessKeyAPI.AccessKeys(ctx, accessKeyOpts...)
	if err != nil {
		return nil, accesskeys.HTTPErrorFormatter("retrieve access keys page ", err).FormatGetErr()
	}

	return res, nil
}

// SubscribeWebhook registers a webhook subscription using the Admin Webhooks API.
// The provided command contains the parameters required to define the webhook subscription.
// Accepts context for controlling cancellation and timeout for the API request.
// The CreateWebhookSubscription command includes the webhook URL and authentication details.
// Returns a formatted error if the API request fails. A nil error indicates the webhook subscription was successful.
func (a *AdminAPI) SubscribeWebhook(ctx context.Context, cmd *commands.CreateWebhookSubscription) error {
	err := a.webhooksAPI.SubscribeWebhook(ctx, cmd)
	if err != nil {
		msg := fmt.Sprintf("subscribe webhook URL address: %s", cmd.URL)
		return webhooks.HTTPErrorFormatter(msg, err).FormatPostErr()
	}

	return nil
}

// UnsubscribeWebhook removes a webhook subscription using the Admin Webhooks API.
// Accepts the context for controlling cancellation and timeout for the API request.
// CancelWebhookSubscription command specifies the webhook URL to be unsubscribed.
// Returns a formatted error if the API request fails. A nil error indicates the webhook subscription was successfully deleted.
func (a *AdminAPI) UnsubscribeWebhook(ctx context.Context, cmd *commands.CancelWebhookSubscription) error {
	err := a.webhooksAPI.UnsubscribeWebhook(ctx, cmd)
	if err != nil {
		msg := fmt.Sprintf("unsubscribe webhook URL address: %s", cmd.URL)
		return webhooks.HTTPErrorFormatter(msg, err).FormatDeleteErr()
	}

	return nil
}

// UTXOs fetches a paginated list of UTXOs via the Admin XPubs API.
// The response includes UTXOs along with pagination details, such as page number,
// sort order, and sorting field.
//
// Optional query parameters can be applied using the provided query options.
// The response is unmarshaled into a *queries.UtxosPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (a *AdminAPI) UTXOs(ctx context.Context, opts ...queries.AdminUtxoQueryOption) (*queries.UtxosPage, error) {
	res, err := a.utxosAPI.UTXOs(ctx, opts...)
	if err != nil {
		return nil, utxos.HTTPErrorFormatter("retrieve utxos page ", err).FormatGetErr()
	}

	return res, nil
}

// Paymails retrieves a paginated list of paymail addresses via the Admin Paymails API.
// The response includes user paymails along with pagination metadata, such as
// the current page number, sort order, and the field used for sorting (sortBy).
//
// Query parameters can be configured using optional query options. These options allow
// filtering based on metadata, pagination settings, or specific paymail attributes.
//
// The API response is unmarshaled into a *queries.PaymailAddressPage struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) Paymails(ctx context.Context, opts ...queries.PaymailQueryOption) (*queries.PaymailAddressPage, error) {
	res, err := a.paymailsAPI.Paymails(ctx, opts...)
	if err != nil {
		return nil, paymails.HTTPErrorFormatter("retrieve paymail addresses page", err).FormatGetErr()
	}

	return res, nil
}

// Paymail retrieves the paymail address associated with the specified ID via the Admin Paymails API.
// The response is expected to be unmarshaled into a *response.PaymailAddress struct.
// Returns an error if the request fails or the response cannot be decoded.
func (a *AdminAPI) Paymail(ctx context.Context, ID string) (*response.PaymailAddress, error) {
	res, err := a.paymailsAPI.Paymail(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("retrieve paymail address with ID: %s", ID)
		return nil, paymails.HTTPErrorFormatter(msg, err).FormatGetErr()
	}

	return res, nil
}

// CreatePaymail creates a new paymail address record via the Admin Paymails API.
// The provided command contains the necessary parameters to define the paymail address record.
//
// The API response is unmarshaled into a *response.Xpub PaymailAddress.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) CreatePaymail(ctx context.Context, cmd *commands.CreatePaymail) (*response.PaymailAddress, error) {
	res, err := a.paymailsAPI.CreatePaymail(ctx, cmd)
	if err != nil {
		return nil, paymails.HTTPErrorFormatter("create paymail address", err).FormatPostErr()
	}

	return res, nil
}

// DeletePaymail deletes a paymail address with via the Admin Paymails API.
// It returns an error if the API request fails. A nil error indicates that the paymail
// was successfully deleted.
func (a *AdminAPI) DeletePaymail(ctx context.Context, address string) error {
	err := a.paymailsAPI.DeletePaymail(ctx, address)
	if err != nil {
		msg := fmt.Sprintf("remove paymail address: %s", address)
		return paymails.HTTPErrorFormatter(msg, err).FormatGetErr()
	}

	return nil
}

// Stats retrieves information about the login status via the Admin XPubs API.
// It accepts a context for controlling cancellation and timeout of the API request.
// The response is expected to be unmarshaled into a *models.AdminStats struct.
// A nil error with a valid response indicates the request was successful.
// Returns a formatted error if the API request fails.
func (a *AdminAPI) Stats(ctx context.Context) (*models.AdminStats, error) {
	res, err := a.statsAPI.Stats(ctx)
	if err != nil {
		return nil, stats.HTTPErrorFormatter("retrieve stats", err).FormatGetErr()
	}

	return res, nil
}

// Status retrieves information about the key type used during the authentication phase.
// If the key corresponds to the admin key, the method returns true with a nil error.
// Otherwise, it returns false with a nil error, indicating that the key used does not match
// the SPV Wallet API admin key. A non-nil error is returned if the API request fails.
func (a *AdminAPI) Status(ctx context.Context) (bool, error) {
	ok, err := a.statusAPI.Status(ctx)
	if err != nil {
		return false, status.HTTPErrorFormatter("retrieve information about the used key type: %w", err).FormatGetErr()
	}

	return ok, nil
}

// NewAdminAPIWithXPriv initializes a new AdminAPI instance using an extended private key (xPriv).
// This function configures the API client with the provided configuration and uses the xPriv key for authentication.
// If any step fails, an appropriate error is returned.
//
// Note: Requests made with this instance will be securely signed.
func NewAdminAPIWithXPriv(cfg config.Config, xPriv string) (*AdminAPI, error) {
	authenticator, err := auth.NewXprivAuthenticator(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xPriv authenticator: %w", err)
	}

	return initAdminAPI(cfg, authenticator)
}

// NewAdminWithXPub initializes a new AdminAPI instance using an extended public key (xPub).
// This function configures the API client with the provided configuration and uses the xPub key for authentication.
// If any configuration or initialization step fails, an appropriate error is returned.
//
// Note: Requests made with this instance will not be signed.
// For enhanced security, it is strongly recommended to use `NewAdminAPIWithXPriv` instead.
func NewAdminWithXPub(cfg config.Config, xPub string) (*AdminAPI, error) {
	authenticator, err := auth.NewXpubOnlyAuthenticator(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xPub authenticator: %w", err)
	}

	return initAdminAPI(cfg, authenticator)
}

func initAdminAPI(cfg config.Config, auth authenticator) (*AdminAPI, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr to url.URL: %w", err)
	}

	httpClient := restyutil.NewHTTPClient(cfg, auth)
	if httpClient == nil {
		return nil, fmt.Errorf("failed to initialize HTTP client - nil value.")
	}

	return &AdminAPI{
		paymailsAPI:     paymails.NewAPI(url, httpClient),
		transactionsAPI: transactions.NewAPI(url, httpClient),
		xpubsAPI:        xpubs.NewAPI(url, httpClient),
		utxosAPI:        utxos.NewAPI(url, httpClient),
		accessKeyAPI:    accesskeys.NewAPI(url, httpClient),
		webhooksAPI:     webhooks.NewAPI(url, httpClient),
		contactsAPI:     contacts.NewAPI(url, httpClient),
		invitationsAPI:  invitations.NewAPI(url, httpClient),
		statusAPI:       status.NewAPI(url, httpClient),
		statsAPI:        stats.NewAPI(url, httpClient),
	}, nil
}
