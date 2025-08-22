package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConfig(t *testing.T) {
	config := &Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:9999/auth/callback",
		IssuerURL:    "https://accounts.google.com",
		CookieDomain: "",
		CookieSecure: false,
		StaticToken:  "test-token",
	}

	assert.Equal(t, "test-client-id", config.ClientID)
	assert.Equal(t, "test-client-secret", config.ClientSecret)
	assert.Equal(t, "http://localhost:9999/auth/callback", config.RedirectURL)
	assert.Equal(t, "https://accounts.google.com", config.IssuerURL)
	assert.Empty(t, config.CookieDomain)
	assert.False(t, config.CookieSecure)
	assert.Equal(t, "test-token", config.StaticToken)
}

func TestNewOIDCAuth_InvalidConfig(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "empty client ID",
			config: &Config{
				ClientID:     "",
				ClientSecret: "test-secret",
			},
		},
		{
			name: "empty client secret",
			config: &Config{
				ClientID:     "test-id",
				ClientSecret: "",
			},
		},
		{
			name: "both empty",
			config: &Config{
				ClientID:     "",
				ClientSecret: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := NewOIDCAuth(ctx, tt.config, logger)
			require.Error(t, err)
			assert.Nil(t, auth)
		})
	}
}

func TestUserInfo(t *testing.T) {
	userInfo := UserInfo{
		Email:         "test@example.com",
		Name:          "Test User",
		Picture:       "https://example.com/avatar.jpg",
		EmailVerified: true,
	}

	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.Equal(t, "https://example.com/avatar.jpg", userInfo.Picture)
	assert.True(t, userInfo.EmailVerified)
}

func TestGenerateRandomState(t *testing.T) {
	state1, err1 := generateRandomState()
	require.NoError(t, err1)
	assert.NotEmpty(t, state1)

	state2, err2 := generateRandomState()
	require.NoError(t, err2)
	assert.NotEmpty(t, state2)

	// States should be different
	assert.NotEqual(t, state1, state2)

	// States should be valid base64
	_, err := base64.URLEncoding.DecodeString(state1)
	require.NoError(t, err)
}

func TestGetUserFromRequest_NoSession(t *testing.T) {
	// Create a mock auth instance (we can't easily test with real OIDC)
	auth := &OIDCAuth{
		logger: zap.NewNop(),
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	user := auth.GetUserFromRequest(req)
	assert.Nil(t, user)
}

func TestGetUserFromRequest_ValidSession(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
	}

	// Create a test user
	testUser := UserInfo{
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: true,
	}

	// Marshal to JSON and encode
	userJSON, err := json.Marshal(testUser)
	require.NoError(t, err)

	encodedUser := base64.StdEncoding.EncodeToString(userJSON)

	// Create request with session cookie
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "user_session",
		Value: encodedUser,
	})

	user := auth.GetUserFromRequest(req)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Email, user.Email)
	assert.Equal(t, testUser.Name, user.Name)
	assert.Equal(t, testUser.EmailVerified, user.EmailVerified)
}

func TestGetUserFromRequest_InvalidSession(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "user_session",
		Value: "invalid-base64-data",
	})

	user := auth.GetUserFromRequest(req)
	assert.Nil(t, user)
}

func TestRequireAuth_NoAuth(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
		config: &Config{},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := auth.RequireAuth(handler)

	// Test with API request (no HTML accept header)
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()

	protectedHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Test with browser request (HTML accept header)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	rr = httptest.NewRecorder()

	protectedHandler(rr, req)
	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, "/auth/login", rr.Header().Get("Location"))
}

func TestRequireAuth_StaticToken(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
		config: &Config{
			StaticToken: "test-static-token",
		},
	}

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		// Check that user context is set
		user := r.Context().Value(userContextKey)
		assert.NotNil(t, user)

		userInfo, ok := user.(*UserInfo)
		assert.True(t, ok)
		assert.Equal(t, "static-token-user", userInfo.Email)

		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := auth.RequireAuth(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer test-static-token")
	rr := httptest.NewRecorder()

	protectedHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
}

func TestRequireAuth_InvalidStaticToken(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
		config: &Config{
			StaticToken: "test-static-token",
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := auth.RequireAuth(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	req.Header.Set("Accept", "application/json")
	rr := httptest.NewRecorder()

	protectedHandler(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRegisterRoutes(t *testing.T) {
	auth := &OIDCAuth{
		logger: zap.NewNop(),
		config: &Config{},
	}

	router := mux.NewRouter()
	auth.RegisterRoutes(router)

	// Test that routes are registered by making requests
	paths := []string{
		"/auth/login",
		"/auth/callback",
		"/api/user",
	}

	for _, path := range paths {
		var req *http.Request
		var method string

		switch path {
		case "/auth/logout":
			method = "POST"
		default:
			method = "GET"
		}

		req = httptest.NewRequest(method, path, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Should not return 404 for registered routes
		assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route %s should be registered", path)
	}
}
