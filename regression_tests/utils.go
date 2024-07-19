package regressiontests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/bitcoin-sv/spv-wallet/models"
)

const (
	domainPrefix             = "https://"
	domainSuffixSharedConfig = "/v1/shared-config"

	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"
)

var (
	explicitHTTPURLRegex = regexp.MustCompile(`^https?://`)
	envVariableError     = errors.New("missing xpriv variables")
)

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
	client := &http.Client{}
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
