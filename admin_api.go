package spvwallet

import (
	"context"
	"fmt"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	xpubs "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/users"
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
	xpubsAPI *xpubs.API // Internal API for managing operations related to XPubs.
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
	return &AdminAPI{xpubsAPI: xpubs.NewAPI(url, httpClient)}, nil
}
