package config_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/stretchr/testify/require"
)

func TestConfig_New(t *testing.T) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		Proxy:               http.ProxyFromEnvironment,
	}

	tests := []struct {
		name     string
		options  []config.Option
		expected config.Config
	}{
		{
			name:    "All defaults",
			options: nil,
			expected: config.Config{
				Addr:      "http://localhost:3003",
				Timeout:   1 * time.Minute,
				Transport: http.DefaultTransport,
			},
		},
		{
			name: "Partial customization",
			options: []config.Option{
				config.WithAddr("http://api.example.com"),
			},
			expected: config.Config{
				Addr:      "http://api.example.com",
				Timeout:   1 * time.Minute,
				Transport: http.DefaultTransport,
			},
		},
		{
			name: "Full customization",
			options: []config.Option{
				config.WithAddr("http://custom.example.com"),
				config.WithTimeout(2 * time.Minute),
				config.WithTransport(transport),
			},
			expected: config.Config{
				Addr:      "http://custom.example.com",
				Timeout:   2 * time.Minute,
				Transport: transport,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.New(test.options...)
			require.Equal(t, test.expected, cfg)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.Config
		expectedErr error
	}{
		{
			name:        "Valid configuration with constructor defaults",
			cfg:         config.New(),
			expectedErr: nil,
		},
		{
			name: "Valid configuration",
			cfg: config.Config{
				Addr:      "http://api.example.com",
				Timeout:   30 * time.Second,
				Transport: http.DefaultTransport,
			},
			expectedErr: nil, // No error expected
		},
		{
			name: "Missing Addr",
			cfg: config.Config{
				Timeout:   30 * time.Second,
				Transport: http.DefaultTransport,
			},
			expectedErr: goclienterr.ErrConfigValidationMissingAddress,
		},
		{
			name: "Invalid Addr URL",
			cfg: config.Config{
				Addr:      "invalid-url",
				Timeout:   30 * time.Second,
				Transport: http.DefaultTransport,
			},
			expectedErr: goclienterr.ErrConfigValidationInvalidAddress,
		},
		{
			name: "Zero Timeout - default 1m",
			cfg: config.Config{
				Addr:      "http://api.example.com",
				Timeout:   0,
				Transport: http.DefaultTransport,
			},
			expectedErr: nil,
		},
		{
			name: "negative Timeout",
			cfg: config.Config{
				Addr:      "http://api.example.com",
				Timeout:   -10 * time.Second,
				Transport: http.DefaultTransport,
			},
			expectedErr: goclienterr.ErrConfigValidationInvalidTimeout,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.cfg.Validate()
			require.ErrorIs(t, test.expectedErr, err)
		})
	}
}
