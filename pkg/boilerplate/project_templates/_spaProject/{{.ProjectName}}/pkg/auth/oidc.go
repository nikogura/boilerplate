package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const userContextKey contextKey = "user"

// Config holds OIDC authentication configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	IssuerURL    string
	CookieDomain string
	CookieSecure bool
	StaticToken  string // Bearer token for bypassing OIDC
}

// OIDCAuth provides OIDC authentication functionality.
type OIDCAuth struct {
	config   *Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   oauth2.Config
	logger   *zap.Logger
}

// UserInfo represents authenticated user information.
type UserInfo struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// NewOIDCAuth creates a new OIDC authentication handler.
func NewOIDCAuth(ctx context.Context, config *Config, logger *zap.Logger) (auth *OIDCAuth, err error) {
	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, errors.New("OIDC client ID and secret are required")
	}

	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	auth = &OIDCAuth{
		config:   config,
		provider: provider,
		verifier: verifier,
		oauth2:   oauth2Config,
		logger:   logger,
	}

	err = nil
	return auth, err
}

// RegisterRoutes registers authentication routes.
func (a *OIDCAuth) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", a.LoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", a.CallbackHandler).Methods("GET")
	router.HandleFunc("/auth/logout", a.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/user", a.UserHandler).Methods("GET")
}

// LoginHandler initiates OIDC login flow.
func (a *OIDCAuth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		a.logger.Error("Failed to generate state", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state in secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   a.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	url := a.oauth2.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// CallbackHandler handles OIDC callback.
func (a *OIDCAuth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		a.logger.Error("State cookie not found", zap.Error(err))
		http.Error(w, "State verification failed", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		a.logger.Error("State mismatch")
		http.Error(w, "State verification failed", http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	token, err := a.oauth2.Exchange(ctx, code)
	if err != nil {
		a.logger.Error("Failed to exchange code", zap.Error(err))
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		a.logger.Error("ID token not found in response")
		http.Error(w, "ID token not found", http.StatusInternalServerError)
		return
	}

	// Verify ID token
	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		a.logger.Error("Failed to verify ID token", zap.Error(err))
		http.Error(w, "Token verification failed", http.StatusUnauthorized)
		return
	}

	// Extract user info from claims
	var userInfo UserInfo
	err = idToken.Claims(&userInfo)
	if err != nil {
		a.logger.Error("Failed to extract claims", zap.Error(err))
		http.Error(w, "Claims extraction failed", http.StatusInternalServerError)
		return
	}

	// Store user session (in production, use proper session storage)
	userJSON, _ := json.Marshal(userInfo)
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    base64.StdEncoding.EncodeToString(userJSON),
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   a.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	a.logger.Info("User authenticated", zap.String("email", userInfo.Email))
	http.Redirect(w, r, "/", http.StatusFound)
}

// LogoutHandler handles user logout.
func (a *OIDCAuth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	a.logger.Info("User logged out")
	http.Redirect(w, r, "/", http.StatusFound)
}

// UserHandler returns current user information.
func (a *OIDCAuth) UserHandler(w http.ResponseWriter, r *http.Request) {
	userInfo := a.GetUserFromRequest(r)
	if userInfo == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user":           userInfo.Email,
		"name":           userInfo.Name,
		"email_verified": userInfo.EmailVerified,
	})
}

// RequireAuth middleware that requires authentication.
func (a *OIDCAuth) RequireAuth(next http.HandlerFunc) (middleware http.HandlerFunc) {
	middleware = func(w http.ResponseWriter, r *http.Request) {
		// Check for static token bypass
		if a.config.StaticToken != "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if token == a.config.StaticToken {
					// Create anonymous user context
					ctx := context.WithValue(r.Context(), userContextKey, &UserInfo{
						Email: "static-token-user",
						Name:  "Static Token User",
					})
					next(w, r.WithContext(ctx))
					return
				}
			}
		}

		userInfo := a.GetUserFromRequest(r)
		if userInfo == nil {
			// Redirect to login for browser requests
			if strings.Contains(r.Header.Get("Accept"), "text/html") {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			// Return 401 for API requests
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)
		next(w, r.WithContext(ctx))
	}

	return middleware
}

// GetUserFromRequest extracts user information from request.
func (a *OIDCAuth) GetUserFromRequest(r *http.Request) (userInfo *UserInfo) {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return userInfo
	}

	userJSON, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return userInfo
	}

	userInfo = &UserInfo{}
	err = json.Unmarshal(userJSON, userInfo)
	if err != nil {
		userInfo = nil
		return userInfo
	}

	return userInfo
}

// generateRandomState generates a random state string for CSRF protection.
func generateRandomState() (state string, err error) {
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		state = ""
		return state, err
	}
	state = base64.URLEncoding.EncodeToString(bytes)
	err = nil
	return state, err
}
