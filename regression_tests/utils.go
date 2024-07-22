package regressiontests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

const (
	atSign                   = "@"
	domainPrefix             = "https://"
	domainSuffixSharedConfig = "/v1/shared-config"

	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"

	timeoutDuration = 30 * time.Second
)

var (
	explicitHTTPURLRegex = regexp.MustCompile(`^https?://`)
	envVariableError     = errors.New("missing xpriv variables")
)

type regressionTestUser struct {
	XPriv   string `json:"xpriv"`
	XPub    string `json:"xpub"`
	Paymail string `json:"paymail"`
}

type regressionTestConfig struct {
	ClientOneURL         string
	ClientTwoURL         string
	ClientOneLeaderXPriv string
	ClientTwoLeaderXPriv string
}

// getEnvVariables retrieves the environment variables needed for the regression tests.
func getEnvVariables() (*regressionTestConfig, error) {
	rtConfig := regressionTestConfig{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}

	if rtConfig.ClientOneLeaderXPriv == "" || rtConfig.ClientTwoLeaderXPriv == "" {
		return nil, envVariableError
	}
	if rtConfig.ClientOneURL == "" || rtConfig.ClientTwoURL == "" {
		rtConfig.ClientOneURL = "http://localhost:3003"
		rtConfig.ClientTwoURL = "http://localhost:3003"
	}
	return &rtConfig, nil
}

// getSharedConfig retrieves the shared configuration from the SPV Wallet.
func getSharedConfig(xpub string, clientUrl string) (*models.SharedConfig, error) {
	req, err := http.NewRequest(http.MethodGet, clientUrl+domainSuffixSharedConfig, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(models.AuthHeader, xpub)
	client := http.Client{
		Timeout: timeoutDuration,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get shared config: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var configResponse models.SharedConfig
	if err := json.Unmarshal(body, &configResponse); err != nil {
		return nil, err
	}

	if len(configResponse.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected 1 paymail domain, got %d", len(configResponse.PaymailDomains))
	}
	return &configResponse, nil
}

// createUser creates a set of keys and new paymail in the SPV Wallet.
func createUser(ctx context.Context, paymail string, paymailDomain string, instanceUrl string, adminXPriv string) (*regressionTestUser, error) {
	keys, err := xpriv.Generate()
	if err != nil {
		return nil, err
	}

	user := &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: preparePaymail(paymail, paymailDomain),
	}

	adminClient := walletclient.NewWithAdminKey(addPrefixIfNeeded(instanceUrl), adminXPriv)

	if err := adminClient.AdminNewXpub(ctx, user.XPub, map[string]any{"some_metadata": "remove"}); err != nil {
		return nil, err
	}

	_, err = adminClient.AdminCreatePaymail(ctx, user.XPub, user.Paymail, "Regression tests", "")
	if err != nil {
		return nil, err
	}

	return user, nil
}

// deleteUser soft deletes paymail from the SPV Wallet.
func deleteUser(ctx context.Context, paymail string, instanceURL string, adminXPriv string) error {
	adminClient := walletclient.NewWithAdminKey(addPrefixIfNeeded(instanceURL), adminXPriv)
	err := adminClient.AdminDeletePaymail(ctx, paymail)
	if err != nil {
		return err
	}
	return nil
}

// getBalance retrieves the balance from the SPV Wallet.
func getBalance(ctx context.Context, fromInstance string, fromXPriv string) (int, error) {
	client := walletclient.NewWithXPriv(addPrefixIfNeeded(fromInstance), fromXPriv)

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		return -1, err
	}
	return int(xpubInfo.CurrentBalance), nil
}

// getTransactions retrieves the transactions from the SPV Wallet.
func getTransactions(ctx context.Context, fromInstance string, fromXPriv string) ([]*models.Transaction, error) {
	client := walletclient.NewWithXPriv(addPrefixIfNeeded(fromInstance), fromXPriv)

	metadata := map[string]any{}
	conditions := filter.TransactionFilter{}
	queryParams := filter.QueryParams{}

	txs, err := client.GetTransactions(ctx, &conditions, metadata, &queryParams)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// sendFunds sends funds from one paymail to another.
func sendFunds(ctx context.Context, fromInstance string, fromXPriv string, toPamail string, howMuch int) (*models.Transaction, error) {
	client := walletclient.NewWithXPriv(fromInstance, fromXPriv)

	balance, err := getBalance(ctx, fromInstance, fromXPriv)
	if err != nil {
		return nil, err
	}
	if balance < howMuch {
		return nil, fmt.Errorf("insufficient funds: %d", balance)
	}

	recipient := walletclient.Recipients{To: toPamail, Satoshis: uint64(howMuch)}
	recipients := []*walletclient.Recipients{&recipient}
	metadata := map[string]any{
		"description": "regression-test",
	}

	transaction, err := client.SendToRecipients(ctx, recipients, metadata)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// preparePaymail prepares the paymail address by combining the alias and domain.
func preparePaymail(paymailAlias string, domain string) string {
	if isValidURL(domain) {
		splitedDomain := strings.SplitAfter(domain, "//")
		domain = splitedDomain[1]
	}
	url := paymailAlias + atSign + domain
	return url
}

// addPrefixIfNeeded adds the HTTPS prefix to the URL if it is missing.
func addPrefixIfNeeded(url string) string {
	if !isValidURL(url) {
		return domainPrefix + url
	}
	return url
}

// isValidURL validates the URL if it has http or https prefix.
func isValidURL(rawURL string) bool {
	return explicitHTTPURLRegex.MatchString(rawURL)
}
