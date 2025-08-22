package {{.ProjectPackageName}}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
		wantErr  bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expected: &Config{
				ServerAddress:    "0.0.0.0:9999",
				MetricsAddress:   "0.0.0.0:8080",
				LogLevel:         "info",
				OIDCClientID:     "",
				OIDCClientSecret: "",
				OIDCIssuerURL:    "https://accounts.google.com",
				OIDCRedirectURL:  "http://localhost:9999/auth/callback",
				OIDCCookieDomain: "",
				OIDCCookieSecure: false,
				OIDCStaticToken:  "",
			},
			wantErr: false,
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"EXAMPLE_SPA_SERVER_ADDRESS":     "127.0.0.1:8080",
				"EXAMPLE_SPA_METRICS_ADDRESS":    "127.0.0.1:9090",
				"EXAMPLE_SPA_LOG_LEVEL":          "debug",
				"EXAMPLE_SPA_OIDC_CLIENT_ID":     "test-client-id",
				"EXAMPLE_SPA_OIDC_CLIENT_SECRET": "test-client-secret",
				"EXAMPLE_SPA_OIDC_REDIRECT_URL":  "https://example.com/callback",
				"EXAMPLE_SPA_OIDC_COOKIE_SECURE": "true",
				"EXAMPLE_SPA_OIDC_STATIC_TOKEN":  "test-token",
			},
			expected: &Config{
				ServerAddress:    "127.0.0.1:8080",
				MetricsAddress:   "127.0.0.1:9090",
				LogLevel:         "debug",
				OIDCClientID:     "test-client-id",
				OIDCClientSecret: "test-client-secret",
				OIDCIssuerURL:    "https://accounts.google.com",
				OIDCRedirectURL:  "https://example.com/callback",
				OIDCCookieDomain: "",
				OIDCCookieSecure: true,
				OIDCStaticToken:  "test-token",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"EXAMPLE_SPA_LOG_LEVEL": "invalid",
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing client secret with client ID",
			envVars: map[string]string{
				"EXAMPLE_SPA_OIDC_CLIENT_ID": "test-client-id",
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "missing client ID with client secret",
			envVars: map[string]string{
				"EXAMPLE_SPA_OIDC_CLIENT_SECRET": "test-client-secret",
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			config, err := LoadConfig()

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, config)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.expected, config)
			}

			// Environment variables are automatically cleaned up by t.Setenv
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ServerAddress:  "0.0.0.0:9999",
				MetricsAddress: "0.0.0.0:8080",
				LogLevel:       "info",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: &Config{
				ServerAddress:  "0.0.0.0:9999",
				MetricsAddress: "0.0.0.0:8080",
				LogLevel:       "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty server address",
			config: &Config{
				ServerAddress:  "",
				MetricsAddress: "0.0.0.0:8080",
				LogLevel:       "info",
			},
			wantErr: true,
		},
		{
			name: "empty metrics address",
			config: &Config{
				ServerAddress:  "0.0.0.0:9999",
				MetricsAddress: "",
				LogLevel:       "info",
			},
			wantErr: true,
		},
		{
			name: "OIDC client ID without secret",
			config: &Config{
				ServerAddress:    "0.0.0.0:9999",
				MetricsAddress:   "0.0.0.0:8080",
				LogLevel:         "info",
				OIDCClientID:     "test-client-id",
				OIDCClientSecret: "",
			},
			wantErr: true,
		},
		{
			name: "OIDC client secret without ID",
			config: &Config{
				ServerAddress:    "0.0.0.0:9999",
				MetricsAddress:   "0.0.0.0:8080",
				LogLevel:         "info",
				OIDCClientID:     "",
				OIDCClientSecret: "test-client-secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
