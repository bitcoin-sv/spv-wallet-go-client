package spvwallet

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/configs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/accesskeys"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/invitations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/paymails"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/utxos"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/xpubs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/constants"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

// UserAPI provides methods for interacting with user-related APIs.
// It abstracts the details of HTTP request and response handling,
// simplifying interaction with the endpoints.
//
// A zero-value UserAPI is not usable. Use one of the constructors
// (e.g., NewUserAPIWithAccessKey, NewUserAPIWithXPriv, or NewUserAPIWithXPub)
// to create a properly initialized instance.
//
// UserAPI methods may return wrapped errors, including models.SPVError or
// ErrUnrecognizedAPIResponse, depending on the behavior of the SPV Wallet API.
type UserAPI struct {
	xpubAPI         *xpubs.API
	accessKeyAPI    *accesskeys.API
	configsAPI      *configs.API
	merkleRootsAPI  *merkleroots.API
	contactsAPI     *contacts.API
	invitationsAPI  *invitations.API
	transactionsAPI *transactions.API
	utxosAPI        *utxos.API
	paymailsAPI     *paymails.API
	totpAPI         *totp.API //only available when using xPriv
}

// Contacts retrieves a paginated list of user contacts from the user contacts API.
//
// The response includes contact data along with pagination details, such as the
// current page, sort order, and sortBy field. Optional query parameters can be
// provided using query options. The result is unmarshaled into a *queries.ContactsPage.
// Returns an error if the API request fails or the response cannot be decoded.
func (u *UserAPI) Contacts(ctx context.Context, contactOpts ...queries.QueryOption[filter.ContactFilter]) (*queries.ContactsPage, error) {
	res, err := u.contactsAPI.Contacts(ctx, contactOpts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserContactsAPI, "retrieve contacts page", err).FormatGetErr()
	}

	return res, nil
}

// ContactWithPaymail retrieves a user contact by their paymail address.
// The response is unmarshaled into a *response.Contact.
// Returns an error if the API request fails or the response cannot be decoded.
func (u *UserAPI) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	res, err := u.contactsAPI.ContactWithPaymail(ctx, paymail)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserContactsAPI, "retrieve contact with paymail", err).FormatGetErr()
	}

	return res, nil
}

// UpsertContact adds or updates a user contact via the user contacts API.
// The response is unmarshaled into a *response.Contact.
// Returns an error if the API request fails or the response cannot be decoded.
func (u *UserAPI) UpsertContact(ctx context.Context, cmd commands.UpsertContact) (*response.Contact, error) {
	res, err := u.contactsAPI.UpsertContact(ctx, cmd)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserContactsAPI, "upsert contact", err).FormatPutErr()
	}

	return res, nil
}

// RemoveContact deletes a user contact with the given paymail via the user contacts API.
// Returns an error if the API request fails or the response cannot be decoded.
// A nil error indicates the deleting contact was successful.
func (u *UserAPI) RemoveContact(ctx context.Context, paymail string) error {
	err := u.contactsAPI.RemoveContact(ctx, paymail)
	if err != nil {
		return errutil.NewHTTPErrorFormatter(constants.AdminContactsAPI, "remove contact", err).FormatDeleteErr()
	}

	return nil
}

// ConfirmContact checks the TOTP code and if it's ok, confirms user's contact using the user contacts API.
func (u *UserAPI) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	if err := u.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits); err != nil {
		return fmt.Errorf("failed to validate TOTP for contact: %w", err)
	}

	err := u.contactsAPI.ConfirmContact(ctx, contact.Paymail)
	if err != nil {
		return errutil.NewHTTPErrorFormatter(constants.AdminContactsAPI, "confirm contact", err).FormatPostErr()
	}

	return nil
}

// UnconfirmContact unconfirms a user contact with the given paymail via the user contacts API.
// Returns an error if the API request fails or the response cannot be decoded. A nil error indicates the deleting confirmation was successful.
func (u *UserAPI) UnconfirmContact(ctx context.Context, paymail string) error {
	err := u.contactsAPI.UnconfirmContact(ctx, paymail)
	if err != nil {
		return errutil.NewHTTPErrorFormatter(constants.AdminContactsAPI, "unconfirm contact", err).FormatDeleteErr()
	}

	return nil
}

// AcceptInvitation accepts a user contact with the given paymail via the user contacts API.
// Returns an error if the API request fails or the response cannot be decoded. A nil error indicates the acceptation was successful.
func (u *UserAPI) AcceptInvitation(ctx context.Context, paymail string) error {
	err := u.invitationsAPI.AcceptInvitation(ctx, paymail)
	if err != nil {
		return errutil.NewHTTPErrorFormatter(constants.UserInvitationsAPI, "accept invitation", err).FormatPostErr()
	}

	return nil
}

// RejectInvitation rejects a user contact with the given paymail via the user contacts API.
// Returns an error if the API request fails or the response cannot be decoded.
// A nil error indicates the rejection was successful.
func (u *UserAPI) RejectInvitation(ctx context.Context, paymail string) error {
	err := u.invitationsAPI.RejectInvitation(ctx, paymail)
	if err != nil {
		return errutil.NewHTTPErrorFormatter(constants.UserInvitationsAPI, "reject invitation", err).FormatDeleteErr()
	}

	return nil
}

// SharedConfig retrieves the shared configuration via the configurations API.
// The response is unmarshaled into a response.SharedConfig.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	res, err := u.configsAPI.SharedConfig(ctx)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserSharedConfigAPI, "retrieve shared configuration", err).FormatGetErr()
	}

	return res, nil
}

// DraftTransaction creates a new draft transaction using the user transactions API.
// The response is expected to be unmarshaled into a *response.DraftTransaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (u *UserAPI) DraftTransaction(ctx context.Context, cmd *commands.DraftTransaction) (*response.DraftTransaction, error) {
	res, err := u.transactionsAPI.DraftTransaction(ctx, cmd)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, "create a draft transaction", err).FormatPostErr()
	}

	return res, nil
}

// RecordTransaction submits a transaction for recording via the user transactions API.
// The response is unmarshaled into a *response.Transaction.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) RecordTransaction(ctx context.Context, cmd *commands.RecordTransaction) (*response.Transaction, error) {
	res, err := u.transactionsAPI.RecordTransaction(ctx, cmd)
	if err != nil {
		msg := fmt.Sprintf("record a transaction with reference ID: %s", cmd.ReferenceID)
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, msg, err).FormatPostErr()
	}

	return res, nil
}

// UpdateTransactionMetadata updates the metadata of a transaction via the user transactions API.
// The response is expected to be unmarshaled into a *response.Transaction struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) UpdateTransactionMetadata(ctx context.Context, cmd *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	res, err := u.transactionsAPI.UpdateTransactionMetadata(ctx, cmd)
	if err != nil {
		msg := fmt.Sprintf("record a transaction with ID: %s", cmd.ID)
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, msg, err).FormatPutErr()
	}

	return res, nil
}

// Transactions retrieves a paginated list of transactions via the user transactions API.
// The returned response includes transactions and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to be to unmarshal into a *response.PageModel[response.Transaction] struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) Transactions(ctx context.Context, opts ...queries.QueryOption[filter.TransactionFilter]) (*queries.TransactionPage, error) {
	res, err := u.transactionsAPI.Transactions(ctx, opts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, "retrieve transactions page", err).FormatGetErr()
	}

	return res, nil
}

// Transaction retrieves a specific transaction by its ID via the user transactions API.
// The response is expected to be unmarshaled into a *response.Transaction struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	res, err := u.transactionsAPI.Transaction(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("retrieve a transaction with ID: %s", ID)
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, msg, err).FormatGetErr()
	}

	return res, nil
}

// FinalizeTransaction finalizes a draft transaction and returns its signed hex representation.
// It uses the draft transaction details to construct, enrich, and sign the transaction
// through the `transactionsigner.TransactionSignedHex` utility function.
// The response is the signed transaction in hex format.
// Returns an error if the transaction cannot be finalized.
func (u *UserAPI) FinalizeTransaction(draft *response.DraftTransaction) (string, error) {
	res, err := u.transactionsAPI.FinalizeTransaction(draft)
	if err != nil {
		return "", fmt.Errorf("couldn't finalize transaction with ID: %s, %w", draft.ID, err)
	}

	return res, nil
}

// SendToRecipients creates, finalizes, and broadcasts a transaction to multiple recipients.
// This method handles the complete process of drafting, finalizing, and recording the transaction
// using the recipient details provided in the command.
// The response is unmarshalled into a *response.Transaction struct.
// Returns an error if the transaction fails at any step, such as drafting, finalization or recording.
func (u *UserAPI) SendToRecipients(ctx context.Context, cmd *commands.SendToRecipients) (*response.Transaction, error) {
	res, err := u.transactionsAPI.SendToRecipients(ctx, cmd)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserTransactionsAPI, "send to recipients", err).FormatPostErr()
	}

	return res, nil
}

// XPub retrieves the full xpub information for the current user via the users API.
// The response is unmarshaled into a *response.Xpub.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) XPub(ctx context.Context) (*response.Xpub, error) {
	res, err := u.xpubAPI.XPub(ctx)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserXPubsAPI, "retrieve xpub information", err).FormatGetErr()
	}

	return res, nil
}

// UpdateXPubMetadata updates the metadata associated with the current user's xpub via the users API.
// The response is unmarshaled into a *response.Xpub.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) UpdateXPubMetadata(ctx context.Context, cmd *commands.UpdateXPubMetadata) (*response.Xpub, error) {
	res, err := u.xpubAPI.UpdateXPubMetadata(ctx, cmd)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserXPubsAPI, "update xpub metadata ", err).FormatGetErr()
	}

	return res, nil
}

// GenerateAccessKey creates a new access key associated with the current user's xpub via the users access key API.
// The response is unmarshaled into a *response.AccessKey.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) GenerateAccessKey(ctx context.Context, cmd *commands.GenerateAccessKey) (*response.AccessKey, error) {
	res, err := u.accessKeyAPI.GenerateAccessKey(ctx, cmd)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserAccessKeyAPI, "generate access key ", err).FormatPostErr()
	}

	return res, nil
}

// AccessKeys retrieves a paginated list of access keys via the user access keys API.
// The response includes access keys and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to unmarshal into a *queries.AccessKeyPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) AccessKeys(ctx context.Context, accessKeyOpts ...queries.QueryOption[filter.AccessKeyFilter]) (*queries.AccessKeyPage, error) {
	res, err := u.accessKeyAPI.AccessKeys(ctx, accessKeyOpts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.AdminAccessKeyAPI, "retrieve access keys page ", err).FormatGetErr()
	}

	return res, nil
}

// AccessKey retrieves the access key associated with the specified ID via the user access keys API.
// The response is expected to be unmarshaled into a *response.AccessKey struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) AccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	res, err := u.accessKeyAPI.AccessKey(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("retrieve access key with ID: %s", ID)
		return nil, errutil.NewHTTPErrorFormatter(constants.UserAccessKeyAPI, msg, err).FormatGetErr()
	}

	return res, nil
}

// RevokeAccessKey revokes the access key associated with the given ID via the user access keys API.
// If the request fails or the response cannot be processed, an error is returned.
// A nil error indicates the revoking access key was successful.
func (u *UserAPI) RevokeAccessKey(ctx context.Context, ID string) error {
	err := u.accessKeyAPI.RevokeAccessKey(ctx, ID)
	if err != nil {
		msg := fmt.Sprintf("revoke access key with ID: %s", ID)
		return errutil.NewHTTPErrorFormatter(constants.AdminAccessKeyAPI, msg, err).FormatDeleteErr()
	}

	return nil
}

// UTXOs fetches a paginated list of UTXOs via the user UTXOs API.
// The response includes UTXOs along with pagination details, such as page number,
// sort order, and sorting field.
//
// Optional query parameters can be applied using the provided query options.
// The response is unmarshaled into a *queries.UtxosPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) UTXOs(ctx context.Context, opts ...queries.QueryOption[filter.UtxoFilter]) (*queries.UtxosPage, error) {
	res, err := u.utxosAPI.UTXOs(ctx, opts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserUtxosAPI, "retrieve UTXOs page", err).FormatGetErr()
	}

	return res, nil
}

// MerkleRoots retrieves a paginated list of Merkle roots via the user Merkle roots API.
// The API response includes Merkle roots along with pagination details, such as the current
// page number, sort order, and sorting field (sortBy).
//
// This method supports optional query parameters, which can be specified using the provided
// query options. These options customize the behavior of the API request, such as setting
// batch size or applying filters for pagination.
//
// The response is unmarshaled into a *queries.MerkleRootPage struct.
// Returns an error if the request fails or the response cannot be decoded.
func (u *UserAPI) MerkleRoots(ctx context.Context, opts ...queries.MerkleRootsQueryOption) (*queries.MerkleRootPage, error) {
	res, err := u.merkleRootsAPI.MerkleRoots(ctx, opts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserMerkleRootAPI, "retrieve Merkle root page", err).FormatGetErr()
	}

	return res, nil
}

// SyncMerkleRoots synchronizes Merkle roots known to the SPV Wallet with the client database.
// This method sends a series of HTTP GET requests to the "/merkleroots" endpoint, fetching
// Merkle roots and storing them in the client database. The process continues until all
func (u *UserAPI) SyncMerkleRoots(ctx context.Context, repo merkleroots.MerkleRootsRepository) error {
	err := u.merkleRootsAPI.SyncMerkleRoots(ctx, repo)
	if err != nil {
		return fmt.Errorf("failed to sync Merkle roots: %w", err)
	}

	return nil
}

// GenerateTotpForContact generates a TOTP code for the specified contact.
func (u *UserAPI) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	if u.totpAPI == nil {
		return "", errors.New("totp client not initialized - xPriv authentication required")
	}

	totp, err := u.totpAPI.GenerateTotpForContact(contact, period, digits)
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP for contact: %w", err)
	}

	return totp, nil
}

// ValidateTotpForContact validates a TOTP code for the specified contact.
func (u *UserAPI) ValidateTotpForContact(contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	if u.totpAPI == nil {
		return errors.New("totp client not initialized - xPriv authentication required")
	}

	if err := u.totpAPI.ValidateTotpForContact(contact, passcode, requesterPaymail, period, digits); err != nil {
		return fmt.Errorf("failed to validate TOTP for contact: %w", err)
	}

	return nil
}

// Paymails retrieves a paginated list of paymail addresses via the User Paymails API.
// The response includes user paymails along with pagination metadata, such as
// the current page number, sort order, and the field used for sorting (sortBy).
//
// Query parameters can be configured using optional query options. These options allow
// filtering based on metadata, pagination settings, or specific paymail attributes.
//
// The API response is unmarshaled into a *queries.PaymailAddressPage struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (u *UserAPI) Paymails(ctx context.Context, opts ...queries.QueryOption[filter.PaymailFilter]) (*queries.PaymailsPage, error) {
	res, err := u.paymailsAPI.Paymails(ctx, opts...)
	if err != nil {
		return nil, errutil.NewHTTPErrorFormatter(constants.UserPaymailAPI, "retrieve paymail addresses page", err).FormatGetErr()
	}

	return res, nil
}

// NewUserAPIWithXPub initializes a new UserAPI instance using an extended public key (xPub).
// This function configures the API client with the provided configuration and uses the xPub key for authentication.
// If any configuration or initialization step fails, an appropriate error is returned.
//
// Note: Requests made with this instance will not be signed.
// For enhanced security, it is strongly recommended to use `NewUserAPIWithXPriv` or `NewUserAPIWithAccessKey` instead.
func NewUserAPIWithXPub(cfg config.Config, xPub string) (*UserAPI, error) {
	authenticator, err := auth.NewXpubOnlyAuthenticator(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xPub authenticator: %w", err)
	}

	return initUserAPI(cfg, authenticator)
}

// NewUserAPIWithXPriv initializes a new UserAPI instance using an extended private key (xPriv).
// This function configures the API client with the provided configuration and uses the xPriv key for authentication.
// If any step fails, an appropriate error is returned.
//
// Note: Requests made with this instance will be securely signed.
func NewUserAPIWithXPriv(cfg config.Config, xPriv string) (*UserAPI, error) {
	authenticator, err := auth.NewXprivAuthenticator(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xPriv authenticator: %w", err)
	}

	userAPI, err := initUserAPIWithXPriv(cfg, xPriv, authenticator)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	return userAPI, nil
}

// NewUserAPIWithAccessKey initializes a new UserAPI instance using an access key.
// This function configures the API client and converts the provided access key from either hex or WIF format into a private key.
// This private key is used for authentication. If any step in the process fails, an appropriate error is returned.
//
// Note: Requests made with this instance will be securely signed.
func NewUserAPIWithAccessKey(cfg config.Config, accessKey string) (*UserAPI, error) {
	authenticator, err := auth.NewAccessKeyAuthenticator(accessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized access key authenticator: %w", err)
	}

	return initUserAPI(cfg, authenticator)
}

type authenticator interface {
	Authenticate(r *resty.Request) error
}

func initUserAPIWithXPriv(cfg config.Config, xPriv string, auth authenticator) (*UserAPI, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr to url.URL: %w", err)
	}

	httpClient := restyutil.NewHTTPClient(cfg, auth)
	transactionsAPI, err := transactions.NewAPIWithXPriv(url, httpClient, xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactionsAPI: %w", err)
	}

	totpAPI, err := totp.NewAPI(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to create totpAPI: %w", err)
	}

	return &UserAPI{
		merkleRootsAPI:  merkleroots.NewAPI(url, httpClient),
		configsAPI:      configs.NewAPI(url, httpClient),
		transactionsAPI: transactionsAPI,
		utxosAPI:        utxos.NewAPI(url, httpClient),
		accessKeyAPI:    accesskeys.NewAPI(url, httpClient),
		xpubAPI:         xpubs.NewAPI(url, httpClient),
		contactsAPI:     contacts.NewAPI(url, httpClient),
		invitationsAPI:  invitations.NewAPI(url, httpClient),
		paymailsAPI:     paymails.NewAPI(url, httpClient),
		totpAPI:         totpAPI,
	}, nil
}

func initUserAPI(cfg config.Config, auth authenticator) (*UserAPI, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr to url.URL: %w", err)
	}

	httpClient := restyutil.NewHTTPClient(cfg, auth)
	if httpClient == nil {
		return nil, fmt.Errorf("failed to initialize HTTP client - nil value")
	}

	transactionsAPI, err := transactions.NewAPI(url, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactionsAPI: %w", err)
	}

	return &UserAPI{
		merkleRootsAPI:  merkleroots.NewAPI(url, httpClient),
		configsAPI:      configs.NewAPI(url, httpClient),
		transactionsAPI: transactionsAPI,
		utxosAPI:        utxos.NewAPI(url, httpClient),
		accessKeyAPI:    accesskeys.NewAPI(url, httpClient),
		xpubAPI:         xpubs.NewAPI(url, httpClient),
		contactsAPI:     contacts.NewAPI(url, httpClient),
		invitationsAPI:  invitations.NewAPI(url, httpClient),
		paymailsAPI:     paymails.NewAPI(url, httpClient),
	}, nil
}
