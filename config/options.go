package config

import (
	"net/http"
	"strings"
	"time"
)

// Option defines a function signature for modifying a Config.
type Option func(*Config)

// WithAddr sets the address in the configuration.
func WithAddr(addr string) Option {
	return func(cfg *Config) {
		cfg.Addr = strings.TrimSpace(addr)
	}
}

// WithTimeout sets the timeout duration in the configuration.
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

// WithTransport sets the HTTP transport in the configuration.
func WithTransport(transport http.RoundTripper) Option {
	return func(cfg *Config) {
		cfg.Transport = transport
	}
}
