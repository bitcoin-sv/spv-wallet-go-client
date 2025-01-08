package regressiontests

import (
	"fmt"
)

// spvWalletServerConfig contains configuration settings for initializing a SPVWalletAPI instance.
// These include the environment URL and private keys required for admin and user operations.
type spvWalletServerConfig struct {
	envURL     string // URL of the SPV Wallet API environment.
	envXPriv   string // Extended private key (xPriv) for the user account.
	adminXPriv string // Extended private key (xPriv) for the admin account.
}

// Validate validates the spvWalletServerConfig.
// It ensures that required fields like EnvURL and keys are not empty.
func (c *spvWalletServerConfig) validate() error {
	if c.envURL == "" {
		return fmt.Errorf("validation failed: environment URL is required")
	}

	if c.adminXPriv == "" {
		return fmt.Errorf("validation failed: admin xPriv is required")
	}

	if c.envXPriv == "" {
		return fmt.Errorf("validation failed: leader xPriv is required")
	}

	return nil
}

// spvWalletServer represents the core API for interacting with the SPV Wallet ecosystem.
// It holds configuration and client instances for admin, user, and leader operations.
type spvWalletServer struct {
	cfg    *spvWalletServerConfig // Configuration for the SPV Wallet API Config (e.g., environment URL, keys).
	admin  *admin                 // Admin client for performing administrative tasks like creating xPubs and paymails.
	user   *user                  // User client for standard wallet operations, such as transactions and balance retrieval.
	leader *user                  // Leader user client with potentially elevated privileges, managing broader wallet operations.
}

// setPaymailDomains sets SPV Wallet server clients to have their paymail addresses with the given domain address part.
func (s *spvWalletServer) setPaymailDomains(domain string) {
	type paymailSetter interface{ setPaymail(string) }

	clients := []paymailSetter{s.leader, s.admin, s.user}
	for _, client := range clients {
		client.setPaymail(domain)
	}
}

// initSPVWalletServer initializes the spvWalletAPI with Admin, Leader, and User clients.
// It accepts user alias and spvWalletServerConfig to be created as input parameters.
// On success, it returns an initialized SPVWalletAPI instance and nil error.
// If initialization of any component fails, a non-nil error is returned.
func initSPVWalletServer(alias string, cfg *spvWalletServerConfig) (*spvWalletServer, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	admin, err := initAdmin(cfg.envURL, cfg.adminXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin: %w", err)
	}

	leader, err := initUserWithXPriv("Leader", cfg.envURL, cfg.envXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize leader user: %w", err)
	}

	user, err := initUser(alias, cfg.envURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user %q: %w", alias, err)
	}

	return &spvWalletServer{
		cfg:    cfg,
		admin:  admin,
		leader: leader,
		user:   user,
	}, nil
}
