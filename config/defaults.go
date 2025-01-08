package config

import (
	"net/http"
	"time"
)

const (
	// defaultAddr is the default base address of the SPV Wallet API.
	defaultAddr string = "http://localhost:3003"
	// DefaultTimeout is the default HTTP requests timeout duration.
	defaultTimeout time.Duration = 1 * time.Minute
)

// setDefaultValues assigns default values to fields that are not explicitly set.
func (cfg *Config) setDefaultValues() {
	if cfg.Addr == "" {
		cfg.Addr = defaultAddr
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
}
