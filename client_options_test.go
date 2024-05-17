package walletclient

import "testing"

func TestValidateAndCleanURL(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		expected string
		wantErr  bool
	}{
		{"Empty URL", "", "", true},
		{"Valid URL with path", "http://example.com/path", "http://example.com", false},
		{"Valid URL without path", "http://example.com", "http://example.com", false},
		{"Valid URL with port", "http://example.com:8080", "http://example.com:8080", false},
		{"Invalid URL", "http://%41:8080/", "", true},
		{"HTTPS URL", "https://example.com", "https://example.com", false},
		{"HTTPS URL with path", "https://example.com/path", "https://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndCleanURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndCleanURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("validateAndCleanURL() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
