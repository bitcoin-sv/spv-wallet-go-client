package spvwallet

import (
	"context"
	"fmt"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/invitations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/xpubs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
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
	xpubsAPI       *xpubs.API // Internal API for managing operations related to XPubs.
	contactsAPI    *contacts.API
	invitationsAPI *invitations.API
}

// CreateXPub creates a new XPub record via the Admin XPubs API.
// The provided command contains the necessary parameters to define the XPub record.
//
// The API response is unmarshaled into a *response.Xpub struct.
// Returns an error if the API request fails or the response cannot be decoded.
func (a *AdminAPI) CreateXPub(ctx context.Context, cmd *commands.CreateUserXpub) (*response.Xpub, error) {
	res, err := a.xpubsAPI.CreateXPub(ctx, cmd)
	if err != nil {
		return nil, xpubs.HTTPErrorFormatter("failed to create XPub", err).FormatPostErr()
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
		return nil, xpubs.HTTPErrorFormatter("failed to retrieve XPubs page", err).FormatGetErr()
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

// NewAdminAPIWithXPriv initializes a new AdminAPI instance using an extended private key (xPriv).
// This function configures the API client with the provided configuration and uses the xPriv key for authentication.
// If any step fails, an appropriate error is returned.
//
// Note: Requests made with this instance will be securely signed.
func NewAdminAPIWithXPriv(cfg config.Config, xPriv string) (*AdminAPI, error) {
	key, err := bip32.GenerateHDKeyFromString(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPriv: %w", err)
	}

	authenticator, err := auth.NewXprivAuthenticator(key)
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
	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
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
	return &AdminAPI{
		xpubsAPI:       xpubs.NewAPI(url, httpClient),
		contactsAPI:    contacts.NewAPI(url, httpClient),
		invitationsAPI: invitations.NewAPI(url, httpClient),
	}, nil
}
